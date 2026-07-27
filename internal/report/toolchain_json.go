package report

import (
	"io"

	"github.com/wbd2023/quill/internal/toolchain"
)

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

func writeToolchainJSON(
	writer io.Writer,
	command string,
	view ToolchainView,
) (allValid bool, err error) {
	err = writeResultEnvelope(writer, command, newToolchainJSON(view))
	return view.AllValid, err
}

func newToolchainJSON(view ToolchainView) (payload toolchainJSON) {
	return toolchainJSON{
		Result:   toolchainResultJSON{Statuses: toolStatusListJSON(view.Result.Statuses)},
		AllValid: view.AllValid,
	}
}

// toolStatusListJSON builds the status payload shared by toolchain-driven commands (doctor,
// install, and fix).
func toolStatusListJSON(statuses []toolchain.Status) (payload []toolStatusJSON) {
	payload = make([]toolStatusJSON, 0, len(statuses))
	for _, status := range statuses {
		payload = append(payload, toolStatusJSON{
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
