#!/bin/sh
#
# Universal installer for dguide (https://github.com/Appdynamics/dguide).
#
# Local install (release tarball): run where the dguide binary is present.
#   sh install.sh
#
# Remote install (download from GitHub Releases): POSIX sh, pipe-friendly.
#   curl -fsSL https://raw.githubusercontent.com/Appdynamics/dguide/main/install.sh | sh
#   curl -fsSL …/install.sh | sh -s -- --version 0.2.0
#

set -e

# --- terminal colors (only when stdout is a TTY) ---
if [ -t 1 ]; then
	_ATTR_BOLD="$(printf '\033[1m')"
	_ATTR_RESET="$(printf '\033[0m')"
	_COLOR_RED="$(printf '\033[31m')"
	_COLOR_GREEN="$(printf '\033[32m')"
	_COLOR_YELLOW="$(printf '\033[33m')"
	_COLOR_CYAN="$(printf '\033[36m')"
else
	_ATTR_BOLD=
	_ATTR_RESET=
	_COLOR_RED=
	_COLOR_GREEN=
	_COLOR_YELLOW=
	_COLOR_CYAN=
fi

_die() {
	printf '%b\n' "${_COLOR_RED}${_ATTR_BOLD}error:${_ATTR_RESET} ${_COLOR_RED}$*${_ATTR_RESET}" >&2
	exit 1
}

_warn() {
	printf '%b\n' "${_COLOR_YELLOW}warning:${_ATTR_RESET} ${_COLOR_YELLOW}$*${_ATTR_RESET}" >&2
}

_info() {
	printf '%b\n' "${_COLOR_CYAN}$*${_ATTR_RESET}"
}

_ok() {
	printf '%b\n' "${_COLOR_GREEN}$*${_ATTR_RESET}"
}

_section() {
	printf '\n%b\n' "${_ATTR_BOLD}${_COLOR_CYAN}== $* ==${_ATTR_RESET}\n"
}

_have_cmd() {
	command -v "$1" >/dev/null 2>&1
}

_download_to() {
	_dl_url="$1"
	_dl_dest="$2"
	if _have_cmd curl; then
		curl -fsSL --proto '=https' --tlsv1.2 \
			-H 'Accept: application/vnd.github+json' \
			-o "$_dl_dest" "$_dl_url"
	elif _have_cmd wget; then
		if wget --help 2>&1 | grep -q secure-protocol; then
			wget --quiet \
				--secure-protocol=TLSv1_2 \
				-O "$_dl_dest" "$_dl_url"
		else
			wget --quiet -O "$_dl_dest" "$_dl_url"
		fi
	else
		_die 'Neither curl nor wget was found in PATH; install one of them and retry.'
	fi
}

_download_to_progress() {
	_dl_url="$1"
	_dl_dest="$2"
	if _have_cmd curl; then
		if [ -t 2 ]; then
			curl -fSL --proto '=https' --tlsv1.2 -# \
				-o "$_dl_dest" "$_dl_url"
		else
			curl -fSL --proto '=https' --tlsv1.2 \
				-o "$_dl_dest" "$_dl_url"
		fi
	elif _have_cmd wget; then
		if wget --help 2>&1 | grep -q secure-protocol; then
			wget --secure-protocol=TLSv1_2 --progress=bar:force -O "$_dl_dest" "$_dl_url"
		else
			wget --progress=bar:force -O "$_dl_dest" "$_dl_url"
		fi
	else
		_die 'Neither curl nor wget was found in PATH; install one of them and retry.'
	fi
}

_show_help() {
	printf '%s\n' 'Usage: install.sh [--version X.Y.Z]'
	printf '%s\n' '  Local:  run from an extracted release tarball (dguide binary in cwd).'
	printf '%s\n' '  Remote: download from GitHub Releases and install to /usr/local/bin.'
	printf '%s\n' '  Pipe:   curl -fsSL …/install.sh | sh -s -- --version 0.2.0'
}

# Install dguide from the current directory into /usr/local/bin.
_install_local() {
	_BINARY_NAME="dguide"
	_DEST_DIR="/usr/local/bin"
	_OS="$(uname -s)"

	if [ ! -f "./${_BINARY_NAME}" ]; then
		_die "Binary ${_BINARY_NAME} not found in $(pwd). Extract the release tarball first, or run without a local binary to download from GitHub."
	fi

	# _section 'Installing to /usr/local/bin'

	if [ ! -w "$_DEST_DIR" ]; then
		if [ "$_OS" = "Linux" ] || [ "$_OS" = "Darwin" ]; then
			_info "The destination directory $_DEST_DIR is not writable. Trying with sudo..."
			sudo mv "./${_BINARY_NAME}" "${_DEST_DIR}/" || _die "Failed to move the binary to ${_DEST_DIR}"
		else
			_die "The destination directory ${_DEST_DIR} is not writable. Please run with appropriate permissions."
		fi
	else
		mv "./${_BINARY_NAME}" "${_DEST_DIR}/" || _die "Failed to move the binary to ${_DEST_DIR}"
	fi

	if command -v "$_BINARY_NAME" >/dev/null 2>&1; then
		_ok "${_BINARY_NAME} has been successfully installed"
	else
		_die "Failed to verify the installation of ${_BINARY_NAME}"
	fi
}

