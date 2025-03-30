// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package database

import (
	"context"
	"database/sql"
)

const addConflict = `-- name: AddConflict :one
INSERT INTO conflicts (
  source_path, target_path, source_hash, target_hash, source_time, target_time, detected_at, resolution_status
) VALUES ( ?, ?, ?, ?, ?, ?, ?, 'unresolved' )
ON CONFLICT DO NOTHING
RETURNING id
`

type AddConflictParams struct {
	SourcePath string `json:"source_path"`
	TargetPath string `json:"target_path"`
	SourceHash string `json:"source_hash"`
	TargetHash string `json:"target_hash"`
	SourceTime int64  `json:"source_time"`
	TargetTime int64  `json:"target_time"`
	DetectedAt int64  `json:"detected_at"`
}

func (q *Queries) AddConflict(ctx context.Context, arg AddConflictParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, addConflict,
		arg.SourcePath,
		arg.TargetPath,
		arg.SourceHash,
		arg.TargetHash,
		arg.SourceTime,
		arg.TargetTime,
		arg.DetectedAt,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const addIgnorePattern = `-- name: AddIgnorePattern :one
INSERT INTO ignore_patterns (
  profile_id, pattern, type
) VALUES ( ?, ?, ? )
RETURNING id
`

type AddIgnorePatternParams struct {
	ProfileID int64  `json:"profile_id"`
	Pattern   string `json:"pattern"`
	Type      string `json:"type"`
}

func (q *Queries) AddIgnorePattern(ctx context.Context, arg AddIgnorePatternParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, addIgnorePattern, arg.ProfileID, arg.Pattern, arg.Type)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const addSyncRule = `-- name: AddSyncRule :one
INSERT INTO sync_rules (
  profile_id, source_dir, target_dir, created_at, updated_at
) VALUES ( ?, ?, ?, strftime('%s', 'now'), strftime('%s', 'now'))
RETURNING id
`

type AddSyncRuleParams struct {
	ProfileID int64  `json:"profile_id"`
	SourceDir string `json:"source_dir"`
	TargetDir string `json:"target_dir"`
}

func (q *Queries) AddSyncRule(ctx context.Context, arg AddSyncRuleParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, addSyncRule, arg.ProfileID, arg.SourceDir, arg.TargetDir)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createProfile = `-- name: CreateProfile :one
INSERT INTO profiles (
  name, created_at, updated_at
) VALUES (?, strftime('%s', 'now'), strftime('%s', 'now')) 
RETURNING id
`

func (q *Queries) CreateProfile(ctx context.Context, name string) (int64, error) {
	row := q.db.QueryRowContext(ctx, createProfile, name)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deleteFile = `-- name: DeleteFile :exec
DELETE FROM files
  WHERE source_path = ? AND target_path = ?
`

type DeleteFileParams struct {
	SourcePath string `json:"source_path"`
	TargetPath string `json:"target_path"`
}

func (q *Queries) DeleteFile(ctx context.Context, arg DeleteFileParams) error {
	_, err := q.db.ExecContext(ctx, deleteFile, arg.SourcePath, arg.TargetPath)
	return err
}

const deleteProfileByID = `-- name: DeleteProfileByID :exec
DELETE FROM profiles 
  WHERE id = ?
`

func (q *Queries) DeleteProfileByID(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteProfileByID, id)
	return err
}

const deleteProfileByName = `-- name: DeleteProfileByName :execrows
DELETE FROM profiles 
  WHERE name = ?
`

func (q *Queries) DeleteProfileByName(ctx context.Context, name string) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteProfileByName, name)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const deleteSyncRuleByProfileName = `-- name: DeleteSyncRuleByProfileName :execrows
DELETE FROM sync_rules
  WHERE profile_id = (SELECT id FROM profiles WHERE name = ?) AND source_dir = ?
`

type DeleteSyncRuleByProfileNameParams struct {
	Name      string `json:"name"`
	SourceDir string `json:"source_dir"`
}

func (q *Queries) DeleteSyncRuleByProfileName(ctx context.Context, arg DeleteSyncRuleByProfileNameParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteSyncRuleByProfileName, arg.Name, arg.SourceDir)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const getConflict = `-- name: GetConflict :one
SELECT id, source_path, target_path, source_hash, target_hash, source_time, target_time, detected_at, resolution_status, resolved_at
  FROM conflicts 
  WHERE id = ?
`

