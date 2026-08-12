#!/usr/bin/env python3
"""
SyncSnipe Watcher Load & Stress Test Suite
Simultaneously performs N file modifications across M files with configurable concurrency,
burst rates, subdirectories, chaos modes, and automatic integrity verification.
"""

import argparse
import concurrent.futures
import hashlib
import os
import random
import shutil
import string
import sys
import time
from dataclasses import dataclass, field
from pathlib import Path
from typing import Dict, List, Optional, Tuple

# ANSI Colors for terminal output
class Colors:
    HEADER = "\033[95m"
    BLUE = "\033[94m"
    CYAN = "\033[96m"
    GREEN = "\033[92m"
    YELLOW = "\033[93m"
    RED = "\033[91m"
    BOLD = "\033[1m"
    DIM = "\033[2m"
    RESET = "\033[0m"


@dataclass
class TestMetrics:
    total_creates: int = 0
    total_modifies: int = 0
    total_deletes: int = 0
    total_renames: int = 0
    total_bytes_written: int = 0
    start_time: float = 0.0
    end_time: float = 0.0
    errors: List[str] = field(default_factory=list)

    @property
    def total_ops(self) -> int:
        return self.total_creates + self.total_modifies + self.total_deletes + self.total_renames

    @property
    def duration_seconds(self) -> float:
        return max(0.001, self.end_time - self.start_time)

    @property
    def ops_per_second(self) -> float:
        return self.total_ops / self.duration_seconds

    @property
    def mb_per_second(self) -> float:
        return (self.total_bytes_written / (1024 * 1024)) / self.duration_seconds


def compute_sha256(filepath: Path) -> str:
    """Compute SHA-256 hash of a file."""
    hasher = hashlib.sha256()
    with open(filepath, "rb") as f:
        while chunk := f.read(65536):
            hasher.update(chunk)
    return hasher.hexdigest()


def generate_payload(size_bytes: int, seed_tag: str = "") -> bytes:
    """Generate random string payload with identifiable timestamp/seed."""
    header = f"--- SYNCSNIPE LOAD TEST | TAG: {seed_tag} | TIME: {time.time()} ---\n".encode("utf-8")
    if len(header) >= size_bytes:
        return header[:size_bytes]
    body_len = size_bytes - len(header)
    # Fast pseudo-random alphanumeric fill
    body = "".join(random.choices(string.ascii_letters + string.digits + " \n", k=body_len)).encode("utf-8")
    return header + body


def worker_modify_task(
    file_path: Path,
    num_changes: int,
    size_bytes: int,
    delay_ms: float,
    burst: bool,
    metrics: TestMetrics,
    lock: concurrent.futures.ThreadPoolExecutor,
) -> Tuple[int, int, int]:
    """Worker task modifying a single file N times."""
    creates = 0
    modifies = 0
    bytes_written = 0

    file_path.parent.mkdir(parents=True, exist_ok=True)

    # Initial creation
    initial_content = generate_payload(size_bytes, f"init_{file_path.name}")
    try:
        with open(file_path, "wb") as f:
            f.write(initial_content)
        creates += 1
        bytes_written += len(initial_content)
    except Exception as e:
        metrics.errors.append(f"Create error on {file_path}: {e}")
        return creates, modifies, bytes_written

    # Sequential modifications
    for c in range(1, num_changes + 1):
        if not burst and delay_ms > 0:
            time.sleep(delay_ms / 1000.0)

        payload = generate_payload(size_bytes, f"change_{c}_{file_path.name}")
        try:
            with open(file_path, "wb") as f:
                f.write(payload)
                f.flush()
                os.fsync(f.fileno())  # Ensure flush to trigger filesystem events immediately
            modifies += 1
            bytes_written += len(payload)
        except Exception as e:
            metrics.errors.append(f"Modify error on {file_path} (step {c}): {e}")

    return creates, modifies, bytes_written