# Download release archive, verify checksum, extract, install binary.
_install_remote() {
	_VERSION_OVERRIDE=""
	while [ "$#" -gt 0 ]; do
		case "$1" in
		--version)
			if [ -z "${2:-}" ]; then
				_die '--version requires a value (e.g. 0.2.0)'
			fi
			_VERSION_OVERRIDE="$2"
			shift 2
			;;
		-h | --help)
			_show_help
			exit 0
			;;
		*)
			_die "Unknown option: $1 (try --help)"
			;;
		esac
	done

	_OS_RAW="$(uname -s)"
	_ARCH_RAW="$(uname -m)"

	case "$_OS_RAW" in
	Linux) _OS_SLUG=linux ;;
	Darwin) _OS_SLUG=darwin ;;
	*)
		_die "Unsupported operating system: $_OS_RAW (expected Linux or Darwin)."
		;;
	esac

	case "$_ARCH_RAW" in
	x86_64 | amd64) _ARCH_SLUG=amd64 ;;
	aarch64 | arm64) _ARCH_SLUG=arm64 ;;
	i386 | i686) _ARCH_SLUG=386 ;;
	*)
		_die "Unsupported CPU architecture: $_ARCH_RAW (expected x86_64, aarch64/arm64, or i386/i686)."
		;;
	esac

	_section "dguide installer"
	_info "Detected platform: ${_OS_SLUG}/${_ARCH_SLUG}"

	_TAG_FOR_URL=""
	_VER_PLAIN=""

	if [ -n "$_VERSION_OVERRIDE" ]; then
		_VER_PLAIN="$(printf '%s' "$_VERSION_OVERRIDE" | sed 's/^v//')"
		if [ -z "$_VER_PLAIN" ]; then
			_die "Invalid version string: $_VERSION_OVERRIDE (use e.g. 0.2.0 or v0.2.0)"
		fi
		case "$_VER_PLAIN" in
		*[!0-9.]*)
			_die "Invalid version string: $_VERSION_OVERRIDE (use e.g. 0.2.0 or v0.2.0)"
			;;
		esac
		_TAG_FOR_URL="v${_VER_PLAIN}"
	else
		if ! _have_cmd curl && ! _have_cmd wget; then
			_die 'Need curl or wget to query the latest release. Install one, or pass --version X.Y.Z explicitly.'
		fi
		_section "Detecting latest release"
		_JSON_TMP="$(mktemp "${TMPDIR:-/tmp}/dguide-release-json.XXXXXX")"
		trap 'rm -f "$_JSON_TMP"' EXIT
		_rel_url='https://api.github.com/repos/Appdynamics/dguide/releases/latest'
		set +e
		if _have_cmd curl; then
			curl -fsSL --proto '=https' --tlsv1.2 \
				-H 'Accept: application/vnd.github+json' \
				-o "$_JSON_TMP" "$_rel_url"
			_API_EC=$?
		else
			if wget --help 2>&1 | grep -q secure-protocol; then
				wget --quiet --secure-protocol=TLSv1_2 -O "$_JSON_TMP" "$_rel_url"
			else
				wget --quiet -O "$_JSON_TMP" "$_rel_url"
			fi
			_API_EC=$?
		fi
		set -e
		if [ "$_API_EC" -ne 0 ] || [ ! -s "$_JSON_TMP" ]; then
			rm -f "$_JSON_TMP"
			trap - EXIT
			_die "Could not fetch latest release metadata from GitHub (network error or rate limit).
Specify a version explicitly, for example:

  curl -fsSL https://raw.githubusercontent.com/Appdynamics/dguide/main/install.sh | sh -s -- --version 0.2.0"
		fi
		_TAG_RAW="$(grep '"tag_name"' "$_JSON_TMP" | head -n 1 | sed -e 's/.*"tag_name"[[:space:]]*:[[:space:]]*"//' -e 's/".*//')"
		rm -f "$_JSON_TMP"
		trap - EXIT
		if [ -z "$_TAG_RAW" ]; then
			_die 'Could not parse tag_name from GitHub API response.
