#!/usr/bin/env bash
# 生成本仓库「Linux 可直接安装部署」的 tar.gz 包（amd64 / arm64）。
# - 在 macOS 上需要已安装 zig（用于 musl 静态交叉编译 CGO/sqlite）。
# - 在 Linux 上若仅本机构建，可设置 ONLY_NATIVE=1 只打当前架构包。
#
# 产物：dist/x-ui-linux-<arch>-direct-<git-describe>.tar.gz
# 安装：解压后 sudo bash install-binary-linux.sh（安装时会联网拉取 xray 与 geo 数据）

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

VERSION="$(git describe --tags --always 2>/dev/null || echo dev)"
DIST="${ROOT}/dist"
BUILD="${DIST}/build"
mkdir -p "${BUILD}"

GOPROXY="${GOPROXY:-https://goproxy.cn,direct}"
export GOPROXY
export GOSUMDB="${GOSUMDB:-off}"

build_zig_musl() {
  local goarch="$1"
  local out="$2"
  local zig_bin
  zig_bin="$(command -v zig)"
  if [[ -z "${zig_bin}" ]]; then
    echo "错误: 未找到 zig，macOS 交叉编译需要: brew install zig" >&2
    exit 1
  fi

  local target cc cxx
  case "${goarch}" in
    amd64)
      target="x86_64-linux-musl"
      ;;
    arm64)
      target="aarch64-linux-musl"
      ;;
    *)
      echo "不支持的 GOARCH: ${goarch}" >&2
      exit 1
      ;;
  esac

  cc="${BUILD}/zig-cc-${goarch}.sh"
  cxx="${BUILD}/zig-cxx-${goarch}.sh"
  cat > "${cc}" <<EOF
#!/usr/bin/env bash
exec "${zig_bin}" cc -target ${target} "\$@"
EOF
  cat > "${cxx}" <<EOF
#!/usr/bin/env bash
exec "${zig_bin}" c++ -target ${target} "\$@"
EOF
  chmod +x "${cc}" "${cxx}"

  CGO_ENABLED=1 GOOS=linux GOARCH="${goarch}" CC="${cc}" CXX="${cxx}" \
    go build -ldflags="-w -s -linkmode external -extldflags '-static'" \
    -o "${out}" main.go
}

build_native_linux() {
  local out="$1"
  CGO_ENABLED=1 go build -ldflags="-w -s" -o "${out}" main.go
}

assemble_pkg() {
  local goarch="$1"
  local binary="$2"
  local pkg="x-ui-linux-${goarch}-direct-${VERSION}"
  local stage="${DIST}/${pkg}"

  rm -rf "${stage}"
  mkdir -p "${stage}/x-ui/bin"

  install -m 755 "${ROOT}/install-binary-linux.sh" "${stage}/install-binary-linux.sh"
  install -m 755 "${binary}" "${stage}/x-ui/x-ui"

  for f in x-ui.sh DockerInit.sh x-ui.service.debian x-ui.service.rhel x-ui.service.arch LICENSE README.md; do
    install -m 644 "${ROOT}/${f}" "${stage}/x-ui/${f}"
  done

  tar -czf "${DIST}/${pkg}.tar.gz" -C "${DIST}" "${pkg}"
  echo "已生成: ${DIST}/${pkg}.tar.gz"
}

main() {
  rm -f "${BUILD}/x-ui-linux-amd64" "${BUILD}/x-ui-linux-arm64"

  if [[ "$(uname -s)" == "Linux" && "${ONLY_NATIVE:-0}" == "1" ]]; then
    case "$(uname -m)" in
      x86_64) build_native_linux "${BUILD}/x-ui-linux-amd64"; assemble_pkg amd64 "${BUILD}/x-ui-linux-amd64" ;;
      aarch64|arm64) build_native_linux "${BUILD}/x-ui-linux-arm64"; assemble_pkg arm64 "${BUILD}/x-ui-linux-arm64" ;;
      *) echo "ONLY_NATIVE 不支持当前架构 $(uname -m)" >&2; exit 1 ;;
    esac
    return
  fi

  if [[ "$(uname -s)" == "Darwin" ]] || [[ "${CROSS:-1}" == "1" ]]; then
    build_zig_musl amd64 "${BUILD}/x-ui-linux-amd64"
    build_zig_musl arm64 "${BUILD}/x-ui-linux-arm64"
    file "${BUILD}/x-ui-linux-amd64" "${BUILD}/x-ui-linux-arm64"
    assemble_pkg amd64 "${BUILD}/x-ui-linux-amd64"
    assemble_pkg arm64 "${BUILD}/x-ui-linux-arm64"
    return
  fi

  echo "请在 Linux 上设置 ONLY_NATIVE=1 或安装 zig 后重试" >&2
  exit 1
}

main "$@"
