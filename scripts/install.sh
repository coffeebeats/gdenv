#!/bin/sh
set -e

# This script installs 'gdenv' by downloading prebuilt binaries from the
# project's GitHub releases page. By default the latest version is installed,
# but a different release can be used instead by setting $GDENV_VERSION.
#
# The script will set up a 'gdenv' cache at '$HOME/.gdenv'. This behavior can
# be customized by setting '$GDENV_HOME' prior to running the script. Existing
# Godot artifacts cached in a 'gdenv' store won't be lost, but this script will
# overwrite any 'gdenv' binary artifacts in '$GDENV_HOME/bin'.

# ------------------------------ Define: Cleanup ----------------------------- #

trap cleanup EXIT

cleanup() {
    if [ -d "${GDENV_TMP=}" ]; then
        rm -rf "${GDENV_TMP}"
    fi
}

# ------------------------------ Define: Logging ----------------------------- #

info() {
    if [ "$1" != "" ]; then
        echo info: "$@"
    fi
}

warn() {
    if [ "$1" != "" ]; then
        echo warning: "$1"
    fi
}

error() {
    if [ "$1" != "" ]; then
        echo error: "$1" >&2
    fi
}

fatal() {
    error "$1"
    exit 1
}

unsupported_platform() {
    error "$1"
    echo "See https://github.com/coffeebeats/gdenv/blob/main/docs/installation.md#install-from-source for instructions on compiling from source."
    exit 1
}

# ------------------------------- Define: Usage ------------------------------ #

usage() {
    cat <<EOF
gdenv-install: Install 'gdenv' for managing multiple versions of the Godot editor.

Usage: gdenv-install [OPTIONS]

NOTE: The following dependencies are required:
    - curl OR wget
    - grep
    - sha256sum OR shasum
    - tar/unzip
    - tr
    - uname

Available options:
    -h, --help          Print this help and exit
    -v, --verbose       Print script debug info (default=false)
    --no-modify-path    Do not modify the \$PATH environment variable
EOF
    exit
}

check_cmd() {
    command -v "$1" >/dev/null 2>&1
}

need_cmd() {
    if ! check_cmd "$1"; then
        fatal "required command not found: '$1'"
    fi
}

# ------------------------------ Define: Params ------------------------------ #

parse_params() {
    MODIFY_PATH=1

    while :; do
        case "${1:-}" in
        -h | --help) usage ;;
        -v | --verbose) set -x ;;

        --no-modify-path) MODIFY_PATH=0 ;;

        -?*) fatal "Unknown option: $1" ;;
        "") break ;;
        esac
        shift
    done

    return 0
}

parse_params "$@"

# ------------------------------ Define: Version ----------------------------- #

GDENV_VERSION="${GDENV_VERSION:-0.6.15}" # x-release-please-version
GDENV_VERSION="v${GDENV_VERSION#v}"

# ----------------------------- Define: Platform ----------------------------- #

need_cmd tr
need_cmd uname

GDENV_CLI_OS="$(echo "${GDENV_CLI_OS=$(uname -s)}" | tr '[:upper:]' '[:lower:]')"
case "$GDENV_CLI_OS" in
darwin*) GDENV_CLI_OS="macos" ;;
linux*) GDENV_CLI_OS="linux" ;;
mac | macos | osx) GDENV_CLI_OS="macos" ;;
cygwin*) GDENV_CLI_OS="windows" ;;
msys* | mingw64*) GDENV_CLI_OS="windows" ;;
uwin* | win*) GDENV_CLI_OS="windows" ;;
*) unsupported_platform "no prebuilt binaries available for operating system: $GDENV_CLI_OS" ;;
esac

GDENV_CLI_ARCH="$(echo ${GDENV_CLI_ARCH=$(uname -m)} | tr '[:upper:]' '[:lower:]')"
case "$GDENV_CLI_ARCH" in
aarch64 | arm64)
    GDENV_CLI_ARCH="arm64"
    if [ "$GDENV_CLI_OS" != "macos" ] && [ "$GDENV_CLI_OS" != "linux" ]; then
        fatal "no prebuilt '$GDENV_CLI_ARCH' binaries available for operating system: $GDENV_CLI_OS"
    fi

    ;;
amd64 | x86_64) GDENV_CLI_ARCH="x86_64" ;;
*) unsupported_platform "no prebuilt binaries available for CPU architecture: $GDENV_CLI_ARCH" ;;
esac

GDENV_CLI_ARCHIVE_EXT=""
case "$GDENV_CLI_OS" in
windows) GDENV_CLI_ARCHIVE_EXT="zip" ;;
*) GDENV_CLI_ARCHIVE_EXT="tar.gz" ;;
esac

GDENV_CLI_ARCHIVE="gdenv-$GDENV_VERSION-$GDENV_CLI_OS-$GDENV_CLI_ARCH.$GDENV_CLI_ARCHIVE_EXT"

