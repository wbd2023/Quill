#!/bin/bash
#
# tools/scripts/check-go-architecture-imports.sh
# Checks Clean Architecture dependency rules (STYLE.md 1.3).
#
# Rules:
#   - core must not import from adapters or app.
#   - adapters must not import from app.
#
# Exit code: 0 if no violations, 1 if violations found.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
FOUND=0
RULE_LABEL="1.3"
GO_LIST_IMPORT_TEMPLATE='{{join .Imports "\n"}}'
CORE_PACKAGES_QUERY="./internal/core/..."
ADAPTERS_PACKAGES_QUERY="./internal/adapters/..."
ADAPTERS_PATH_SUFFIX="/internal/adapters"
APP_PATH_SUFFIX="/internal/app"
MESSAGE_CORE_IMPORTS="[${RULE_LABEL}] core must not import from adapters or app:"
MESSAGE_ADAPTER_IMPORTS="[${RULE_LABEL}] adapters must not import from app:"

# ---------------------------------------------- Paths ---------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tools/scripts/lib/style-common.sh
source "$SCRIPT_DIR/lib/style-common.sh"

PROJECT_ROOT="$(style_project_root_from_dir "$SCRIPT_DIR")"

if ! style_require_command "go" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

MODULE=$(cd "$PROJECT_ROOT" && go list -m -f '{{.Path}}')

# --------------------------------------------- Checks ---------------------------------------------

# core -> adapters or app.
core_violations=$(
	cd "$PROJECT_ROOT"
	for package in $(go list "$CORE_PACKAGES_QUERY"); do
		go list -f "$GO_LIST_IMPORT_TEMPLATE" "$package" |
			awk -v package_name="$package" \
				-v module_path="$MODULE" \
				-v adapters_suffix="$ADAPTERS_PATH_SUFFIX" \
				-v app_suffix="$APP_PATH_SUFFIX" '
					$0 ~ "^" module_path adapters_suffix || $0 ~ "^" module_path app_suffix {
						printf("%s imports %s\n", package_name, $0)
					}
				'
	done
)

if [ -n "$core_violations" ]; then
	echo "$MESSAGE_CORE_IMPORTS"
	echo "$core_violations"
	echo ""
	FOUND=1
fi

# adapters -> app.
adapters_violations=$(
	cd "$PROJECT_ROOT"
	for package in $(go list "$ADAPTERS_PACKAGES_QUERY"); do
		go list -f "$GO_LIST_IMPORT_TEMPLATE" "$package" |
			awk -v package_name="$package" \
				-v module_path="$MODULE" \
				-v app_suffix="$APP_PATH_SUFFIX" '
					$0 ~ "^" module_path app_suffix {
						printf("%s imports %s\n", package_name, $0)
					}
				'
	done
)

if [ -n "$adapters_violations" ]; then
	echo "$MESSAGE_ADAPTER_IMPORTS"
	echo "$adapters_violations"
	echo ""
	FOUND=1
fi

# ------------------------------------------- Validation -------------------------------------------

if [ "$FOUND" -eq 1 ]; then
	exit 1
fi

exit 0
