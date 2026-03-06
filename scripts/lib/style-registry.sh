#!/bin/bash
#
# tools/scripts/lib/style-registry.sh
# Loads style-check registry entries from a table file.

set -euo pipefail

# ---------------------------------------------- Paths ---------------------------------------------

STYLE_REGISTRY_LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tools/scripts/lib/style-registry-constants.sh
source "$STYLE_REGISTRY_LIB_DIR/style-registry-constants.sh"

# --------------------------------------------- Config ---------------------------------------------

REGISTRY_FIELD_SEPARATOR='|'
REGISTRY_COMMENT_PREFIX="#"
REGISTRY_LEVEL_REQUIRED="${LEVEL_REQUIRED:-required}"
REGISTRY_LEVEL_RECOMMENDATION="${LEVEL_RECOMMENDATION:-recommendation}"
REGISTRY_RULE_RECOMMENDATION_PREFIX="R"
REGISTRY_SCOPE_APP="${STYLE_SCOPE_APP:-app}"
REGISTRY_SCOPE_TOOLS="${STYLE_SCOPE_TOOLS:-tools}"
REGISTRY_SCOPE_ALL="${STYLE_SCOPE_ALL:-all}"
REGISTRY_RUNNER_SCRIPT="${STYLE_RUNNER_SCRIPT:-script}"
REGISTRY_RUNNER_SCRIPT_SCOPE="${STYLE_RUNNER_SCRIPT_SCOPE:-script_scope}"
REGISTRY_RUNNER_EXECUTOR="${STYLE_RUNNER_EXECUTOR:-runner}"
REGISTRY_RUNNER_TARGET_GOLANGCI_APP="${STYLE_RUNNER_TARGET_GOLANGCI_APP:-golangci_app}"
REGISTRY_RUNNER_TARGET_GOLANGCI_TOOLS="${STYLE_RUNNER_TARGET_GOLANGCI_TOOLS:-golangci_tools}"
REGISTRY_RUNNER_TARGET_AST_APP="${STYLE_RUNNER_TARGET_AST_APP:-ast_app}"
REGISTRY_RUNNER_TARGET_AST_TOOLS="${STYLE_RUNNER_TARGET_AST_TOOLS:-ast_tools}"
REGISTRY_TIER_ONE="${STYLE_TIER_ONE:-tier1}"
REGISTRY_TIER_TWO="${STYLE_TIER_TWO:-tier2}"
REGISTRY_TIER_THREE="${STYLE_TIER_THREE:-tier3}"
REGISTRY_SCRIPT_TARGET_SUFFIX=".sh"
REGISTRY_MESSAGE_INCONSISTENT="style check registry is inconsistent"
REGISTRY_MESSAGE_TABLE_MISSING="style check registry table is not readable:"
REGISTRY_MESSAGE_TABLE_ROW_INVALID="invalid style check registry row:"
REGISTRY_TABLE_FILENAME="style-registry.table"
REGISTRY_MIN_CHECK_COUNT=1
REGISTRY_ROW_KEY_SEPARATOR=$'\x1f'

declare -a CHECK_TIERS=()
declare -a CHECK_LEVELS=()
declare -a CHECK_RULES=()
declare -a CHECK_NAMES=()
declare -a CHECK_SCOPES=()
declare -a CHECK_RUNNERS=()
declare -a CHECK_TARGETS=()

# ---------------------------------------------- Paths ---------------------------------------------

STYLE_REGISTRY_SCRIPTS_DIR="$(cd "$STYLE_REGISTRY_LIB_DIR/.." && pwd)"
DEFAULT_TABLE_FILE="$STYLE_REGISTRY_LIB_DIR/$REGISTRY_TABLE_FILENAME"
STYLE_REGISTRY_TABLE_FILE="${STYLE_REGISTRY_TABLE_FILE:-$DEFAULT_TABLE_FILE}"

# ---------------------------------------- Registry Helpers ----------------------------------------

style_registry_clear() {
	CHECK_TIERS=()
	CHECK_LEVELS=()
	CHECK_RULES=()
	CHECK_NAMES=()
	CHECK_SCOPES=()
	CHECK_RUNNERS=()
	CHECK_TARGETS=()
}