func (q *Queries) GetConflict(ctx context.Context, id int64) (Conflict, error) {
	row := q.db.QueryRowContext(ctx, getConflict, id)
	var i Conflict
	err := row.Scan(
		&i.ID,
		&i.SourcePath,
		&i.TargetPath,
		&i.SourceHash,
		&i.TargetHash,
		&i.SourceTime,
		&i.TargetTime,
		&i.DetectedAt,
		&i.ResolutionStatus,
		&i.ResolvedAt,
	)
	return i, err
}

const getFile = `-- name: GetFile :one
SELECT id, source_path, target_path, hash, size, mod_time, last_synced
  FROM files
  WHERE source_path = ? AND target_path = ?
`

type GetFileParams struct {
	SourcePath string `json:"source_path"`
	TargetPath string `json:"target_path"`
}

func (q *Queries) GetFile(ctx context.Context, arg GetFileParams) (File, error) {
	row := q.db.QueryRowContext(ctx, getFile, arg.SourcePath, arg.TargetPath)
	var i File
	err := row.Scan(
		&i.ID,
		&i.SourcePath,
		&i.TargetPath,
		&i.Hash,
		&i.Size,
		&i.ModTime,
		&i.LastSynced,
	)
	return i, err
}

const getIgnorePattern = `-- name: GetIgnorePattern :one
SELECT id, profile_id, pattern, type
  FROM ignore_patterns
  WHERE id = ?
`

func (q *Queries) GetIgnorePattern(ctx context.Context, id int64) (IgnorePattern, error) {
	row := q.db.QueryRowContext(ctx, getIgnorePattern, id)
	var i IgnorePattern
	err := row.Scan(
		&i.ID,
		&i.ProfileID,
		&i.Pattern,
		&i.Type,
	)
	return i, err
}

const getProfile = `-- name: GetProfile :one
SELECT id, name, created_at, updated_at
  FROM profiles
  WHERE id = ?
`