def worker_chaos_task(
    file_path: Path,
    num_changes: int,
    size_bytes: int,
    delay_ms: float,
    metrics: TestMetrics,
) -> Tuple[int, int, int, int]:
    """Randomly creates, appends, overwrites, renames, and occasionally deletes files."""
    creates = 0
    modifies = 0
    deletes = 0
    renames = 0
    bytes_written = 0
    curr_path = file_path

    for c in range(num_changes):
        if delay_ms > 0:
            time.sleep(delay_ms / 1000.0)

        action = random.choices(["create_or_edit", "append", "rename", "delete"], weights=[50, 30, 10, 10])[0]

        try:
            if action == "create_or_edit" or not curr_path.exists():
                curr_path.parent.mkdir(parents=True, exist_ok=True)
                payload = generate_payload(size_bytes, f"chaos_{c}")
                with open(curr_path, "wb") as f:
                    f.write(payload)
                    f.flush()
                    os.fsync(f.fileno())
                creates += 1
                bytes_written += len(payload)

            elif action == "append":
                append_data = f"\n[APPEND_STEP_{c}_{time.time()}]\n".encode("utf-8")
                with open(curr_path, "ab") as f:
                    f.write(append_data)
                    f.flush()
                    os.fsync(f.fileno())
                modifies += 1
                bytes_written += len(append_data)

            elif action == "rename":
                new_path = curr_path.parent / f"{curr_path.stem}_renamed_{c}{curr_path.suffix}"
                curr_path.rename(new_path)
                curr_path = new_path
                renames += 1

            elif action == "delete":
                curr_path.unlink(missing_ok=True)
                deletes += 1

        except Exception as e:
            metrics.errors.append(f"Chaos error on {curr_path}: {e}")

    # Ensure a final valid file exists at the end of chaos so target verification can check it
    if not curr_path.exists():
        try:
            curr_path.parent.mkdir(parents=True, exist_ok=True)
            payload = generate_payload(size_bytes, f"chaos_final_{curr_path.name}")
            with open(curr_path, "wb") as f:
                f.write(payload)
            creates += 1
            bytes_written += len(payload)
        except Exception:
            pass

    return creates, modifies, deletes, renames


def run_load_test(
    source_dir: Path,
    num_files: int,
    changes_per_file: int,
    workers: int,
    size_kb: int,
    delay_ms: float,
    mode: str,
    subdirs_count: int,
) -> TestMetrics:
    """Execute concurrent load test against source directory."""
    metrics = TestMetrics()
    size_bytes = max(64, int(size_kb * 1024))

    # Pre-generate file paths (optionally distributed over subdirectories)
    file_paths: List[Path] = []
    for i in range(num_files):
        if subdirs_count > 0:
            subdir_name = f"folder_{(i % subdirs_count) + 1}"
            p = source_dir / subdir_name / f"load_test_file_{i:04d}.dat"
        else:
            p = source_dir / f"load_test_file_{i:04d}.dat"
        file_paths.append(p)

    print(f"\n{Colors.BOLD}{Colors.CYAN}=== Starting SyncSnipe Load Test ==={Colors.RESET}")
    print(f" {Colors.BOLD}Source Dir:{Colors.RESET}       {source_dir.resolve()}")
    print(f" {Colors.BOLD}Total Files (M):{Colors.RESET}  {num_files}")
    print(f" {Colors.BOLD}Changes/File (N):{Colors.RESET} {changes_per_file}")
    print(f" {Colors.BOLD}Total Ops Target:{Colors.RESET} {num_files * (changes_per_file + 1)}")
    print(f" {Colors.BOLD}Concurrency:{Colors.RESET}      {workers} worker threads")
    print(f" {Colors.BOLD}File Payload:{Colors.RESET}     {size_kb} KB ({size_bytes:,} bytes)")
    print(f" {Colors.BOLD}Test Mode:{Colors.RESET}        {mode.upper()}")
    print(f" {Colors.BOLD}Inter-Op Delay:{Colors.RESET}   {delay_ms} ms")
    print(f" {Colors.BOLD}Subdirectories:{Colors.RESET}   {subdirs_count}")
    print(f"{Colors.CYAN}{'-' * 42}{Colors.RESET}\n")

    metrics.start_time = time.perf_counter()

    with concurrent.futures.ThreadPoolExecutor(max_workers=workers) as executor:
        if mode == "chaos":
            futures = [
                executor.submit(
                    worker_chaos_task,
                    fp,
                    changes_per_file,
                    size_bytes,
                    delay_ms,
                    metrics,
                )
                for fp in file_paths
            ]
            for f in concurrent.futures.as_completed(futures):
                c, m, d, r = f.result()
                metrics.total_creates += c
                metrics.total_modifies += m
                metrics.total_deletes += d
                metrics.total_renames += r
        else:
            is_burst = mode == "burst"
            futures = [
                executor.submit(
                    worker_modify_task,
                    fp,
                    changes_per_file,
                    size_bytes,
                    delay_ms,
                    is_burst,
                    metrics,
                    executor,
                )
                for fp in file_paths
            ]
            for f in concurrent.futures.as_completed(futures):
                c, m, b = f.result()
                metrics.total_creates += c
                metrics.total_modifies += m
                metrics.total_bytes_written += b

    metrics.end_time = time.perf_counter()
    return metrics


