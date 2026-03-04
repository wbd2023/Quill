#!/bin/bash
#
# tools/scripts/check-go-line-length.sh
# Checks Go line length using golangci-lint lll (STYLE.md 1.1).
#
# Rules:
#	- Go lines must be at most 100 characters.
#
# Usage:
#	./tools/scripts/check-go-line-length.sh [--scope app|tools|all]
#
# Exit code: 0 if no violations, 1 if violations found, 2 if invalid usage/tool unavailable.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
CACHE_ROOT="${TMPDIR:-/tmp}/ciphera-stylecheck-${USER:-user}"
GO_BUILD_CACHE="${GO_BUILD_CACHE:-$CACHE_ROOT/go-build}"
GOLANGCI_CACHE="${GOLANGCI_CACHE:-$CACHE_ROOT/golangci-lint}"
RULE_LABEL="1.1"
RULE_MESSAGE_PREFIX="[${RULE_LABEL}] Go line-length findings:"
LINTER_NAME="lll"
PACKAGE_APP_CMD="./cmd/..."
PACKAGE_APP_INTERNAL="./internal/..."
PACKAGE_TOOLS="./tools/..."
PACKAGE_ALL="./..."
LINTER_HELP_ENABLE_ONLY='--enable-only'

# ---------------------------------------------- Paths ---------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tools/scripts/lib/style-common.sh
source "$SCRIPT_DIR/lib/style-common.sh"

PROJECT_ROOT="$(style_project_root_from_dir "$SCRIPT_DIR")"

# ---------------------------------------------- Args ----------------------------------------------

if ! SCOPE="$(style_parse_scope_arg "$USAGE_EXIT_CODE" "$STYLE_SCOPE_ALL" "$@")"; then
	exit "$USAGE_EXIT_CODE"
fi

if ! style_require_command "golangci-lint" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

mkdir -p "$GO_BUILD_CACHE" "$GOLANGCI_CACHE"

linter_flags=()
if golangci-lint run --help 2>/dev/null | grep -q -- "$LINTER_HELP_ENABLE_ONLY"; then
	linter_flags=(--default none --enable-only "$LINTER_NAME")
else
	linter_flags=(--disable-all --enable "$LINTER_NAME")
fi

if ! style_validate_scope "$SCOPE" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

case "$SCOPE" in
"$STYLE_SCOPE_APP")
	packages=("$PACKAGE_APP_CMD" "$PACKAGE_APP_INTERNAL")
	;;
"$STYLE_SCOPE_TOOLS")
	packages=("$PACKAGE_TOOLS")
	;;
"$STYLE_SCOPE_ALL")
	packages=("$PACKAGE_ALL")
	;;
esac

if ! output=$(
	cd "$PROJECT_ROOT" && GOCACHE="$GO_BUILD_CACHE" GOLANGCI_LINT_CACHE="$GOLANGCI_CACHE" \
		golangci-lint run "${linter_flags[@]}" "${packages[@]}" 2>&1
); then
	echo "$RULE_MESSAGE_PREFIX"
	echo "$output"
	echo ""
	exit 1
fi

exit 0
