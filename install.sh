#!/bin/sh
set -e

REPO="ma-cohen/code-cat"
BINARY="ccat"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

detect_os() {
  case "$(uname -s)" in
    Linux*)  echo "linux" ;;
    Darwin*) echo "darwin" ;;
    *)       echo "Unsupported OS: $(uname -s)" >&2; exit 1 ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    *) echo "Unsupported arch: $(uname -m)" >&2; exit 1 ;;
  esac
}

latest_version() {
  curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' \
    | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/'
}

main() {
  OS=$(detect_os)
  ARCH=$(detect_arch)
  VERSION=$(latest_version)
  ARCHIVE="${BINARY}_${OS}_${ARCH}.tar.gz"
  URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"

  echo "Installing ${BINARY} ${VERSION} (${OS}/${ARCH})..."

  TMP=$(mktemp -d)
  trap 'rm -rf "$TMP"' EXIT

  curl -fsSL "$URL" -o "${TMP}/${ARCHIVE}"
  tar -xzf "${TMP}/${ARCHIVE}" -C "$TMP"

  if [ ! -w "$INSTALL_DIR" ]; then
    echo "Installing to ${INSTALL_DIR} (requires sudo)..."
    sudo install -m755 "${TMP}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
  else
    install -m755 "${TMP}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
  fi

  echo "Installed: $(${INSTALL_DIR}/${BINARY} --version)"
}

main
