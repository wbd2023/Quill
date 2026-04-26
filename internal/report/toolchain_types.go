package report

import "ciphera/tools/internal/toolchain"

type ToolchainResult struct {
	Statuses []toolchain.Status
}

type ToolchainView struct {
	Result   ToolchainResult
	AllValid bool
}

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
