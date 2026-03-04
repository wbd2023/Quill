#!/bin/bash
#
# tools/scripts/check-go-error-style.sh
# Checks error handling conventions (STYLE.md 2.1).
#
# Rules:
#   - Error context strings must be lowercase.
#   - Error context strings must not end with punctuation.
#   - Error context strings must not include secrets.
#   - Sentinel errors (var Err...) live in domain/errors.go only.
#
# Exit code: 0 if no violations, 1 if violations found.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
FOUND=0
RULE_LABEL="2.1"
GO_TEST_EXCLUDE='!*_test.go'
ERRORF_UPPERCASE_REGEX='fmt\.Errorf\("[A-Z]'
ERRORS_NEW_UPPERCASE_REGEX='errors\.New\("[A-Z]'
ERRORS_PUNCTUATION_REGEX='(fmt\.Errorf|errors\.New)\("[^"]*[.!?]"\)'
SENTINEL_VAR_REGEX='^\s*var\s+Err\w+\s*='
SENTINEL_ASSIGN_REGEX='^\s*Err\w+\s*='
SENTINEL_ALLOWED_PATH='domain/errors.go'
VENDOR_FILTER='vendor'
MESSAGE_ERRORF_CASE="[${RULE_LABEL}] Error context must be lowercase (fmt.Errorf):"
MESSAGE_ERRORS_NEW_CASE="[${RULE_LABEL}] Error context must be lowercase (errors.New):"
MESSAGE_PUNCTUATION="[${RULE_LABEL}] Error context must not end with punctuation:"
MESSAGE_SECRETS="[${RULE_LABEL}] Error context must not include secrets:"
MESSAGE_SENTINELS="[${RULE_LABEL}] Sentinel errors must live in domain/errors.go:"

# ---------------------------------------------- Paths ---------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tools/scripts/lib/style-common.sh
source "$SCRIPT_DIR/lib/style-common.sh"

PROJECT_ROOT="$(style_project_root_from_dir "$SCRIPT_DIR")"

if ! style_require_command "rg" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

# --------------------------------------------- Checks ---------------------------------------------

# fmt.Errorf starting with uppercase.
uppercase_errorf=$(rg -n --glob "$STYLE_PATTERN_GO" "$ERRORF_UPPERCASE_REGEX" \
	"$PROJECT_ROOT/$STYLE_PATH_INTERNAL" "$PROJECT_ROOT/$STYLE_PATH_CMD" |
	grep -v "$VENDOR_FILTER" || true)

if [ -n "$uppercase_errorf" ]; then
	echo "$MESSAGE_ERRORF_CASE"
	echo "$uppercase_errorf"
	echo ""
	FOUND=1
fi

# errors.New starting with uppercase.
uppercase_new=$(rg -n --glob "$STYLE_PATTERN_GO" "$ERRORS_NEW_UPPERCASE_REGEX" \
	"$PROJECT_ROOT/$STYLE_PATH_INTERNAL" "$PROJECT_ROOT/$STYLE_PATH_CMD" |
	grep -v "$VENDOR_FILTER" || true)

if [ -n "$uppercase_new" ]; then
	echo "$MESSAGE_ERRORS_NEW_CASE"
	echo "$uppercase_new"
	echo ""
	FOUND=1
fi

# Error strings ending with punctuation.
punctuation=$(rg -n --glob "$STYLE_PATTERN_GO" \
	"$ERRORS_PUNCTUATION_REGEX" \
	"$PROJECT_ROOT/$STYLE_PATH_INTERNAL" "$PROJECT_ROOT/$STYLE_PATH_CMD" |
	grep -v "$VENDOR_FILTER" || true)

if [ -n "$punctuation" ]; then
	echo "$MESSAGE_PUNCTUATION"
	echo "$punctuation"
	echo ""
	FOUND=1
fi

# Secret-bearing values passed into fmt.Errorf formatting args.
#
# This catches likely leaks such as:
#   fmt.Errorf("failed auth for %s", passphrase)
# It does not flag plain wording like "passphrase is required".
secrets_pattern='fmt\.Errorf\("[^"]*%[^"]*",[^)]*\b(passphrase|password|privateKey|'
secrets_pattern+='secretKey|secret|token|seed)\b'

secrets_in_errors=$(rg -n --glob "$STYLE_PATTERN_GO" --glob "$GO_TEST_EXCLUDE" \
	"$secrets_pattern" \
	"$PROJECT_ROOT/$STYLE_PATH_INTERNAL" "$PROJECT_ROOT/$STYLE_PATH_CMD" |
	grep -v "$VENDOR_FILTER" || true)

if [ -n "$secrets_in_errors" ]; then
	echo "$MESSAGE_SECRETS"
	echo "$secrets_in_errors"
	echo ""
	FOUND=1
fi

# Sentinel errors outside domain/errors.go.
# Covers both:
#   var ErrX = ...
#   var ( ErrX = ... )
sentinel_outside=$(
	{
		rg -n --glob "$STYLE_PATTERN_GO" --glob "$GO_TEST_EXCLUDE" "$SENTINEL_VAR_REGEX" \
			"$PROJECT_ROOT/$STYLE_PATH_INTERNAL" || true
		rg -n --glob "$STYLE_PATTERN_GO" --glob "$GO_TEST_EXCLUDE" "$SENTINEL_ASSIGN_REGEX" \
			"$PROJECT_ROOT/$STYLE_PATH_INTERNAL" || true
	} | grep -v "$VENDOR_FILTER" | grep -v "$SENTINEL_ALLOWED_PATH" | sort -u || true
)

if [ -n "$sentinel_outside" ]; then
	echo "$MESSAGE_SENTINELS"
	echo "$sentinel_outside"
	echo ""
	FOUND=1
fi

# ------------------------------------------- Validation -------------------------------------------

if [ "$FOUND" -eq 1 ]; then
	exit 1
fi

exit 0
