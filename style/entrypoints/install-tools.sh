#!/bin/bash
#
# tools/style/entrypoints/install-tools.sh
# Installs required third-party tools for STYLE.md checks.
#
# Usage:
#	./tools/style/entrypoints/install-tools.sh

set -euo pipefail

# --------------------------------------------- Config ---------------------------------------------

USAGE_EXIT_CODE=2
SHELLCHECK_VERSION="0.10.0"
MISSPELL_VERSION="v0.3.4"
GOLANGCI_LINT_VERSION="v2.6.2"
SHFMT_VERSION="v3.12.0"
MARKDOWNLINT_CLI_VERSION="0.45.0"
GO_BIN_SUBDIR="/bin"
LOCAL_NPM_PREFIX="${HOME}/.local"
LOCAL_NPM_BIN="$LOCAL_NPM_PREFIX/bin"
MESSAGE_GO_REQUIRED="go is required to install style tools"
MESSAGE_PATH_HINT="ensure \$(go env GOPATH)/bin and \$HOME/.local/bin are on PATH"

# ---------------------------------------------- Paths ---------------------------------------------

if ! command -v go >/dev/null 2>&1; then
	echo "$MESSAGE_GO_REQUIRED"
	exit "$USAGE_EXIT_CODE"
fi

GO_BIN_DIR="$(go env GOPATH)$GO_BIN_SUBDIR"
mkdir -p "$GO_BIN_DIR"

if [ -d "$LOCAL_NPM_BIN" ]; then
	PATH="$LOCAL_NPM_BIN:$PATH"
fi
PATH="$GO_BIN_DIR:$PATH"

# --------------------------------------------- Helpers --------------------------------------------

install_go_tool() {
	local label="$1"
	local package_path="$2"
	local version="$3"

	echo "Installing $label@$version..."
	go install "$package_path@$version"
}

verify_archive_checksum() {
	local archive_path="$1"
	local archive_name="$2"
	local checksum_file="$3"
	local expected_hash
	local actual_hash

	if ! command -v sha256sum >/dev/null 2>&1 && ! command -v shasum >/dev/null 2>&1; then
		echo "cannot verify checksum: no sha256 tool available"
		return 1
	fi

	expected_hash="$(
		awk -v target="$archive_name" '
			{
				file_name = $2
				sub(/^\*/, "", file_name)
			}
			file_name == target {
				print $1
				exit
			}
		' "$checksum_file"
	)"

	if [ -z "$expected_hash" ]; then
		echo "checksum not found for $archive_name"
		return 1
	fi

	if command -v sha256sum >/dev/null 2>&1; then
		actual_hash="$(sha256sum "$archive_path" | awk '{print $1}')"
	else
		actual_hash="$(shasum -a 256 "$archive_path" | awk '{print $1}')"
	fi

	if [ "$actual_hash" != "$expected_hash" ]; then
		echo "checksum mismatch for $archive_name"
		return 1
	fi

	return 0
}

# ------------------------------------------- Installers -------------------------------------------

