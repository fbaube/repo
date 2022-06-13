package repo

// Backupable methods use the current date & time
// and return the full path (or URI/URL) and error.
type Backupable interface {
	MoveToBackup() (string, error)
	CopyToBackup() (string, error)
	RestoreFromMostRecentBackup() (string, error)
}
