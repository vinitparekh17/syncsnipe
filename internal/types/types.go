//go:generate stringer -type=SyncStatus,ConflictResolutionStatus,IgnoreType -output=types_string.go
package types

import "errors"

// SyncStatus represents the current state of a sync rule in the database table sync_rules.status
type SyncStatus int

// Status values as enum
const (
	Active SyncStatus = iota
	Idle
	Scheduled
	Paused
	Disabled
)

func StringToSyncStatus(s string) (SyncStatus, error) {
	switch s {
	case "active":
		return Active, nil
	case "idle":
		return Idle, nil
	case "scheduled":
		return Scheduled, nil
	case "paused":
		return Paused, nil
	}
	return 0, errors.New("invalid sync status")
}

// ConflictResolutionStatus represents the current state of a conflict in the database table conflicts.resolution_status
type ConflictResolutionStatus int

const (
	Unresolved ConflictResolutionStatus = iota
	ResolvedSource
	ResolvedTarget
	Merged
)

func StringToConflictResolutionStatus(s string) (ConflictResolutionStatus, error) {
	switch s {
	case "unresolved":
		return Unresolved, nil
	case "resolved_source":
		return ResolvedSource, nil
	case "resolved_target":
		return ResolvedTarget, nil
	case "merged":
		return Merged, nil
	}
	return 0, errors.New("invalid conflict resolution status")
}

type IgnoreType int

const (
	Glob IgnoreType = iota
	Regex
	Exact
)

func StringToIgnoreType(s string) (IgnoreType, error) {
	switch s {
	case "glob":
		return Glob, nil
	case "regex":
		return Regex, nil
	case "exact":
		return Exact, nil
	}
	return 0, errors.New("invalid ignore type")
}
