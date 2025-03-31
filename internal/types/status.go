package types

// Status represents the current state of a sync rule in the database table sync_rules.status
//
//go:generate stringer -type=Status
type Status int

// Status values as enum
const (
	Active Status = iota
	Idle
	Scheduled
	Paused
	Disabled
)

func (s Status) IsValid() bool {
	switch s {
	case Active, Idle, Scheduled, Paused, Disabled:
		return true
	}
	return false
}