def verify_target_sync(
    source_dir: Path,
    target_dir: Path,
    timeout_sec: float = 15.0,
    poll_interval: float = 0.5,
) -> Tuple[bool, int, int, List[str]]:
    """
    Polls and verifies that all files in source_dir exist in target_dir with identical SHA-256 hashes.
    """
    print(f"\n{Colors.BOLD}{Colors.YELLOW}>>> Waiting for SyncWatcher debounce (250ms) and worker sync queue...{Colors.RESET}")
    start_poll = time.time()
    source_files: Dict[str, str] = {}

    # Scan source files
    for root, _, files in os.walk(source_dir):
        for f in files:
            sp = Path(root) / f
            rel = sp.relative_to(source_dir).as_posix()
            source_files[rel] = compute_sha256(sp)

    total_source = len(source_files)
    print(f" Source active files to verify: {Colors.BOLD}{total_source}{Colors.RESET}")

    matched = 0
    mismatched = []
    missing = []

    while time.time() - start_poll < timeout_sec:
        matched = 0
        mismatched.clear()
        missing.clear()

        for rel, src_hash in source_files.items():
            tp = target_dir / rel
            if not tp.exists():
                missing.append(rel)
                continue
            try:
                tgt_hash = compute_sha256(tp)
                if tgt_hash == src_hash:
                    matched += 1
                else:
                    mismatched.append(f"{rel} (Hash mismatch: src {src_hash[:8]} vs tgt {tgt_hash[:8]})")
            except Exception:
                missing.append(rel)

        if matched == total_source and not mismatched:
            break

        time.sleep(poll_interval)

    success = matched == total_source and not mismatched
    issues = mismatched + [f"Missing in target: {m}" for m in missing]
    return success, matched, total_source, issues