install_shellcheck() {
	local os_name
	local architecture
	local asset_name
	local archive_name
	local download_url
	local checksum_url
	local temporary_directory
	local extracted_directory
	local checksum_file

	if command -v shellcheck >/dev/null 2>&1; then
		return
	fi

	if command -v apt-get >/dev/null 2>&1; then
		if [ "$(id -u)" -eq 0 ]; then
			echo "Installing shellcheck via apt-get (root)..."
			apt-get update
			apt-get install -y shellcheck
			return
		fi

		if command -v sudo >/dev/null 2>&1 && sudo -n true >/dev/null 2>&1; then
			echo "Installing shellcheck via apt-get (sudo)..."
			sudo -n apt-get update
			sudo -n apt-get install -y shellcheck
			return
		fi
	fi

	if ! command -v curl >/dev/null 2>&1 || ! command -v tar >/dev/null 2>&1; then
		echo "shellcheck install fallback requires curl and tar"
		return
	fi

	os_name="$(uname -s | tr '[:upper:]' '[:lower:]')"
	architecture="$(uname -m)"

	case "$os_name:$architecture" in
	linux:x86_64)
		asset_name="linux.x86_64"
		;;
	linux:aarch64 | linux:arm64)
		asset_name="linux.aarch64"
		;;
	darwin:x86_64)
		asset_name="darwin.x86_64"
		;;
	darwin:arm64)
		asset_name="darwin.aarch64"
		;;
	*)
		echo "shellcheck unsupported platform: $os_name/$architecture"
		return
		;;
	esac

	archive_name="shellcheck-v${SHELLCHECK_VERSION}.${asset_name}.tar.xz"
	download_url="https://github.com/koalaman/shellcheck/releases/download/"
	download_url+="v${SHELLCHECK_VERSION}/${archive_name}"
	checksum_url="https://github.com/koalaman/shellcheck/releases/download/"
	checksum_url+="v${SHELLCHECK_VERSION}/sha256sums.txt"
	temporary_directory="$(mktemp -d)"
	extracted_directory="$temporary_directory/shellcheck-v${SHELLCHECK_VERSION}"
	checksum_file="$temporary_directory/sha256sums.txt"

	echo "Installing shellcheck from release archive..."
	if ! curl -fsSL "$download_url" -o "$temporary_directory/$archive_name"; then
		echo "failed to download shellcheck archive"
		rm -rf "$temporary_directory"
		return
	fi

	if ! curl -fsSL "$checksum_url" -o "$checksum_file"; then
		echo "failed to download shellcheck checksums"
		rm -rf "$temporary_directory"
		return
	fi

	if ! verify_archive_checksum \
		"$temporary_directory/$archive_name" \
		"$archive_name" \
		"$checksum_file"; then
		rm -rf "$temporary_directory"
		return
	fi

	if ! tar -xJf "$temporary_directory/$archive_name" -C "$temporary_directory"; then
		echo "failed to extract shellcheck archive"
		rm -rf "$temporary_directory"
		return
	fi

	if [ ! -x "$extracted_directory/shellcheck" ]; then
		echo "shellcheck binary missing from extracted archive"
		rm -rf "$temporary_directory"
		return
	fi

	cp "$extracted_directory/shellcheck" "$GO_BIN_DIR/shellcheck"
	chmod 755 "$GO_BIN_DIR/shellcheck"
	rm -rf "$temporary_directory"
}

install_ripgrep() {
	if command -v rg >/dev/null 2>&1; then
		return
	fi

	if command -v apt-get >/dev/null 2>&1; then
		if [ "$(id -u)" -eq 0 ]; then
			echo "Installing ripgrep via apt-get (root)..."
			apt-get update
			apt-get install -y ripgrep
			return
		fi

		if command -v sudo >/dev/null 2>&1 && sudo -n true >/dev/null 2>&1; then
			echo "Installing ripgrep via apt-get (sudo)..."
			sudo -n apt-get update
			sudo -n apt-get install -y ripgrep
			return
		fi
	fi

	if command -v brew >/dev/null 2>&1; then
		echo "Installing ripgrep via Homebrew..."
		brew install ripgrep
		return
	fi

	echo "ripgrep (rg) is required but could not be installed automatically"
}

install_markdownlint() {
	if command -v markdownlint >/dev/null 2>&1; then
		return
	fi

	if ! command -v npm >/dev/null 2>&1; then
		echo "markdownlint is required but npm is unavailable"
		return
	fi

	echo "Installing markdownlint-cli via npm..."
	if npm install -g "markdownlint-cli@${MARKDOWNLINT_CLI_VERSION}" >/dev/null 2>&1; then
		return
	fi

	echo "Global npm install unavailable; using user-local prefix..."
	npm install --prefix "$LOCAL_NPM_PREFIX" "markdownlint-cli@${MARKDOWNLINT_CLI_VERSION}"
}

# --------------------------------------------- Install --------------------------------------------

if ! command -v misspell >/dev/null 2>&1; then
	install_go_tool "misspell" "github.com/client9/misspell/cmd/misspell" "$MISSPELL_VERSION"
fi

if ! command -v golangci-lint >/dev/null 2>&1; then
	install_go_tool "golangci-lint" "github.com/golangci/golangci-lint/cmd/golangci-lint" \
		"$GOLANGCI_LINT_VERSION"
fi

if ! command -v shfmt >/dev/null 2>&1; then
	install_go_tool "shfmt" "mvdan.cc/sh/v3/cmd/shfmt" "$SHFMT_VERSION"
fi

install_ripgrep
install_shellcheck
install_markdownlint

# ------------------------------------------- Validation -------------------------------------------

required_tools=(
	misspell
	golangci-lint
	shfmt
	rg
	shellcheck
	markdownlint
)

missing_tools=()
for tool_name in "${required_tools[@]}"; do
	if ! command -v "$tool_name" >/dev/null 2>&1; then
		missing_tools+=("$tool_name")
	fi
done

if [ "${#missing_tools[@]}" -gt 0 ]; then
	echo "missing required style tools: ${missing_tools[*]}"
	echo "$MESSAGE_PATH_HINT"
	exit 1
fi

echo "Style tools installed."
