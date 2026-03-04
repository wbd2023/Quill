#!/bin/bash
#
# tools/scripts/check-bash-style.sh
# Checks Bash-specific style rules (STYLE.md 2.11).
#
# Rules:
#	- Scripts start with #!/bin/bash.
#	- Scripts include set -euo pipefail.
#	- Leading indentation uses tabs, not spaces.
#	- No trailing whitespace.
#	- Unix line endings only.
#
# Usage:
#	./tools/scripts/check-bash-style.sh [--scope app|tools|all]
#
# Exit code: 0 if no violations, 1 if violations found.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
FOUND=0
SHEBANG_LINE="#!/bin/bash"
STRICT_MODE_LINE='set -euo pipefail'
STRICT_MODE_REGEX='^set -euo pipefail$'
SPACE_INDENT_REGEX='^[ ]+\S'
TRAILING_WHITESPACE_REGEX='[ \t]+$'
CRLF_REGEX='\r$'
RULE_PREFIX="[2.11]"
MESSAGE_SHEBANG="${RULE_PREFIX} Bash script must start with '${SHEBANG_LINE}':"
MESSAGE_STRICT_MODE="${RULE_PREFIX} Bash script must include '${STRICT_MODE_LINE}':"
MESSAGE_SPACE_INDENT="${RULE_PREFIX} Leading indentation must use tabs, not spaces:"
MESSAGE_TRAILING_WHITESPACE="${RULE_PREFIX} Trailing whitespace detected:"
MESSAGE_CRLF="${RULE_PREFIX} CRLF line endings detected:"

# ---------------------------------------------- Paths ---------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tools/scripts/lib/style-common.sh
source "$SCRIPT_DIR/lib/style-common.sh"

PROJECT_ROOT="$(style_project_root_from_dir "$SCRIPT_DIR")"

# ---------------------------------------------- Args ----------------------------------------------

if ! SCOPE="$(style_parse_scope_arg "$USAGE_EXIT_CODE" "$STYLE_SCOPE_ALL" "$@")"; then
	exit "$USAGE_EXIT_CODE"
fi

if ! style_validate_scope "$SCOPE" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

mapfile -t SCRIPT_FILES < <(style_collect_files "$PROJECT_ROOT" "$SCOPE" "$STYLE_PATTERN_SHELL")

if [ "${#SCRIPT_FILES[@]}" -eq 0 ]; then
	exit 0
fi

if ! style_require_command "rg" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

# --------------------------------------------- Checks ---------------------------------------------

for file in "${SCRIPT_FILES[@]}"; do

	first_line="$(head -n 1 "$file")"
	if [ "$first_line" != "$SHEBANG_LINE" ]; then
		echo "$MESSAGE_SHEBANG"
		echo "  $file:1"
		echo ""
		FOUND=1
	fi

	if ! rg -n "$STRICT_MODE_REGEX" "$file" >/dev/null 2>&1; then
		echo "$MESSAGE_STRICT_MODE"
		echo "  $file"
		echo ""
		FOUND=1
	fi
done

space_indent=$(rg -n "$SPACE_INDENT_REGEX" "${SCRIPT_FILES[@]}" || true)
if [ -n "$space_indent" ]; then
	echo "$MESSAGE_SPACE_INDENT"
	echo "$space_indent"
	echo ""
	FOUND=1
fi

trailing_ws=$(rg -n "$TRAILING_WHITESPACE_REGEX" "${SCRIPT_FILES[@]}" || true)
if [ -n "$trailing_ws" ]; then
	echo "$MESSAGE_TRAILING_WHITESPACE"
	echo "$trailing_ws"
	echo ""
	FOUND=1
fi

crlf_lines=$(rg -n "$CRLF_REGEX" "${SCRIPT_FILES[@]}" || true)
if [ -n "$crlf_lines" ]; then
	echo "$MESSAGE_CRLF"
	echo "$crlf_lines"
	echo ""
	FOUND=1
fi

# ------------------------------------------- Validation -------------------------------------------

if [ "$FOUND" -eq 1 ]; then
	exit 1
fi

exit 0
