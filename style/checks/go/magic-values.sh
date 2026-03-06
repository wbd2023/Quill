#!/bin/bash
#
# tools/style/checks/go/magic-values.sh
# Checks for magic numeric values in Go files (STYLE.md 2.10).
#
# Rules:
#	- No unexplained numeric literals in Go logic.
#	- Enforced via golangci-lint mnd.
#
# Usage:
#	./tools/style/checks/go/magic-values.sh [--scope app|tools|all]
#
# Exit code: 0 if no violations, 1 if violations found, 2 if setup error.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
CACHE_ROOT="${TMPDIR:-/tmp}/ciphera-stylecheck-${USER:-user}"
GO_BUILD_CACHE="$CACHE_ROOT/go-build"
GOLANGCI_CACHE="$CACHE_ROOT/golangci-lint"
FAILED=0
LINTER_NAME="mnd"
LINTER_HELP_ENABLE_ONLY='--enable-only'
PACKAGE_APP_CMD="./cmd/..."
PACKAGE_APP_INTERNAL="./internal/..."
PACKAGE_TOOLS="./..."

# ---------------------------------------------- Paths ---------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
STYLE_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
# shellcheck source=tools/style/internal/common.sh
source "$STYLE_DIR/internal/common.sh"

PROJECT_ROOT="$(style_project_root_from_dir "$STYLE_DIR")"
STYLECHECK_DIRECTORY="$PROJECT_ROOT/$STYLE_PATH_TOOLS/style/ast"

# ---------------------------------------------- Args ----------------------------------------------

if ! SCOPE="$(style_parse_scope_arg "$USAGE_EXIT_CODE" "$STYLE_SCOPE_ALL" "$@")"; then
	exit "$USAGE_EXIT_CODE"
fi

mkdir -p "$GO_BUILD_CACHE" "$GOLANGCI_CACHE"

if ! style_require_command "golangci-lint" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

if ! style_validate_scope "$SCOPE" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

linter_flags=()
if golangci-lint run --help 2>/dev/null | grep -q -- "$LINTER_HELP_ENABLE_ONLY"; then
	linter_flags=(--default none --enable-only "$LINTER_NAME")
else
	linter_flags=(--disable-all --enable "$LINTER_NAME")
fi

run_magic_lint() {
	local working_directory="$1"
	shift

	if (
		cd "$working_directory" &&
			GOCACHE="$GO_BUILD_CACHE" \
				GOLANGCI_LINT_CACHE="$GOLANGCI_CACHE" \
				golangci-lint run "${linter_flags[@]}" "$@"
	); then
		return
	fi

	FAILED=1
}

case "$SCOPE" in
"$STYLE_SCOPE_APP")
	run_magic_lint "$PROJECT_ROOT" "$PACKAGE_APP_CMD" "$PACKAGE_APP_INTERNAL"
	;;
"$STYLE_SCOPE_TOOLS")
	run_magic_lint "$STYLECHECK_DIRECTORY" "$PACKAGE_TOOLS"
	;;
"$STYLE_SCOPE_ALL")
	run_magic_lint "$PROJECT_ROOT" "$PACKAGE_APP_CMD" "$PACKAGE_APP_INTERNAL"
	run_magic_lint "$STYLECHECK_DIRECTORY" "$PACKAGE_TOOLS"
	;;
esac

if [ "$FAILED" -eq 1 ]; then
	exit 1
fi

exit 0
