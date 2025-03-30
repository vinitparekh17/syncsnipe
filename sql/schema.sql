PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA cache_size = -64000;

CREATE TABLE IF NOT EXISTS profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE, -- Profile name (e.g., "Work", "Pictures")
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);

CREATE TABLE IF NOT EXISTS files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_path TEXT NOT NULL,
    target_path TEXT NOT NULL,
    hash TEXT NOT NULL,
    size INTEGER NOT NULL,
    mod_time INTEGER NOT NULL,
    last_synced INTEGER NOT NULL,
    UNIQUE(source_path, target_path)
);

CREATE TABLE IF NOT EXISTS sync_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_id INTEGER NOT NULL,
    source_dir TEXT NOT NULL UNIQUE,
    target_dir TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'idle',
    last_run_successful BOOLEAN DEFAULT NULL, -- Track success or failure of last run
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    UNIQUE(profile_id, source_dir),
    CHECK (source_dir != target_dir),
    CHECK(status IN ('active', 'idle', 'scheduled', 'paused', 'disabled')),
    FOREIGN KEY(profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS conflicts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_path TEXT NOT NULL,
    target_path TEXT NOT NULL,
    source_hash TEXT NOT NULL,
    target_hash TEXT NOT NULL,
    source_time INTEGER NOT NULL,
    target_time INTEGER NOT NULL,
    detected_at INTEGER NOT NULL,
    resolution_status TEXT, -- 'unresolved', 'resolved_source', etc.
    resolved_at INTEGER,    -- Timestamp of resolution
    UNIQUE(source_path, target_path)
);

CREATE TABLE IF NOT EXISTS ignore_patterns (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_id INTEGER NOT NULL,
    pattern TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL DEFAULT 'glob', -- Glob, Regex, or Exact
    FOREIGN KEY(profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_files_source_path ON files(source_path);
CREATE INDEX IF NOT EXISTS idx_files_target_path ON files(target_path);
CREATE INDEX IF NOT EXISTS idx_conflicts_source_path ON conflicts(source_path);
CREATE INDEX IF NOT EXISTS idx_conflicts_detected_at ON conflicts(detected_at);
CREATE INDEX IF NOT EXISTS idx_sync_rules_status ON sync_rules(status);
CREATE INDEX IF NOT EXISTS idx_sync_rules_source_dir ON sync_rules(source_dir);