style_registry_add_check() {
	local tier="$1"
	local level="$2"
	local rule="$3"
	local name="$4"
	local required_scope="$5"
	local runner_kind="$6"
	local target="$7"
	local index="${#CHECK_TIERS[@]}"

	CHECK_TIERS[index]="$tier"
	CHECK_LEVELS[index]="$level"
	CHECK_RULES[index]="$rule"
	CHECK_NAMES[index]="$name"
	CHECK_SCOPES[index]="$required_scope"
	CHECK_RUNNERS[index]="$runner_kind"
	CHECK_TARGETS[index]="$target"
}

style_trim_whitespace() {
	local value="$1"
	value="${value#"${value%%[![:space:]]*}"}"
	value="${value%"${value##*[![:space:]]}"}"
	printf '%s' "$value"
}

style_registry_is_valid_tier() {
	local tier="$1"
	case "$tier" in
	"$REGISTRY_TIER_ONE" | "$REGISTRY_TIER_TWO" | "$REGISTRY_TIER_THREE")
		return 0
		;;
	*)
		return 1
		;;
	esac
}

style_registry_is_valid_level() {
	local level="$1"
	case "$level" in
	"$REGISTRY_LEVEL_REQUIRED" | "$REGISTRY_LEVEL_RECOMMENDATION")
		return 0
		;;
	*)
		return 1
		;;
	esac
}

style_registry_is_valid_rule() {
	local level="$1"
	local rule="$2"

	case "$level" in
	"$REGISTRY_LEVEL_REQUIRED")
		[[ "$rule" != "$REGISTRY_RULE_RECOMMENDATION_PREFIX"* ]]
		return
		;;
	"$REGISTRY_LEVEL_RECOMMENDATION")
		[[ "$rule" == "$REGISTRY_RULE_RECOMMENDATION_PREFIX"* ]]
		return
		;;
	*)
		return 1
		;;
	esac
}

style_registry_is_valid_scope() {
	local required_scope="$1"
	case "$required_scope" in
	"$REGISTRY_SCOPE_APP" | "$REGISTRY_SCOPE_TOOLS" | "$REGISTRY_SCOPE_ALL")
		return 0
		;;
	*)
		return 1
		;;
	esac
}

style_registry_is_valid_runner() {
	local runner_kind="$1"
	case "$runner_kind" in
	"$REGISTRY_RUNNER_SCRIPT" | "$REGISTRY_RUNNER_SCRIPT_SCOPE" | "$REGISTRY_RUNNER_EXECUTOR")
		return 0
		;;
	*)
		return 1
		;;
	esac
}

style_registry_is_valid_executor_target() {
	local target="$1"
	case "$target" in
	"$REGISTRY_RUNNER_TARGET_GOLANGCI_APP" | "$REGISTRY_RUNNER_TARGET_GOLANGCI_TOOLS" | \
		"$REGISTRY_RUNNER_TARGET_AST_APP" | "$REGISTRY_RUNNER_TARGET_AST_TOOLS")
		return 0
		;;
	*)
		return 1
		;;
	esac
}

style_registry_is_valid_target() {
	local runner_kind="$1"
	local target="$2"
	case "$runner_kind" in
	"$REGISTRY_RUNNER_SCRIPT" | "$REGISTRY_RUNNER_SCRIPT_SCOPE")
		[[ "$target" == *"$REGISTRY_SCRIPT_TARGET_SUFFIX" ]] || return 1
		[ -f "$STYLE_REGISTRY_SCRIPTS_DIR/$target" ]
		return
		;;
	"$REGISTRY_RUNNER_EXECUTOR")
		style_registry_is_valid_executor_target "$target"
		return
		;;
	*)
		return 1
		;;
	esac
}

style_registry_is_valid_row() {
	local tier="$1"
	local level="$2"
	local rule="$3"
	local required_scope="$4"
	local runner_kind="$5"
	local target="$6"

	if ! style_registry_is_valid_tier "$tier"; then
		return 1
	fi

	if ! style_registry_is_valid_level "$level"; then
		return 1
	fi

	if ! style_registry_is_valid_rule "$level" "$rule"; then
		return 1
	fi

	if ! style_registry_is_valid_scope "$required_scope"; then
		return 1
	fi

	if ! style_registry_is_valid_runner "$runner_kind"; then
		return 1
	fi

	if ! style_registry_is_valid_target "$runner_kind" "$target"; then
		return 1
	fi

	return 0
}

