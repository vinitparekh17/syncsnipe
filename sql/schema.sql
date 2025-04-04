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
    status INTEGER NOT NULL DEFAULT 1,
    last_run_successful BOOLEAN DEFAULT NULL, -- Track success or failure of last run
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    UNIQUE(profile_id, source_dir),
    CHECK (source_dir != target_dir),
    CHECK(status IN (0, 1, 2, 3, 4)), -- 0: active, 1: idle, 2: scheduled, 3: paused, 4: disabled
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
    resolution_status INTEGER NOT NULL DEFAULT 0 
        CHECK(resolution_status IN (0, 1, 2, 3)),
    resolved_at INTEGER,
    UNIQUE(source_path, target_path)
);

CREATE TABLE IF NOT EXISTS ignore_patterns (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_id INTEGER NOT NULL,
    pattern TEXT NOT NULL,
    type INTEGER NOT NULL DEFAULT 0 
        CHECK(type IN (0, 1, 2)), -- 0: glob, 1: regex, 2: exact
    UNIQUE(profile_id, pattern),  -- Ensures uniqueness within each profile
    FOREIGN KEY(profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_profiles_name ON profiles(name);
CREATE INDEX IF NOT EXISTS idx_files_source_path ON files(source_path);
CREATE INDEX IF NOT EXISTS idx_files_target_path ON files(target_path);
CREATE INDEX IF NOT EXISTS idx_conflicts_source_path ON conflicts(source_path);
CREATE INDEX IF NOT EXISTS idx_conflicts_detected_at ON conflicts(detected_at);
CREATE INDEX IF NOT EXISTS idx_sync_rules_status ON sync_rules(status);
CREATE INDEX IF NOT EXISTS idx_sync_rules_source_dir ON sync_rules(source_dir);
CREATE INDEX IF NOT EXISTS idx_ignore_patterns_profile_id ON ignore_patterns(profile_id);
