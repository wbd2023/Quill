package builtin

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack"
)

// Pack describes one built-in Pack definition.
type Pack = pack.Definition

// PackConfig describes config accepted by a built-in Pack.
type PackConfig = pack.Config

type RuleDefinition = contract.RuleDefinition

type ExecutionSpec = contract.ExecutionSpec
