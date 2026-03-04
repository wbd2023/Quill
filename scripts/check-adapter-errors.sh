#!/bin/bash
#
# tools/scripts/check-adapter-errors.sh
# Checks adapter error wrapping conventions (STYLE.md 2.1).
#
# Rules:
#   - Adapters should not return bare propagated errors.
#   - Low-level errors should be wrapped with context using %w.
#
# Exit code: 0 if no violations, 1 if violations found.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

FOUND=0
RULE_LABEL="2.1"
ADAPTER_ERROR_PATTERN='return\s+err$|return\s+.*,\s*err$'
GO_TEST_EXCLUDE='!*_test.go'
MESSAGE_ADAPTER_WRAP="[${RULE_LABEL}] Adapters must wrap low-level errors with context (%w):"
ADAPTERS_SUBDIRECTORY="adapters"

# ---------------------------------------------- Paths ---------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tools/scripts/lib/style-common.sh
source "$SCRIPT_DIR/lib/style-common.sh"

PROJECT_ROOT="$(style_project_root_from_dir "$SCRIPT_DIR")"

# --------------------------------------------- Checks ---------------------------------------------

# Bare error propagation patterns in adapters.
# Examples flagged:
#   return err
#   return nil, err
#   return value, found, err
bare_returns=$(rg -n --glob "$STYLE_PATTERN_GO" --glob "$GO_TEST_EXCLUDE" \
	"$ADAPTER_ERROR_PATTERN" \
	"$PROJECT_ROOT/$STYLE_PATH_INTERNAL/$ADAPTERS_SUBDIRECTORY" || true)

if [ -n "$bare_returns" ]; then
	echo "$MESSAGE_ADAPTER_WRAP"
	echo "$bare_returns"
	echo ""
	FOUND=1
fi

# ------------------------------------------- Validation -------------------------------------------

if [ "$FOUND" -eq 1 ]; then
	exit 1
fi

exit 0
