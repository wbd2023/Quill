package report

type lockJSON struct {
	Path         string `json:"path"`
	ArchiveCount int    `json:"archive_count"`
}

func newLockJSON(result LockResult) (payload lockJSON) {
	return lockJSON(result)
}
