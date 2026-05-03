package rulepack

import "ciphera/tools/internal/contract"

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
	ConfigReferenceGo             = "go"
	ConfigReferenceQualitySurface = "quality_surface"
	ConfigReferenceRepository     = "repository"
	ConfigReferenceVocabulary     = "vocabulary"
)

const LanguageGo = "go"
