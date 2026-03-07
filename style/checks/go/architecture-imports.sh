#!/bin/bash
#
# tools/style/checks/go/architecture-imports.sh
# Checks repository dependency boundaries (STYLE.md 2.1).
#
# Rules:
#	- core stays independent from client and relay application code.
#	- application ports depend only on core and sibling ports.
#	- application services depend only on core and sibling application packages.
#	- adapters stay inside their own application boundary.
#	- bootstrap packages compose only their own application.
#	- cmd packages import only their matching internal entrypoint/bootstrap package.
#
# Exit code: 0 if no violations, 1 if violations found.

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
FOUND=0
RULE_LABEL="1.3"
GO_LIST_IMPORT_TEMPLATE='{{join .Imports "\n"}}'
MESSAGE_VIOLATIONS="[${RULE_LABEL}] forbidden internal dependency imports:"

# ---------------------------------------------- Paths ---------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
STYLE_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
# shellcheck source=tools/style/internal/common.sh
source "$STYLE_DIR/internal/common.sh"

PROJECT_ROOT="$(style_project_root_from_dir "$STYLE_DIR")"

if ! style_require_command "go" "$USAGE_EXIT_CODE"; then
	exit "$USAGE_EXIT_CODE"
fi

MODULE=$(cd "$PROJECT_ROOT" && go list -m -f '{{.Path}}')

# -------------------------------------------- Helpers --------------------------------------------

classify_package() {
	case "$1" in
		"$MODULE"/internal/core | "$MODULE"/internal/core/*)
			echo "core"
			;;
		"$MODULE"/internal/client/application/port | "$MODULE"/internal/client/application/port/*)
			echo "client_port"
			;;
		"$MODULE"/internal/client/application/service | "$MODULE"/internal/client/application/service/*)
			echo "client_service"
			;;
		"$MODULE"/internal/client/adapters/inbound | "$MODULE"/internal/client/adapters/inbound/*)
			echo "client_inbound"
			;;
		"$MODULE"/internal/client/adapters/outbound | "$MODULE"/internal/client/adapters/outbound/*)
			echo "client_outbound"
			;;
		"$MODULE"/internal/client/bootstrap | "$MODULE"/internal/client/bootstrap/*)
			echo "client_bootstrap"
			;;
		"$MODULE"/internal/relay/application/port | "$MODULE"/internal/relay/application/port/*)
			echo "relay_port"
			;;
		"$MODULE"/internal/relay/application/service | "$MODULE"/internal/relay/application/service/*)
			echo "relay_service"
			;;
		"$MODULE"/internal/relay/adapters/inbound | "$MODULE"/internal/relay/adapters/inbound/*)
			echo "relay_inbound"
			;;
		"$MODULE"/internal/relay/adapters/outbound | "$MODULE"/internal/relay/adapters/outbound/*)
			echo "relay_outbound"
			;;
		"$MODULE"/internal/relay/bootstrap | "$MODULE"/internal/relay/bootstrap/*)
			echo "relay_bootstrap"
			;;
		"$MODULE"/internal/relaywire | "$MODULE"/internal/relaywire/*)
			echo "shared"
			;;
		"$MODULE"/cmd/ciphera | "$MODULE"/cmd/ciphera/*)
			echo "cmd_ciphera"
			;;
		"$MODULE"/cmd/relay | "$MODULE"/cmd/relay/*)
			echo "cmd_relay"
			;;
		*)
			echo ""
			;;
	esac
}

is_allowed_import() {
	case "$1" in
		core)
			[ "$2" = "core" ]
			;;
		client_port)
			[ "$2" = "core" ] || [ "$2" = "client_port" ]
			;;
		client_service)
			[ "$2" = "core" ] || [ "$2" = "client_port" ] || [ "$2" = "client_service" ]
			;;
		client_inbound)
			[ "$2" = "core" ] ||
				[ "$2" = "client_port" ] ||
				[ "$2" = "client_service" ] ||
				[ "$2" = "client_inbound" ] ||
				[ "$2" = "client_bootstrap" ] ||
				[ "$2" = "shared" ]
			;;
		client_outbound)
			[ "$2" = "core" ] ||
				[ "$2" = "client_port" ] ||
				[ "$2" = "client_outbound" ] ||
				[ "$2" = "shared" ]
			;;
		client_bootstrap)
			[ "$2" = "core" ] ||
				[ "$2" = "client_port" ] ||
				[ "$2" = "client_service" ] ||
				[ "$2" = "client_inbound" ] ||
				[ "$2" = "client_outbound" ] ||
				[ "$2" = "client_bootstrap" ] ||
				[ "$2" = "shared" ]
			;;
		relay_port)
			[ "$2" = "core" ] || [ "$2" = "relay_port" ]
			;;
		relay_service)
			[ "$2" = "core" ] ||
				[ "$2" = "relay_port" ] ||
				[ "$2" = "relay_service" ]
			;;
		relay_inbound)
			[ "$2" = "core" ] ||
				[ "$2" = "relay_port" ] ||
				[ "$2" = "relay_service" ] ||
				[ "$2" = "relay_inbound" ] ||
				[ "$2" = "relay_bootstrap" ] ||
				[ "$2" = "shared" ]
			;;
		relay_outbound)
			[ "$2" = "core" ] ||
				[ "$2" = "relay_port" ] ||
				[ "$2" = "relay_outbound" ] ||
				[ "$2" = "shared" ]
			;;
		relay_bootstrap)
			[ "$2" = "core" ] ||
				[ "$2" = "relay_port" ] ||
				[ "$2" = "relay_service" ] ||
				[ "$2" = "relay_inbound" ] ||
				[ "$2" = "relay_outbound" ] ||
				[ "$2" = "relay_bootstrap" ] ||
				[ "$2" = "shared" ]
			;;
		shared)
			[ "$2" = "core" ] || [ "$2" = "shared" ]
			;;
		cmd_ciphera)
			[ "$2" = "client_inbound" ]
			;;
		cmd_relay)
			[ "$2" = "relay_bootstrap" ]
			;;
		*)
			return 0
			;;
	esac
}

# --------------------------------------------- Checks ---------------------------------------------

violations=$(
	cd "$PROJECT_ROOT"
	while IFS= read -r package; do
		from_category="$(classify_package "$package")"
		if [ -z "$from_category" ]; then
			continue
		fi

		while IFS= read -r import_path; do
			if [ -z "$import_path" ]; then
				continue
			fi

			case "$import_path" in
				"$MODULE"/*)
					;;
				*)
					continue
					;;
			esac

			to_category="$(classify_package "$import_path")"
			if [ -z "$to_category" ]; then
				continue
			fi

			if ! is_allowed_import "$from_category" "$to_category"; then
				printf "%s [%s] imports %s [%s]\n" \
					"$package" \
					"$from_category" \
					"$import_path" \
					"$to_category"
			fi
		done < <(go list -f "$GO_LIST_IMPORT_TEMPLATE" "$package")
	done < <(go list ./...)
)

if [ -n "$violations" ]; then
	echo "$MESSAGE_VIOLATIONS"
	echo "$violations"
	echo ""
	FOUND=1
fi

# ------------------------------------------- Validation -------------------------------------------

if [ "$FOUND" -eq 1 ]; then
	exit 1
fi

exit 0