func (q *Queries) GetProfile(ctx context.Context, id int64) (Profile, error) {
	row := q.db.QueryRowContext(ctx, getProfile, id)
	var i Profile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProfileByName = `-- name: GetProfileByName :one
SELECT id, name, created_at, updated_at
  FROM profiles
  WHERE name = ?
`

func (q *Queries) GetProfileByName(ctx context.Context, name string) (Profile, error) {
	row := q.db.QueryRowContext(ctx, getProfileByName, name)
	var i Profile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProfileIDBySourceDir = `-- name: GetProfileIDBySourceDir :one
SELECT profile_id
  FROM sync_rules
  WHERE source_dir = ? AND status IS NOT 'disabled'
  LIMIT 1
`

func (q *Queries) GetProfileIDBySourceDir(ctx context.Context, sourceDir string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getProfileIDBySourceDir, sourceDir)
	var profile_id int64
	err := row.Scan(&profile_id)
	return profile_id, err
}

const getSyncRule = `-- name: GetSyncRule :one
SELECT id, profile_id, source_dir, target_dir, status, last_run_successful, created_at, updated_at
  FROM sync_rules
  WHERE id = ?
`

func (q *Queries) GetSyncRule(ctx context.Context, id int64) (SyncRule, error) {
	row := q.db.QueryRowContext(ctx, getSyncRule, id)
	var i SyncRule
	err := row.Scan(
		&i.ID,
		&i.ProfileID,
		&i.SourceDir,
		&i.TargetDir,
		&i.Status,
		&i.LastRunSuccessful,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSyncStatusByProfileName = `-- name: GetSyncStatusByProfileName :one
SELECT p.name as profile_name,
  sr.id, sr.profile_id, sr.source_dir, sr.target_dir, sr.status, sr.last_run_successful, sr.created_at, sr.updated_at
  FROM sync_rules sr
  JOIN profiles p ON sr.profile_id = p.id
  WHERE p.name = ?
`

type GetSyncStatusByProfileNameRow struct {
	ProfileName       string       `json:"profile_name"`
	ID                int64        `json:"id"`
	ProfileID         int64        `json:"profile_id"`
	SourceDir         string       `json:"source_dir"`
	TargetDir         string       `json:"target_dir"`
	Status            string       `json:"status"`
	LastRunSuccessful sql.NullBool `json:"last_run_successful"`
	CreatedAt         int64        `json:"created_at"`
	UpdatedAt         int64        `json:"updated_at"`
}

func (q *Queries) GetSyncStatusByProfileName(ctx context.Context, name string) (GetSyncStatusByProfileNameRow, error) {
	row := q.db.QueryRowContext(ctx, getSyncStatusByProfileName, name)
	var i GetSyncStatusByProfileNameRow
	err := row.Scan(
		&i.ProfileName,
		&i.ID,
		&i.ProfileID,
		&i.SourceDir,
		&i.TargetDir,
		&i.Status,
		&i.LastRunSuccessful,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const isProfileExists = `-- name: IsProfileExists :one
SELECT COUNT(*) 
  FROM profiles 
  WHERE LOWER(name) = LOWER(?)
`

func (q *Queries) IsProfileExists(ctx context.Context, lower string) (int64, error) {
	row := q.db.QueryRowContext(ctx, isProfileExists, lower)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const listFiles = `-- name: ListFiles :many
SELECT f.id, f.source_path, f.target_path, f.hash, f.size, f.mod_time, f.last_synced 
  FROM files f
  JOIN sync_rules sr ON f.source_path LIKE sr.source_dir || '%'
  WHERE sr.profile_id = ?
`

func (q *Queries) ListFiles(ctx context.Context, profileID int64) ([]File, error) {
	rows, err := q.db.QueryContext(ctx, listFiles, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []File
	for rows.Next() {
		var i File
		if err := rows.Scan(
			&i.ID,
			&i.SourcePath,
			&i.TargetPath,
			&i.Hash,
			&i.Size,
			&i.ModTime,
			&i.LastSynced,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listIgnorePattern = `-- name: ListIgnorePattern :many
SELECT id, profile_id, pattern, type
  FROM ignore_patterns
  WHERE profile_id = ?
`

func (q *Queries) ListIgnorePattern(ctx context.Context, profileID int64) ([]IgnorePattern, error) {
	rows, err := q.db.QueryContext(ctx, listIgnorePattern, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []IgnorePattern
	for rows.Next() {
		var i IgnorePattern
		if err := rows.Scan(
			&i.ID,
			&i.ProfileID,
			&i.Pattern,
			&i.Type,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listProfiles = `-- name: ListProfiles :many
SELECT id, name, created_at, updated_at
  FROM profiles
  ORDER BY created_at ASC
`

func (q *Queries) ListProfiles(ctx context.Context) ([]Profile, error) {
	rows, err := q.db.QueryContext(ctx, listProfiles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Profile
	for rows.Next() {
		var i Profile
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSyncRules = `-- name: ListSyncRules :many
SELECT id, profile_id, source_dir, target_dir, status, last_run_successful, created_at, updated_at
  FROM sync_rules
  WHERE profile_id = ? AND status IS NOT 'disabled'
  ORDER BY source_dir
`

func (q *Queries) ListSyncRules(ctx context.Context, profileID int64) ([]SyncRule, error) {
	rows, err := q.db.QueryContext(ctx, listSyncRules, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SyncRule
	for rows.Next() {
		var i SyncRule
		if err := rows.Scan(
			&i.ID,
			&i.ProfileID,
			&i.SourceDir,
			&i.TargetDir,
			&i.Status,
			&i.LastRunSuccessful,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSyncRulesGroupByProfile = `-- name: ListSyncRulesGroupByProfile :many
SELECT sr.profile_id as pid,
  p.name as profile_name,
  COUNT(sr.id) as rule_count,
  sr.source_dir,
  sr.target_dir
  FROM sync_rules sr
  JOIN profiles p ON sr.profile_id = p.id
  GROUP BY sr.profile_id
  ORDER BY sr.profile_id
`

type ListSyncRulesGroupByProfileRow struct {
	Pid         int64  `json:"pid"`
	ProfileName string `json:"profile_name"`
	RuleCount   int64  `json:"rule_count"`
	SourceDir   string `json:"source_dir"`
	TargetDir   string `json:"target_dir"`
}

func (q *Queries) ListSyncRulesGroupByProfile(ctx context.Context) ([]ListSyncRulesGroupByProfileRow, error) {
	rows, err := q.db.QueryContext(ctx, listSyncRulesGroupByProfile)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListSyncRulesGroupByProfileRow
	for rows.Next() {
		var i ListSyncRulesGroupByProfileRow
		if err := rows.Scan(
			&i.Pid,
			&i.ProfileName,
			&i.RuleCount,
			&i.SourceDir,
			&i.TargetDir,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUnresolvedConflicts = `-- name: ListUnresolvedConflicts :many
SELECT c.id, c.source_path, c.target_path, c.source_hash, c.target_hash, c.source_time, c.target_time, c.detected_at, c.resolution_status, c.resolved_at 
  FROM conflicts c
  JOIN sync_rules sr ON c.source_path LIKE sr.source_dir || '%'
  WHERE sr.profile_id = ? AND c.resolution_status = 'unresolved'
  ORDER BY detected_at DESC
`

func (q *Queries) ListUnresolvedConflicts(ctx context.Context, profileID int64) ([]Conflict, error) {
	rows, err := q.db.QueryContext(ctx, listUnresolvedConflicts, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Conflict
	for rows.Next() {
		var i Conflict
		if err := rows.Scan(
			&i.ID,
			&i.SourcePath,
			&i.TargetPath,
			&i.SourceHash,
			&i.TargetHash,
			&i.SourceTime,
			&i.TargetTime,
			&i.DetectedAt,
			&i.ResolutionStatus,
			&i.ResolvedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeIgnorePattern = `-- name: RemoveIgnorePattern :exec
DELETE FROM ignore_patterns
  WHERE id = ?
`

func (q *Queries) RemoveIgnorePattern(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, removeIgnorePattern, id)
	return err
}

const resolveConflict = `-- name: ResolveConflict :exec
UPDATE conflicts 
  SET resolution_status = ?, resolved_at = ? 
  WHERE id = ?
`

type ResolveConflictParams struct {
	ResolutionStatus sql.NullString `json:"resolution_status"`
	ResolvedAt       sql.NullInt64  `json:"resolved_at"`
	ID               int64          `json:"id"`
}

func (q *Queries) ResolveConflict(ctx context.Context, arg ResolveConflictParams) error {
	_, err := q.db.ExecContext(ctx, resolveConflict, arg.ResolutionStatus, arg.ResolvedAt, arg.ID)
	return err
}

const updateProfileByID = `-- name: UpdateProfileByID :exec
UPDATE profiles
  SET name = ?, updated_at = strftime('%s', 'now')
  WHERE id = ?
`

type UpdateProfileByIDParams struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

func (q *Queries) UpdateProfileByID(ctx context.Context, arg UpdateProfileByIDParams) error {
	_, err := q.db.ExecContext(ctx, updateProfileByID, arg.Name, arg.ID)
	return err
}

const updateProfileByName = `-- name: UpdateProfileByName :execrows
UPDATE profiles
  SET name = ?, updated_at = strftime('%s', 'now')
  WHERE name = ?
`

type UpdateProfileByNameParams struct {
	Name   string `json:"name"`
	Name_2 string `json:"name_2"`
}

func (q *Queries) UpdateProfileByName(ctx context.Context, arg UpdateProfileByNameParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, updateProfileByName, arg.Name, arg.Name_2)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const updateSyncRule = `-- name: UpdateSyncRule :exec
UPDATE sync_rules
  SET status = ?, last_run_successful = ?, updated_at = strftime('%s', 'now')
  WHERE id = ?
`

type UpdateSyncRuleParams struct {
	Status            string       `json:"status"`
	LastRunSuccessful sql.NullBool `json:"last_run_successful"`
	ID                int64        `json:"id"`
}

func (q *Queries) UpdateSyncRule(ctx context.Context, arg UpdateSyncRuleParams) error {
	_, err := q.db.ExecContext(ctx, updateSyncRule, arg.Status, arg.LastRunSuccessful, arg.ID)
	return err
}

const upsertFile = `-- name: UpsertFile :exec
INSERT INTO files (
  source_path, target_path, hash, size, mod_time, last_synced
) VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT(source_path, target_path) 
DO UPDATE SET 
  hash = excluded.hash,
  size = excluded.size,
  mod_time = excluded.mod_time,
  last_synced = excluded.last_synced
`

type UpsertFileParams struct {
	SourcePath string `json:"source_path"`
	TargetPath string `json:"target_path"`
	Hash       string `json:"hash"`
	Size       int64  `json:"size"`
	ModTime    int64  `json:"mod_time"`
	LastSynced int64  `json:"last_synced"`
}

func (q *Queries) UpsertFile(ctx context.Context, arg UpsertFileParams) error {
	_, err := q.db.ExecContext(ctx, upsertFile,
		arg.SourcePath,
		arg.TargetPath,
		arg.Hash,
		arg.Size,
		arg.ModTime,
		arg.LastSynced,
	)
	return err
}
