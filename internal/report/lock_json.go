package report

/* ------------------------------------------ JSON DTOs ----------------------------------------- */

type lockJSON struct {
	Path         string `json:"path"`
	ArchiveCount int    `json:"archive_count"`
}

/* ------------------------------------------ Rendering ----------------------------------------- */

func newLockJSON(result LockResult) (payload lockJSON) {
	return lockJSON(result)
}
