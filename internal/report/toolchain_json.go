package report

import "io"

type toolchainJSON struct {
	Result   toolchainResultJSON `json:"result"`
	AllValid bool                `json:"all_valid"`
}

type toolchainResultJSON struct {
	Statuses []toolStatusJSON `json:"statuses"`
}

type toolStatusJSON struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Path          string `json:"path"`
	Version       string `json:"version"`
	PinnedVersion string `json:"pinned_version"`
	Valid         bool   `json:"valid"`
	Issue         string `json:"issue,omitempty"`
}

func writeToolchainJSON(writer io.Writer, view ToolchainView) (allValid bool, err error) {
	err = writeJSON(writer, struct {
		Toolchain toolchainJSON `json:"toolchain"`
	}{Toolchain: newToolchainJSON(view)})
	return view.AllValid, err
}

func newToolchainJSON(view ToolchainView) (payload toolchainJSON) {
	payload = toolchainJSON{
		Result: toolchainResultJSON{
			Statuses: make([]toolStatusJSON, 0, len(view.Result.Statuses)),
		},
		AllValid: view.AllValid,
	}

	for _, status := range view.Result.Statuses {
		payload.Result.Statuses = append(payload.Result.Statuses, toolStatusJSON{
			ID:            status.Tool.ID,
			Name:          status.Tool.Name,
			Path:          status.Path,
			Version:       status.Version,
			PinnedVersion: status.Tool.PinnedVersion,
			Valid:         status.Valid,
			Issue:         status.Issue,
		})
	}

	return payload
}
