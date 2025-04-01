//go:generate stringer -type=SyncStatus,ConflictResolutionStatus -output=types_string.go
package types

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

func (s SyncStatus) IsValid() bool {
	switch s {
	case Active, Idle, Scheduled, Paused, Disabled:
		return true
	}
	return false
}

// ConflictResolutionStatus represents the current state of a conflict in the database table conflicts.resolution_status
type ConflictResolutionStatus int

const (
	Unresolved ConflictResolutionStatus = iota
	ResolvedSource
	ResolvedTarget
	Merged
)

func (s ConflictResolutionStatus) IsValid() bool {
	switch s {
	case Unresolved, ResolvedSource, ResolvedTarget, Merged:
		return true
	}
	return false
}