Specify --version X.Y.Z explicitly.'
		fi
		case "$_TAG_RAW" in
		v*) _TAG_FOR_URL="$_TAG_RAW" ;;
		*) _TAG_FOR_URL="v${_TAG_RAW}" ;;
		esac
		_VER_PLAIN="$(printf '%s' "$_TAG_FOR_URL" | sed 's/^v//')"
		_ok "Latest release tag: ${_TAG_FOR_URL}"
	fi

	_ARCHIVE_NAME="dguide_${_VER_PLAIN}_${_OS_SLUG}_${_ARCH_SLUG}.tar.gz"
	_DL_BASE="https://github.com/Appdynamics/dguide/releases/download/${_TAG_FOR_URL}"
	_ARCHIVE_URL="${_DL_BASE}/${_ARCHIVE_NAME}"
	_CHECKSUMS_URL="${_DL_BASE}/dguide_checksums.txt"

	_section "Downloading ${_ARCHIVE_NAME}"

	_TMPROOT="$(mktemp -d "${TMPDIR:-/tmp}/dguide-install.XXXXXX")"
	cleanup_tmproot() {
		rm -rf "$_TMPROOT"
	}
	trap cleanup_tmproot EXIT INT TERM HUP

	_ARCHIVE_PATH="${_TMPROOT}/${_ARCHIVE_NAME}"
	_CHECKSUMS_PATH="${_TMPROOT}/dguide_checksums.txt"

	_download_to_progress "$_ARCHIVE_URL" "$_ARCHIVE_PATH"

	_VERIFY=0
	if _have_cmd sha256sum || _have_cmd shasum; then
		_VERIFY=1
	fi

	if [ "$_VERIFY" -eq 1 ]; then
		_info 'Fetching checksums (dguide_checksums.txt)...'
		_download_to "$_CHECKSUMS_URL" "$_CHECKSUMS_PATH"
		tr -d '\015' <"$_CHECKSUMS_PATH" >"${_TMPROOT}/checksums.txt"
		_CHK_LINE="$(grep "[[:space:]]${_ARCHIVE_NAME}[[:space:]]*$" "${_TMPROOT}/checksums.txt" | head -n 1)"
		if [ -z "$_CHK_LINE" ]; then
			_CHK_LINE="$(grep "${_ARCHIVE_NAME}" "${_TMPROOT}/checksums.txt" | head -n 1)"
		fi
		if [ -z "$_CHK_LINE" ]; then
			_die "Checksums file did not contain an entry for ${_ARCHIVE_NAME}."
		fi
		printf '%s\n' "$_CHK_LINE" >"${_TMPROOT}/checksum-line.txt"
		_info 'Verifying SHA-256 checksum...'
		set +e
		if _have_cmd sha256sum; then
			(cd "$_TMPROOT" && sha256sum -c "checksum-line.txt")
			_VE=$?
		elif _have_cmd shasum; then
			(cd "$_TMPROOT" && shasum -a 256 -c "checksum-line.txt")
			_VE=$?
		else
			_VE=0
		fi
		set -e
		if [ "$_VE" -ne 0 ]; then
			_die 'SHA-256 verification failed - the archive may be corrupted or tampered with.'
		fi
		_ok 'Checksum verified.'
	else
		_warn 'sha256sum/shasum not found; skipping checksum verification.'
	fi

	# _section 'Extracting archive'
	(
		cd "$_TMPROOT" || exit 1
		tar -xzf "$_ARCHIVE_PATH"
	)

	xattr -d com.apple.quarantine "${_TMPROOT}/dguide" 2>/dev/null || true
	chmod +x "${_TMPROOT}/dguide" 2>/dev/null || true

	(
		cd "$_TMPROOT" || exit 1
		_install_local
	)

	if [ "$_OS_RAW" = "Linux" ] && _have_cmd getenforce; then
		_SEL="$(getenforce 2>/dev/null || true)"
		if [ "$_SEL" = 'Enforcing' ]; then
			_warn 'SELinux is Enforcing. If dguide fails to run or install is blocked, you may need to adjust SELinux (e.g. permissive mode for testing) or apply an appropriate policy — see INSTALL.md.'
		fi
	fi

	_section 'Verifying CLI'
	if ! command -v dguide >/dev/null 2>&1; then
		_die 'dguide was not found on PATH after install. Ensure /usr/local/bin is on your PATH.'
	fi
	set +e
	_VERSION_OUT="$(dguide version 2>&1)"
	_DV=$?
	set -e
	if [ "$_DV" -ne 0 ]; then
		_die "Post-install check failed: \`dguide version\` exited with status ${_DV}.

${_VERSION_OUT}"
	fi

	_ok "$_VERSION_OUT"

	_section 'Done'
	_ok "dguide ${_VER_PLAIN} is installed successfully."

	cleanup_tmproot
	trap - EXIT INT TERM HUP
}

# --- entry ---
case "${1:-}" in
-h | --help)
	_show_help
	exit 0
	;;
esac

if [ -f ./dguide ]; then
	_install_local
else
	_install_remote "$@"
fi
