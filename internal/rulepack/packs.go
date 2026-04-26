package rulepack

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/toolchain"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	PackControl  = "control"
	PackText     = "text"
	PackMarkdown = "markdown"
	PackShell    = "shell"
	PackGo       = "go"
	PackSecurity = "security"
	PackNaming   = "naming"
)

const (
	RuleGroupControlPlane contract.RuleGroup = "control_plane"
	RuleGroupLanguage     contract.RuleGroup = "language_backends"
	RuleGroupText         contract.RuleGroup = "text_scanners"
	RuleGroupSecurity     contract.RuleGroup = "security_scanners"
	RuleGroupNaming       contract.RuleGroup = "naming_scanners"
	RuleGroupExternal     contract.RuleGroup = "external_tools"
)

const (
	ExecutorToolchain      contract.ExecutorKind = "toolchain"
	ExecutorControlPlane   contract.ExecutorKind = "control_plane"
	ExecutorFileCommand    contract.ExecutorKind = "file_command"
	ExecutorBackendCommand contract.ExecutorKind = "backend_command"
	ExecutorBackendCheck   contract.ExecutorKind = "backend_check"
	ExecutorRepositoryScan contract.ExecutorKind = "repository_scan"
)

const (
	BackendActionGoFormat = "go_format"
	BackendActionGolangci = "golangci"
)

const (
	ConfigRefArchitecture = "architecture"
	ConfigRefControlPlane = "control_plane"
	ConfigRefNaming       = "naming"
	ConfigRefRepository   = "repository"
)

const LanguageGo = "go"

const (
	ToolGo           = "go"
	ToolGoimports    = "goimports"
	ToolMisspell     = "misspell"
	ToolGolangciLint = "golangci-lint"
	ToolShfmt        = "shfmt"
	ToolShellcheck   = "shellcheck"
	ToolMarkdownlint = "markdownlint"
)

const (
	ToolVersionGoCommand  toolchain.VersionKind = "go_command"
	ToolVersionBuildInfo  toolchain.VersionKind = "build_info"
	ToolVersionShellcheck toolchain.VersionKind = "shellcheck"
	ToolVersionNodeCLI    toolchain.VersionKind = "node_cli"
)

const (
	ToolInstallNone              toolchain.InstallKind = "none"
	ToolInstallGoBinary          toolchain.InstallKind = "go_binary"
	ToolInstallNodePackage       toolchain.InstallKind = "node_package"
	ToolInstallShellcheckArchive toolchain.InstallKind = "shellcheck_archive"
)

const (
	GoCheckComments           = "comments"
	GoCheckData               = "data"
	GoCheckDomainIdentifiers  = "domain_identifiers"
	GoCheckErrors             = "errors"
	GoCheckGuardClauseSpacing = "guard_clause_spacing"
	GoCheckLogging            = "logging"
	GoCheckNaming             = "naming"
	GoCheckOrder              = "order"
	GoCheckParameters         = "parameters"
	GoCheckProcess            = "process"
	GoCheckResources          = "resources"
	GoCheckReturns            = "returns"
	GoCheckSecurity           = "security"
	GoCheckSwitchCaseSpacing  = "switch_case_spacing"
	GoCheckTests              = "tests"
)

/* -------------------------------------------- Types ------------------------------------------- */

type Pack struct {
	ID    string
	Name  string
	Tools []toolchain.Capability
	Rules []RuleDefinition
}

type RuleDefinition = contract.RuleDefinition

type ExecutionSpec = contract.ExecutionSpec

type Registry struct {
	packs        []Pack
	capabilities []toolchain.Capability
	rules        []RuleDefinition
}

/* ------------------------------------------ Accessors ----------------------------------------- */

func (registry Registry) Packs() (packs []Pack) {
	return append([]Pack{}, registry.packs...)
}

func (registry Registry) ToolCapabilities() (capabilities []toolchain.Capability) {
	return append([]toolchain.Capability{}, registry.capabilities...)
}

func (registry Registry) Tools() (tools []contract.Tool) {
	return toolchain.Policies(registry.capabilities)
}

func (registry Registry) Rules() (rules []RuleDefinition) {
	return append([]RuleDefinition{}, registry.rules...)
}

func (registry Registry) Definitions() (definitions contract.Definitions) {
	return contract.Definitions{
		Tools: registry.Tools(),
		Rules: registry.Rules(),
	}
}

func (registry Registry) ToolByID(id string) (capability toolchain.Capability, found bool) {
	for _, capability := range registry.capabilities {
		if capability.ID == id {
			return capability, true
		}
	}

	return toolchain.Capability{}, false
}