style_registry_add_table_row() {
	local row="$1"
	local tier=""
	local level=""
	local rule=""
	local name=""
	local required_scope=""
	local runner_kind=""
	local target=""
	local extra=""

	IFS="$REGISTRY_FIELD_SEPARATOR" read -r \
		tier level rule name required_scope runner_kind target extra <<<"$row"

	tier="$(style_trim_whitespace "$tier")"
	level="$(style_trim_whitespace "$level")"
	rule="$(style_trim_whitespace "$rule")"
	name="$(style_trim_whitespace "$name")"
	required_scope="$(style_trim_whitespace "$required_scope")"
	runner_kind="$(style_trim_whitespace "$runner_kind")"
	target="$(style_trim_whitespace "$target")"
	extra="$(style_trim_whitespace "$extra")"

	if [ -n "$extra" ] || [ -z "$tier" ] || [ -z "$level" ] || [ -z "$rule" ] ||
		[ -z "$name" ] || [ -z "$required_scope" ] || [ -z "$runner_kind" ] ||
		[ -z "$target" ]; then
		return 1
	fi

	if ! style_registry_is_valid_row \
		"$tier" "$level" "$rule" "$required_scope" "$runner_kind" "$target"; then
		return 1
	fi

	style_registry_add_check \
		"$tier" \
		"$level" \
		"$rule" \
		"$name" \
		"$required_scope" \
		"$runner_kind" \
		"$target"
}

style_registry_load_table() {
	local table_file="$1"
	local row_number=0
	local row=""
	local trimmed_row=""

	# shellcheck disable=SC2094  # Reads from "$table_file" only; no writes occur in this loop.
	while IFS= read -r row || [ -n "$row" ]; do
		row_number=$((row_number + 1))

		if [[ ! "$row" =~ [^[:space:]] ]]; then
			continue
		fi

		trimmed_row="${row#"${row%%[![:space:]]*}"}"
		if [[ "$trimmed_row" == "$REGISTRY_COMMENT_PREFIX"* ]]; then
			continue
		fi

		if ! style_registry_add_table_row "$row"; then
			echo "$REGISTRY_MESSAGE_TABLE_ROW_INVALID $table_file:$row_number"
			return 1
		fi
	done <"$table_file"
}

# ------------------------------------------- Validation -------------------------------------------

style_registry_validate() {
	local expected_count="${#CHECK_TIERS[@]}"
	local index
	local row_key

	declare -A seen_row_keys=()

	if [ "$expected_count" -lt "$REGISTRY_MIN_CHECK_COUNT" ]; then
		return 1
	fi

	[ "${#CHECK_LEVELS[@]}" -eq "$expected_count" ] || return 1
	[ "${#CHECK_RULES[@]}" -eq "$expected_count" ] || return 1
	[ "${#CHECK_NAMES[@]}" -eq "$expected_count" ] || return 1
	[ "${#CHECK_SCOPES[@]}" -eq "$expected_count" ] || return 1
	[ "${#CHECK_RUNNERS[@]}" -eq "$expected_count" ] || return 1
	[ "${#CHECK_TARGETS[@]}" -eq "$expected_count" ] || return 1

	for index in "${!CHECK_TIERS[@]}"; do
		row_key="${CHECK_TIERS[$index]}$REGISTRY_ROW_KEY_SEPARATOR"
		row_key+="${CHECK_LEVELS[$index]}$REGISTRY_ROW_KEY_SEPARATOR"
		row_key+="${CHECK_RULES[$index]}$REGISTRY_ROW_KEY_SEPARATOR"
		row_key+="${CHECK_NAMES[$index]}$REGISTRY_ROW_KEY_SEPARATOR"
		row_key+="${CHECK_SCOPES[$index]}$REGISTRY_ROW_KEY_SEPARATOR"
		row_key+="${CHECK_RUNNERS[$index]}$REGISTRY_ROW_KEY_SEPARATOR"
		row_key+="${CHECK_TARGETS[$index]}"

		if [[ -n "${seen_row_keys[$row_key]+x}" ]]; then
			return 1
		fi

		seen_row_keys[$row_key]=1
	done

	return 0
}

# -------------------------------------------- Registry --------------------------------------------

style_register_default_checks() {
	style_registry_clear

	if [ ! -r "$STYLE_REGISTRY_TABLE_FILE" ]; then
		echo "$REGISTRY_MESSAGE_TABLE_MISSING $STYLE_REGISTRY_TABLE_FILE"
		return 1
	fi

	if ! style_registry_load_table "$STYLE_REGISTRY_TABLE_FILE"; then
		return 1
	fi

	if ! style_registry_validate; then
		echo "$REGISTRY_MESSAGE_INCONSISTENT"
		return 1
	fi
}