def main():
    parser = argparse.ArgumentParser(
        description="SyncSnipe Multi-File Concurrent Load & Stress Testing Tool",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    parser.add_argument(
        "--source",
        "-s",
        type=Path,
        default=Path("./test-src"),
        help="Source directory watched by SyncSnipe",
    )
    parser.add_argument(
        "--target",
        "-t",
        type=Path,
        default=Path("./test-target"),
        help="Target directory where files are synced",
    )
    parser.add_argument(
        "--files",
        "-m",
        type=int,
        default=50,
        help="Number of distinct files (M) to generate",
    )
    parser.add_argument(
        "--changes",
        "-n",
        type=int,
        default=10,
        help="Number of modifications per file (N)",
    )
    parser.add_argument(
        "--workers",
        "-w",
        type=int,
        default=16,
        help="Number of concurrent worker threads",
    )
    parser.add_argument(
        "--size-kb",
        "-k",
        type=int,
        default=4,
        help="Size of each written file payload in KB",
    )
    parser.add_argument(
        "--delay-ms",
        "-d",
        type=float,
        default=0.0,
        help="Delay in milliseconds between edits on the same file (0 for max stress)",
    )
    parser.add_argument(
        "--mode",
        choices=["modify", "burst", "chaos"],
        default="modify",
        help="Load pattern: 'modify' (continuous updates), 'burst' (rapid back-to-back writes), 'chaos' (mix of creates, appends, renames, deletes)",
    )
    parser.add_argument(
        "--subdirs",
        type=int,
        default=5,
        help="Number of nested subdirectories to distribute files across",
    )
    parser.add_argument(
        "--clean",
        action="store_true",
        help="Clean the source directory before running test",
    )
    parser.add_argument(
        "--verify",
        action="store_true",
        default=True,
        help="Verify target directory consistency after load test completes",
    )
    parser.add_argument(
        "--no-verify",
        dest="verify",
        action="store_false",
        help="Skip verification phase",
    )
    parser.add_argument(
        "--timeout",
        type=float,
        default=20.0,
        help="Timeout in seconds to wait for target sync reconciliation",
    )

    args = parser.parse_args()

    # Prepare directories
    args.source.mkdir(parents=True, exist_ok=True)
    if args.clean:
        print(f"{Colors.YELLOW}Cleaning source directory {args.source}...{Colors.RESET}")
        for item in args.source.iterdir():
            if item.is_dir():
                shutil.rmtree(item)
            else:
                item.unlink(missing_ok=True)

    # Run load test
    metrics = run_load_test(
        source_dir=args.source,
        num_files=args.files,
        changes_per_file=args.changes,
        workers=args.workers,
        size_kb=args.size_kb,
        delay_ms=args.delay_ms,
        mode=args.mode,
        subdirs_count=args.subdirs,
    )

    # Print Load Test Summary
    print(f"\n{Colors.BOLD}{Colors.GREEN}=== Load Generation Complete ==={Colors.RESET}")
    print(f" {Colors.BOLD}Total Operations:{Colors.RESET}   {metrics.total_ops:,}")
    print(f"   - Creates:           {metrics.total_creates:,}")
    print(f"   - Modifications:     {metrics.total_modifies:,}")
    if metrics.total_renames or metrics.total_deletes:
        print(f"   - Renames:           {metrics.total_renames:,}")
        print(f"   - Deletes:           {metrics.total_deletes:,}")
    print(f" {Colors.BOLD}Total Data Written:{Colors.RESET} {metrics.total_bytes_written / 1024 / 1024:.2f} MB")
    print(f" {Colors.BOLD}Elapsed Time:{Colors.RESET}       {metrics.duration_seconds:.3f} s")
    print(f" {Colors.BOLD}Throughput:{Colors.RESET}         {Colors.CYAN}{metrics.ops_per_second:.1f} ops/sec{Colors.RESET} ({metrics.mb_per_second:.2f} MB/s)")

    if metrics.errors:
        print(f"\n{Colors.BOLD}{Colors.RED}Warnings/Errors Encountered ({len(metrics.errors)}):{Colors.RESET}")
        for err in metrics.errors[:5]:
            print(f"  - {err}")
        if len(metrics.errors) > 5:
            print(f"  ... and {len(metrics.errors) - 5} more")

    # Verification Phase
    if args.verify and args.target.exists():
        success, matched, total, issues = verify_target_sync(
            source_dir=args.source,
            target_dir=args.target,
            timeout_sec=args.timeout,
        )

        print(f"\n{Colors.BOLD}{Colors.CYAN}=== Target Synchronization Verification ==={Colors.RESET}")
        print(f" Target Directory:      {args.target.resolve()}")
        print(f" Verified Match Rate:   {matched}/{total} files ({(matched / max(1, total)) * 100:.1f}%)")

        if success:
            print(f" {Colors.BOLD}{Colors.GREEN}✓ PASS: All files in target perfectly match source SHA-256 hashes!{Colors.RESET}\n")
        else:
            print(f" {Colors.BOLD}{Colors.RED}✗ FAIL: Incomplete sync or hash mismatches within {args.timeout}s timeout:{Colors.RESET}")
            for issue in issues[:10]:
                print(f"  {Colors.YELLOW}• {issue}{Colors.RESET}")
            if len(issues) > 10:
                print(f"  ... and {len(issues) - 10} more")
            print()
            sys.exit(1)
    elif args.verify:
        print(f"\n{Colors.YELLOW}Note: Target directory {args.target} does not exist yet. Skipping verification.{Colors.RESET}")


if __name__ == "__main__":
    main()
