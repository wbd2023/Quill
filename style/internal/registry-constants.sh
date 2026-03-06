#!/bin/bash
#
# tools/style/internal/registry-constants.sh
# Shared registry constants for style-check orchestration.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

# shellcheck disable=SC2034  # Consumed by sourced scripts.
readonly \
	STYLE_TIER_ONE="tier1" \
	STYLE_TIER_TWO="tier2" \
	STYLE_TIER_THREE="tier3" \
	STYLE_RUNNER_SCRIPT="script" \
	STYLE_RUNNER_SCRIPT_SCOPE="script_scope" \
	STYLE_RUNNER_EXECUTOR="runner" \
	STYLE_RUNNER_TARGET_GOLANGCI_APP="golangci_app" \
	STYLE_RUNNER_TARGET_GOLANGCI_TOOLS="golangci_tools" \
	STYLE_RUNNER_TARGET_AST_APP="ast_app" \
	STYLE_RUNNER_TARGET_AST_TOOLS="ast_tools"
