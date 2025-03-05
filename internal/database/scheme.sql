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
    source_dir TEXT NOT NULL UNIQUE,
    target_dir TEXT NOT NULL,
    enabled BOOLEAN DEFAULT 1,
    last_run INTEGER,
    last_run_successful BOOLEAN DEFAULT NULL, -- Track success or failure of last run
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    CHECK (source_dir != target_dir)
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
    pattern TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL DEFAULT 'glob' -- Glob, Regex, or Exact
);

CREATE INDEX IF NOT EXISTS idx_files_source_path ON files(source_path);
CREATE INDEX IF NOT EXISTS idx_files_target_path ON files(target_path);
CREATE INDEX IF NOT EXISTS idx_conflicts_source_path ON conflicts(source_path);
CREATE INDEX IF NOT EXISTS idx_conflicts_detected_at ON conflicts(detected_at);
CREATE INDEX IF NOT EXISTS idx_sync_rules_enabled ON sync_rules(enabled);
CREATE INDEX IF NOT EXISTS idx_sync_rules_source_dir ON sync_rules(source_dir);
