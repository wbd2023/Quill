package report

import "io"

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
