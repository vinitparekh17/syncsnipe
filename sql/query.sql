-- name: CreateProfile :one
INSERT INTO profiles (
  name, created_at, updated_at
) VALUES (?, strftime('%s', 'now'), strftime('%s', 'now')) 
RETURNING id;

-- name: GetProfile :one
SELECT *
  FROM profiles
  WHERE id = ?;

-- name: IsProfileExists :one
SELECT COUNT(*) 
  FROM profiles 
  WHERE LOWER(name) = LOWER(?);

-- name: GetProfileIDByName :one
SELECT id
  FROM profiles
  WHERE name = ?;

-- name: ListProfiles :many
SELECT *
  FROM profiles
  ORDER BY created_at ASC;

-- name: UpdateProfileByID :exec
UPDATE profiles
  SET name = ?, updated_at = strftime('%s', 'now')
  WHERE id = ?;

-- name: UpdateProfileByName :execrows
UPDATE profiles
  SET name = ?, updated_at = strftime('%s', 'now')
  WHERE name = ?;

-- name: DeleteProfileByID :exec
DELETE FROM profiles 
  WHERE id = ?;

-- name: DeleteProfileByName :execrows
DELETE FROM profiles 
  WHERE name = ?;

-- name: AddSyncRule :one
INSERT INTO sync_rules (
  profile_id, source_dir, target_dir, created_at, updated_at
) VALUES ( ?, ?, ?, strftime('%s', 'now'), strftime('%s', 'now'))
RETURNING id;

-- name: GetSyncRule :one
SELECT *
  FROM sync_rules
  WHERE id = ?;

-- name: GetSyncStatusByProfileName :one
SELECT sr.id,
  p.name as profile_name,
  sr.source_dir,
  sr.target_dir,
  sr.status,
  sr.created_at,
  sr.updated_at
  FROM sync_rules sr
  JOIN profiles p ON sr.profile_id = p.id
  WHERE p.name = ?;

-- name: GetProfileIDBySourceDir :one
SELECT profile_id
  FROM sync_rules
  WHERE source_dir = ? AND status IS NOT 'disabled'
  LIMIT 1;

-- name: ListSyncRules :many
SELECT *
  FROM sync_rules
  WHERE profile_id = ? AND status IS NOT 'disabled'
  ORDER BY source_dir;

-- name: ListSyncRulesGroupByProfile :many
SELECT sr.profile_id as pid,
  p.name as profile_name,
  COUNT(sr.id) as rule_count,
  sr.source_dir,
  sr.target_dir
  FROM sync_rules sr
  JOIN profiles p ON sr.profile_id = p.id
  GROUP BY sr.profile_id
  ORDER BY sr.profile_id;

-- name: ListSyncRulesByProfileID :many
SELECT *
  FROM sync_rules
  WHERE profile_id = ? AND status IS NOT 0;

-- name: UpdateSyncRule :exec
UPDATE sync_rules
  SET status = ?, last_run_successful = ?, updated_at = strftime('%s', 'now')
  WHERE id = ?;

-- name: DeleteSyncRuleByProfileName :execrows
DELETE FROM sync_rules
  WHERE profile_id = (SELECT id FROM profiles WHERE name = ?) AND source_dir = ?;

-- name: UpsertFile :exec
INSERT INTO files (
  source_path, target_path, hash, size, mod_time, last_synced
) VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT(source_path, target_path) 
DO UPDATE SET 
  hash = excluded.hash,
  size = excluded.size,
  mod_time = excluded.mod_time,
  last_synced = excluded.last_synced;

-- name: GetFile :one
SELECT *
  FROM files
  WHERE source_path = ? AND target_path = ?;

-- name: ListFiles :many
SELECT f.* 
  FROM files f
  JOIN sync_rules sr ON f.source_path LIKE sr.source_dir || '%'
  WHERE sr.profile_id = ?;

-- name: DeleteFile :exec
DELETE FROM files
  WHERE source_path = ? AND target_path = ?;

-- name: AddConflict :one
INSERT INTO conflicts (
  source_path, target_path, source_hash, target_hash, source_time, target_time, detected_at, resolution_status
) VALUES ( ?, ?, ?, ?, ?, ?, ?, 0 )
ON CONFLICT DO NOTHING
RETURNING id;

-- name: GetConflict :one
SELECT *
  FROM conflicts 
  WHERE id = ?;

-- name: ListUnresolvedConflicts :many
SELECT c.* 
  FROM conflicts c
  JOIN sync_rules sr ON c.source_path LIKE sr.source_dir || '%'
  WHERE sr.profile_id = ? AND c.resolution_status = 0
  ORDER BY detected_at DESC;

-- name: ResolveConflict :exec
UPDATE conflicts 
  SET resolution_status = ?, resolved_at = ? 
  WHERE id = ?;

-- name: AddIgnorePattern :one
INSERT INTO ignore_patterns (
  profile_id, pattern, type
) VALUES ( ?, ?, ? )
RETURNING id;

-- name: GetIgnorePattern :one
SELECT *
  FROM ignore_patterns
  WHERE id = ?;

-- name: ListIgnorePattern :many
SELECT *
  FROM ignore_patterns
  WHERE profile_id = ?;

-- name: RemoveIgnorePattern :exec
DELETE FROM ignore_patterns
  WHERE id = ?;

-- name: RemoveIgnorePatternByProfileName :execrows
DELETE FROM ignore_patterns
  WHERE profile_id = ? AND pattern = ?;