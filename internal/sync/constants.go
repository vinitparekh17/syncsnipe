package sync

import "time"

const (
	CreateOrModifyEvent string = "CREATE_OR_MODIFY"
	DeleteEvent         string = "DELETE"
	RenameEvent         string = "RENAME"

	DebounceTime time.Duration = 250 * time.Millisecond
)
