#!/bin/bash
#
# tools/scripts/check-spelling.sh
# Checks Australian-English-compatible spelling using misspell (UK locale).
#
# Rules:
#	- Non-Go text files should use Australian English conventions.
#	- Go files are checked separately via golangci-lint misspell.
#
# Usage:
#	./tools/scripts/check-spelling.sh [--scope app|tools|all]
#
# Exit code: 0 if no violations, 1 if violations found, 2 if tool unavailable.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
RULE_LABEL="1.2"
RULE_MESSAGE_PREFIX="[${RULE_LABEL}] Non-Go spelling issues (UK locale, AU-compatible):"
GO_FILE_EXCLUDE='\.go$'
LICENSE_FILE_EXCLUDE='(^|/)(LICENSE|COPYING|NOTICE)(\..+)?$'
ALL_SCOPE_EXCLUDE='^(vendor/|bin/|\.git/)'
APP_SCOPE_QUERY_CMD='cmd/**'
APP_SCOPE_QUERY_INTERNAL='internal/**'
TOOLS_SCOPE_QUERY='tools/**'
MISSPELL_LOCALE='UK'

# ---------------------------------------------- Paths ---------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tools/scripts/lib/style-common.sh
source "$SCRIPT_DIR/lib/style-common.sh"

PROJECT_ROOT="$(style_project_root_from_dir "$SCRIPT_DIR")"

# ---------------------------------------------- Args ----------------------------------------------

if ! SCOPE="$(style_parse_scope_arg "$USAGE_EXIT_CODE" "$STYLE_SCOPE_ALL" "$@")"; then
	exit "$USAGE_EXIT_CODE"
fi

if ! style_require_command "rg" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

if ! style_require_command "misspell" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

if ! style_validate_scope "$SCOPE" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

case "$SCOPE" in
"$STYLE_SCOPE_APP")
	mapfile -t FILES < <(
		git -C "$PROJECT_ROOT" ls-files "$APP_SCOPE_QUERY_CMD" "$APP_SCOPE_QUERY_INTERNAL" |
			rg -v "$GO_FILE_EXCLUDE" |
			rg -v "$LICENSE_FILE_EXCLUDE"
	)
	;;
"$STYLE_SCOPE_TOOLS")
	mapfile -t FILES < <(
		git -C "$PROJECT_ROOT" ls-files "$TOOLS_SCOPE_QUERY" |
			rg -v "$GO_FILE_EXCLUDE" |
			rg -v "$LICENSE_FILE_EXCLUDE"
	)
	;;
"$STYLE_SCOPE_ALL")
	mapfile -t FILES < <(
		git -C "$PROJECT_ROOT" ls-files |
			rg -v "$ALL_SCOPE_EXCLUDE" |
			rg -v "$GO_FILE_EXCLUDE" |
			rg -v "$LICENSE_FILE_EXCLUDE"
	)
	;;
esac

if [ "${#FILES[@]}" -eq 0 ]; then
	exit 0
fi

if ! output=$(misspell -error -locale "$MISSPELL_LOCALE" "${FILES[@]/#/$PROJECT_ROOT/}" 2>&1); then
	echo "$RULE_MESSAGE_PREFIX"
	echo "$output"
	echo ""
	exit 1
fi

exit 0
