#!/bin/bash
#
# tools/scripts/lib/style-common.sh
# Shared helpers for STYLE.md shell check scripts.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

# Shared constant values for style scripts.
STYLE_SCOPE_APP="app"
STYLE_SCOPE_TOOLS="tools"
STYLE_SCOPE_ALL="all"
STYLE_SCOPE_USAGE="app|tools|all"

STYLE_PATH_CMD="cmd"
STYLE_PATH_INTERNAL="internal"
STYLE_PATH_TOOLS="tools"
STYLE_PATH_GIT=".git"
STYLE_PATH_VENDOR="vendor"
STYLE_PATH_BIN="bin"

STYLE_PATTERN_GO="*.go"
STYLE_PATTERN_SHELL="*.sh"
STYLE_PATTERN_MARKDOWN="*.md"
STYLE_MESSAGE_UNSUPPORTED_PATTERN="unsupported file pattern"

# --------------------------------------------- Helpers --------------------------------------------

style_project_root_from_dir() {
	local script_dir="$1"
	(cd "$script_dir/../.." && pwd)
}

# Returns success if a file pattern is one of the shared supported patterns.
style_is_supported_pattern() {
	local pattern="$1"

	case "$pattern" in
	"$STYLE_PATTERN_GO" | "$STYLE_PATTERN_SHELL" | "$STYLE_PATTERN_MARKDOWN")
		return 0
		;;
	*)
		return 1
		;;
	esac
}

style_parse_scope_arg() {
	local usage_exit_code="${1:-2}"
	local default_scope="${2:-$STYLE_SCOPE_ALL}"
	shift 2

	if [ "$#" -eq 0 ]; then
		printf '%s\n' "$default_scope"
		return 0
	fi

	if [ "$#" -eq 2 ] && [ "$1" = "--scope" ]; then
		printf '%s\n' "$2"
		return 0
	fi

	echo "usage: $0 [--scope $STYLE_SCOPE_USAGE]" >&2
	return "$usage_exit_code"
}

style_validate_scope() {
	local scope="$1"
	local usage_exit_code="${2:-2}"

	case "$scope" in
	"$STYLE_SCOPE_APP" | "$STYLE_SCOPE_TOOLS" | "$STYLE_SCOPE_ALL")
		return 0
		;;
	*)
		echo "invalid scope: $scope (expected $STYLE_SCOPE_USAGE)"
		return "$usage_exit_code"
		;;
	esac
}

style_collect_files() {
	local project_root="$1"
	local scope="$2"
	local pattern="$3"

	if ! style_is_supported_pattern "$pattern"; then
		echo "$STYLE_MESSAGE_UNSUPPORTED_PATTERN: $pattern" >&2
		return 1
	fi

	case "$scope" in
	"$STYLE_SCOPE_APP")
		find "$project_root/$STYLE_PATH_INTERNAL" "$project_root/$STYLE_PATH_CMD" \
			-type f -name "$pattern" 2>/dev/null | sort -u
		;;
	"$STYLE_SCOPE_TOOLS")
		find "$project_root/$STYLE_PATH_TOOLS" \
			-type f -name "$pattern" 2>/dev/null | sort -u
		;;
	"$STYLE_SCOPE_ALL")
		find "$project_root" \
			-type d \( \
			-name "$STYLE_PATH_GIT" -o \
			-name "$STYLE_PATH_VENDOR" -o \
			-name "$STYLE_PATH_BIN" \
			\) -prune -o \
			-type f -name "$pattern" -print 2>/dev/null | sort -u
		;;
	esac
}

style_require_command() {
	local command_name="$1"
	local usage_exit_code="${2:-2}"

	if command -v "$command_name" >/dev/null 2>&1; then
		return 0
	fi

	echo "$command_name is not installed"
	return "$usage_exit_code"
}
