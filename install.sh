#!/bin/sh
set -eu

REPO="zhravan/defenux"
INSTALL_DIR="${DEFENUX_INSTALL_DIR:-}"

case "$(uname -s)" in
    Linux) ;;
    *)
        echo "defenux: Linux is required" >&2
        exit 1
        ;;
esac

case "$(uname -m)" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    armv7l|armv7*)
        ARCH="armv7"
        ;;
    armv6l|armv6*)
        ARCH="armv6"
        ;;
    i386|i686)
        ARCH="386"
        ;;
    ppc64le)
        ARCH="ppc64le"
        ;;
    s390x)
        ARCH="s390x"
        ;;
    riscv64)
        ARCH="riscv64"
        ;;
    loongarch64)
        ARCH="loong64"
        ;;
    *)
        echo "defenux: unsupported Linux architecture: $(uname -m)" >&2
        exit 1
        ;;
esac

if [ -n "${DEFENUX_VERSION:-}" ]; then
    case "$DEFENUX_VERSION" in
        v*) TAG="$DEFENUX_VERSION" ;;
        *)  TAG="v$DEFENUX_VERSION" ;;
    esac
    BASE="https://github.com/$REPO/releases/download/$TAG"
else
    BASE="https://github.com/$REPO/releases/latest/download"
fi

ASSET="defenux_linux_${ARCH}.tar.gz"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Installing defenux ($ARCH)..."

curl -fsSL "$BASE/$ASSET" -o "$TMP_DIR/$ASSET"
curl -fsSL "$BASE/sha256sums.txt" -o "$TMP_DIR/sha256sums.txt"

EXPECTED="$(awk -v file="$ASSET" '$2 == file { print $1; exit }' "$TMP_DIR/sha256sums.txt")"
[ -n "$EXPECTED" ] || {
    echo "defenux: checksum not found for $ASSET" >&2
    exit 1
}

ACTUAL="$(sha256sum "$TMP_DIR/$ASSET" | awk '{ print $1 }')"
[ "$EXPECTED" = "$ACTUAL" ] || {
    echo "defenux: checksum verification failed" >&2
    exit 1
}

tar -xzf "$TMP_DIR/$ASSET" -C "$TMP_DIR"

if [ -z "$INSTALL_DIR" ]; then
    if [ -w /usr/local/bin ]; then
        INSTALL_DIR="/usr/local/bin"
    elif command -v sudo >/dev/null 2>&1; then
        INSTALL_DIR="/usr/local/bin"
    else
        INSTALL_DIR="${HOME}/.local/bin"
    fi
fi

if [ -w "$INSTALL_DIR" ] || { [ ! -e "$INSTALL_DIR" ] && [ -w "$(dirname "$INSTALL_DIR")" ]; }; then
    mkdir -p "$INSTALL_DIR"
    install -m 0755 "$TMP_DIR/defenux" "$INSTALL_DIR/defenux"
elif command -v sudo >/dev/null 2>&1; then
    sudo mkdir -p "$INSTALL_DIR"
    sudo install -m 0755 "$TMP_DIR/defenux" "$INSTALL_DIR/defenux"
else
    INSTALL_DIR="${HOME}/.local/bin"
    mkdir -p "$INSTALL_DIR"
    install -m 0755 "$TMP_DIR/defenux" "$INSTALL_DIR/defenux"
fi

echo "Installed defenux to $INSTALL_DIR/defenux"
"$INSTALL_DIR/defenux" version
