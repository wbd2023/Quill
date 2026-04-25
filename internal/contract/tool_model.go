package contract

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	ToolVersionGoCommand  ToolVersionKind = "go_command"
	ToolVersionBuildInfo  ToolVersionKind = "build_info"
	ToolVersionShellcheck ToolVersionKind = "shellcheck"
	ToolVersionNodeCLI    ToolVersionKind = "node_cli"
)

const (
	ToolInstallNone              ToolInstallKind = "none"
	ToolInstallGoBinary          ToolInstallKind = "go_binary"
	ToolInstallNodePackage       ToolInstallKind = "node_package"
	ToolInstallShellcheckArchive ToolInstallKind = "shellcheck_archive"
)

const (
	ToolGo           = "go"
	ToolGoimports    = "goimports"
	ToolMisspell     = "misspell"
	ToolGolangciLint = "golangci-lint"
	ToolShfmt        = "shfmt"
	ToolShellcheck   = "shellcheck"
	ToolMarkdownlint = "markdownlint"
)

/* -------------------------------------------- Types ------------------------------------------- */

type ToolVersionKind string

type ToolInstallKind string

type Tool struct {
	ID            string
	Name          string
	Command       string
	PinnedVersion string
	VersionKind   ToolVersionKind
	ModulePath    string
	InstallKind   ToolInstallKind
	InstallSource string
}
