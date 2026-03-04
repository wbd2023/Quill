#!/bin/bash
#
# tools/scripts/check-section-headers.sh
# Checks section header formatting (STYLE.md 2.4).
#
# Rules:
#	- Section headers must be exactly 100 characters.
#	- Section header text must be centred.
#	- If centring is uneven by one, the left dash padding must be longer.
#	- 100+ line Go/Bash files must include at least one section header.
#
# Usage:
#	./tools/scripts/check-section-headers.sh [--scope app|tools|all]
#
# Exit code: 0 if no violations, 1 if violations found.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
SECTION_HEADER_LENGTH=100
MIN_LINES_FOR_HEADERS=100
RULE_LABEL="2.4"
GO_HEADER_REGEX='^/\*\ -+\ .+\ -+\ \*/$'
SHELL_HEADER_REGEX='^#\ -+\ .+\ -+$'
HEADER_BODY_REGEX='^(-+)\ (.+)\ (-+)$'
MESSAGE_LENGTH="[${RULE_LABEL}] Section header not ${SECTION_HEADER_LENGTH} chars"
MESSAGE_CENTRE="[${RULE_LABEL}] Section header text is not centred with left-side precedence:"
MESSAGE_MISSING="[${RULE_LABEL}] Missing section headers in ${MIN_LINES_FOR_HEADERS}+ line file:"
FOUND=0

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

declare -a FILES=()
case "$SCOPE" in
"$STYLE_SCOPE_APP")
	mapfile -t FILES < <(style_collect_files "$PROJECT_ROOT" "$STYLE_SCOPE_APP" "$STYLE_PATTERN_GO")
	;;
"$STYLE_SCOPE_TOOLS")
	mapfile -t FILES < <(
		{
			style_collect_files "$PROJECT_ROOT" "$STYLE_SCOPE_TOOLS" "$STYLE_PATTERN_GO"
			style_collect_files "$PROJECT_ROOT" "$STYLE_SCOPE_TOOLS" "$STYLE_PATTERN_SHELL"
		} | sort -u
	)
	;;
"$STYLE_SCOPE_ALL")
	mapfile -t FILES < <(
		{
			style_collect_files "$PROJECT_ROOT" "$STYLE_SCOPE_APP" "$STYLE_PATTERN_GO"
			style_collect_files "$PROJECT_ROOT" "$STYLE_SCOPE_APP" "$STYLE_PATTERN_SHELL"
			style_collect_files "$PROJECT_ROOT" "$STYLE_SCOPE_TOOLS" "$STYLE_PATTERN_GO"
			style_collect_files "$PROJECT_ROOT" "$STYLE_SCOPE_TOOLS" "$STYLE_PATTERN_SHELL"
		} | sort -u
	)
	;;
esac

# --------------------------------------------- Helpers --------------------------------------------

check_centred() {
	local body="$1"
	if [[ ! "$body" =~ $HEADER_BODY_REGEX ]]; then
		return 1
	fi

	local left="${#BASH_REMATCH[1]}"
	local right="${#BASH_REMATCH[3]}"
	[ "$left" -eq "$right" ] || [ "$left" -eq $((right + 1)) ]
}

for file in "${FILES[@]}"; do
	file_lines=$(wc -l <"$file")
	header_count=0
	linenum=0

	# shellcheck disable=SC2094  # Reads from "$file" only; no writes occur in this loop.
	while IFS= read -r line; do
		linenum=$((linenum + 1))
		is_header=false
		body=""

		case "$file" in
		*.go)
			if [[ "$line" =~ $GO_HEADER_REGEX ]]; then
				is_header=true
				body="${line#/\* }"
				body="${body% \*/}"
			fi
			;;
		*.sh)
			if [[ "$line" =~ $SHELL_HEADER_REGEX ]]; then
				is_header=true
				body="${line#\# }"
			fi
			;;
		esac

		if [ "$is_header" = false ]; then
			continue
		fi

		header_count=$((header_count + 1))
		length=${#line}
		if [ "$length" -ne "$SECTION_HEADER_LENGTH" ]; then
			echo "${MESSAGE_LENGTH} (got $length):"
			echo "  $file:$linenum"
			echo ""
			FOUND=1
		fi

		if ! check_centred "$body"; then
			echo "$MESSAGE_CENTRE"
			echo "  $file:$linenum"
			echo ""
			FOUND=1
		fi
	done <"$file"

	if [ "$file_lines" -ge "$MIN_LINES_FOR_HEADERS" ] && [ "$header_count" -eq 0 ]; then
		echo "$MESSAGE_MISSING"
		echo "  $file"
		echo ""
		FOUND=1
	fi
done

# ------------------------------------------- Validation -------------------------------------------

if [ "$FOUND" -eq 1 ]; then
	exit 1
fi

exit 0
