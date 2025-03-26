package sync

import "time"

const (
	CREATE_OR_MODIFY string = "CREATE_OR_MODIFY"
	DELETE           string = "DELETE"
	RENAME           string = "RENAME"

	DEBOUNCE_TIME time.Duration = 250 * time.Millisecond
)
