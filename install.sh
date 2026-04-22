#!/bin/sh
set -e

REPO="ma-cohen/code-cat"
BINARY="ccat"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

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

warn_install_dir_not_traversable() {
  echo "Warning: ${INSTALL_DIR} is not executable for your user, so you cannot run ${BINARY} until that is fixed." >&2
  echo "  Example: sudo chmod 755 ${INSTALL_DIR}" >&2
  echo "  Or use a user-writable INSTALL_DIR (default is \$HOME/.local/bin)." >&2
}

print_path_hint() {
  case ":${PATH}:" in
    *:"${INSTALL_DIR}":*) return 0 ;;
  esac
  echo ""
  echo "Add ${BINARY} to your PATH (then open a new terminal):"
  echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
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

  mkdir -p "$INSTALL_DIR"

  used_sudo=0
  if [ ! -w "$INSTALL_DIR" ]; then
    echo "Installing to ${INSTALL_DIR} (requires sudo)..."
    used_sudo=1
    sudo install -m755 "${TMP}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
  else
    install -m755 "${TMP}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
  fi

  version_line=""
  if version_line=$("${INSTALL_DIR}/${BINARY}" --version 2>/dev/null); then
    echo "Installed: ${version_line}"
  elif [ "$used_sudo" = 1 ] && version_line=$(sudo "${INSTALL_DIR}/${BINARY}" --version 2>/dev/null); then
    echo "Installed: ${version_line}"
    if [ ! -x "$INSTALL_DIR" ]; then
      warn_install_dir_not_traversable
    fi
  else
    echo "Installed to ${INSTALL_DIR}/${BINARY}"
    if [ "$used_sudo" = 1 ] && [ ! -x "$INSTALL_DIR" ]; then
      warn_install_dir_not_traversable
    fi
  fi

  print_path_hint
}

main