# ------------------------------- Define: Store ------------------------------ #

GDENV_HOME_PREV="${GDENV_HOME_PREV=}" # save for later in script

GDENV_HOME="${GDENV_HOME=}"
if [ "$GDENV_HOME" = "" ]; then
    if [ "${HOME=}" = "" ]; then
        fatal "both '\$GDENV_HOME' and '\$HOME' unset; one must be specified to determine 'gdenv' installation path"
    fi

    GDENV_HOME="$HOME/.gdenv"
fi

info "using 'gdenv' store path: '$GDENV_HOME'"

# ----------------------------- Define: Download ----------------------------- #

need_cmd grep
need_cmd mktemp

GDENV_TMP=$(mktemp -d --tmpdir gdenv-XXXXXXXXXX)
cd "$GDENV_TMP"

GDENV_RELEASE_URL="https://github.com/coffeebeats/gdenv/releases/download/$GDENV_VERSION"

download_with_curl() {
    curl \
        --fail \
        --location \
        --parallel \
        --retry 3 \
        --retry-delay 1 \
        --show-error \
        --silent \
        -o "$GDENV_CLI_ARCHIVE" \
        "$GDENV_RELEASE_URL/$GDENV_CLI_ARCHIVE" \
        -o "checksums.txt" \
        "$GDENV_RELEASE_URL/checksums.txt"
}

download_with_wget() {
    wget -q -t 4 -O "$GDENV_CLI_ARCHIVE" "$GDENV_RELEASE_URL/$GDENV_CLI_ARCHIVE" 2>&1
    wget -q -t 4 -O "checksums.txt" "$GDENV_RELEASE_URL/checksums.txt" 2>&1
}

if check_cmd curl; then
    download_with_curl
elif check_cmd wget; then
    download_with_wget
else
    fatal "missing one of 'curl' or 'wget' commands"
fi

# -------------------------- Define: Verify checksum ------------------------- #

verify_with_sha256sum() {
    cat "checksums.txt" | grep "$GDENV_CLI_ARCHIVE" | sha256sum --check --status
}

verify_with_shasum() {
    cat "checksums.txt" | grep "$GDENV_CLI_ARCHIVE" | shasum -a 256 -p --check --status
}

if check_cmd sha256sum; then
    verify_with_sha256sum
elif check_cmd shasum; then
    verify_with_shasum
else
    fatal "missing one of 'sha256sum' or 'shasum' commands"
fi

# ------------------------------ Define: Extract ----------------------------- #

case "$GDENV_CLI_OS" in
windows)
    need_cmd unzip

    mkdir -p "$GDENV_HOME/bin"
    unzip -u "$GDENV_CLI_ARCHIVE" -d "$GDENV_HOME/bin"
    ;;
*)
    need_cmd tar

    mkdir -p "$GDENV_HOME/bin"
    tar -C "$GDENV_HOME/bin" --no-same-owner -xzf "$GDENV_CLI_ARCHIVE"
    ;;
esac

info "successfully installed 'gdenv@$GDENV_VERSION' to '$GDENV_HOME/bin'"

if [ $MODIFY_PATH -eq 0 ]; then
    exit 0
fi

# The $PATH modification and $GDENV_HOME export is already done.
if check_cmd gdenv && [ "$GDENV_HOME_PREV" != "" ]; then
    exit 0
fi

# Simplify the exported $GDENV_HOME if possible.
if [ "$HOME" != "" ]; then
    case "$GDENV_HOME" in
    $HOME*) GDENV_HOME="\$HOME${GDENV_HOME#$HOME}" ;;
    esac
fi

CMD_EXPORT_HOME="export GDENV_HOME=\"$GDENV_HOME\""
CMD_MODIFY_PATH="export PATH=\"\$GDENV_HOME/bin:\$PATH\""

case $(basename $SHELL) in
sh) OUT="$HOME/.profile" ;;
bash) OUT="$HOME/.bashrc" ;;
zsh) OUT="$HOME/.zshrc" ;;
*)
    echo ""
    echo "Add the following to your shell profile script:"
    echo "    $CMD_EXPORT_HOME"
    echo "    $CMD_MODIFY_PATH"
    ;;
esac

if [ "$OUT" != "" ]; then
    if [ -f "$OUT" ] && $(cat "$OUT" | grep -q 'export GDENV_HOME'); then
        info "Found 'GDENV_HOME' export in shell Rc file; skipping modification."
        exit 0
    fi

    if [ -f "$OUT" ] && [ "$(tail -n 1 "$OUT")" != "" ]; then
        echo "" >>"$OUT"
    fi

    echo "# Added by 'gdenv' install script." >>"$OUT"
    echo "$CMD_EXPORT_HOME" >>"$OUT"
    echo "$CMD_MODIFY_PATH" >>"$OUT"

    info "Updated shell Rc file: $OUT\n      Open a new terminal to start using 'gdenv'."
fi
