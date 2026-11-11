#!/bin/bash
module=$(grep "module" go.mod | cut -d ' ' -f 2)
options=("linux:amd64" "linux:arm64" "windows:amd64" "darwin:arm64" "darwin:amd64")
version=$(git tag -l "v[0-99]*.[0-99]*.[0-99]*" --sort=-creatordate | head -n 1)
versionDir="$module/pkg"
appname=$(basename $module)
githubUser="xxl6097"
githubRepo="TVWeb"

function buildgo() {
  builddir=$1
  appname=$2
  version=$3
  appdir=$4
  os=$5
  arch=$6
  extra=$7

  dstFilePath=${builddir}/${appname}_${version}_${os}_${arch}
  flags=''
  if [ "${os}" = "linux" ] && [ "${arch}" = "arm" ] && [ "${extra}" != "" ]; then
    if [ "${extra}" = "7" ]; then
      flags=GOARM=7
      dstFilePath=${builddir}/${appname}_${version}_${os}_${arch}hf
    elif [ "${extra}" = "5" ]; then
      flags=GOARM=5
    fi
  elif [ "${os}" = "windows" ]; then
    dstFilePath=${builddir}/${appname}_${version}_${os}_${arch}.exe
  fi

  printf "build: GOOS=%-7s GOARCH=%-8s ==> %s\n" "${os}" "${arch}" "${dstFilePath}"

  filename=$(basename "$dstFilePath")
  binName="-X '${versionDir}.BinName=${filename}'"
  githubUserFlag="-X '${versionDir}.GithubUser=${githubUser}'"
  githubRepoFlag="-X '${versionDir}.GithubRepo=${githubRepo}'"
  # 同步注入到 TVBoxWeb 内部的 api.AppVersion / api.BuildTime，
  # 让 githubupgrade.New() 在运行时拿到正确的版本号做比对，
  # 以及让界面"构建时间"显示真实时间（默认值 "unknown"）
  tvboxVersionFlag="-X 'github.com/xxl6097/TVBoxWeb/internal/api.AppVersion=${version}'"
  tvboxBuildTimeFlag="-X 'github.com/xxl6097/TVBoxWeb/internal/api.BuildTime=$(TZ=Asia/Shanghai date "+%Y-%m-%d %H:%M:%S")'"

  env CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} ${flags} \
    go build -trimpath \
    -ldflags "$ldflags $binName $githubUserFlag $githubRepoFlag $tvboxVersionFlag $tvboxBuildTimeFlag -linkmode internal" \
    -o "${dstFilePath}" "${appdir}"

  # TVBoxWeb v1.6.0 的升级器查找固定名称 yezhen-<os>-<arch>，创建对应副本
  yezhenName="yezhen-${os}-${arch}"
  if [ "${os}" = "windows" ]; then
    yezhenName="${yezhenName}.exe"
  fi
  yezhenPath="${builddir}/${yezhenName}"
  cp "${dstFilePath}" "${yezhenPath}"
  printf "       alias: %s\n" "${yezhenPath}"
}

function buildAll() {
  builddir=$1
  appname=$2
  version=$3
  appdir=$4
  disname=$5
  describe=$6
  ldflags=$(buildLdflags "$appname" "$disname" "$describe")
  for entry in "${options[@]}"; do
    IFS=":" read -r os arch extra <<< "$entry"
    buildgo "$builddir" "$appname" "$version" "$appdir" "$os" "$arch" "$extra"
  done
}

version::get_version_vars() {
  GIT_COMMIT="$(git rev-parse HEAD^{commit})"

  if git_status=$(git status --porcelain 2>/dev/null) && [[ -z ${git_status} ]]; then
    GIT_TREE_STATE="clean"
  else
    GIT_TREE_STATE="dirty"
  fi

  if [[ -n ${GIT_VERSION-} ]] || GIT_VERSION=$(git describe --tags --abbrev=14 --match "v[0-9]*" "${GIT_COMMIT}" 2>/dev/null); then
    DASHES_IN_VERSION=$(echo "${GIT_VERSION}" | sed "s/[^-]//g")
    if [[ "${DASHES_IN_VERSION}" == "---" ]]; then
      GIT_VERSION=$(echo "${GIT_VERSION}" | sed "s/-\([0-9]\{1,\}\)-g\([0-9a-f]\{14\}\)$/.\1\-\2/")
    elif [[ "${DASHES_IN_VERSION}" == "--" ]]; then
      GIT_VERSION=$(echo "${GIT_VERSION}" | sed "s/-g\([0-9a-f]\{14\}\)$/-\1/")
    fi
    if [[ "${GIT_TREE_STATE}" == "dirty" ]]; then
      GIT_VERSION+="-dirty"
    fi
    if [[ "${GIT_VERSION}" =~ ^v([0-9]+)\.([0-9]+)(\.[0-9]+)?([-].*)?([+].*)?$ ]]; then
      GIT_MAJOR=${BASH_REMATCH[1]}
      GIT_MINOR=${BASH_REMATCH[2]}
    fi
    if ! [[ "${GIT_VERSION}" =~ ^v([0-9]+)\.([0-9]+)(\.[0-9]+)?(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$ ]]; then
      echo "GIT_VERSION should be a valid Semantic Version. Current value: ${GIT_VERSION}"
      echo "Please see more details here: https://semver.org"
      exit 1
    fi
  fi

  GIT_RELEASE_TAG=$(git describe --abbrev=0 --tags 2>/dev/null || echo "v0.0.0")
  GIT_RELEASE_COMMIT=$(git rev-list -n 1 "${GIT_RELEASE_TAG}" 2>/dev/null || echo "")
}

function buildLdflags() {
  local ldflags
  ldflags="-s -w"
  function add_ldflag() {
    local key=${1}
    local val=${2}
    ldflags+=("-X '${versionDir}.${key}=${val}'")
  }
  appname=$1
  DisplayName=$2
  Description=$3
  APP_NAME=${appname}
  BUILD_VERSION=$(if [ "$(git describe --tags --abbrev=0 2>/dev/null)" != "" ]; then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
  BUILD_TIME=$(TZ=Asia/Shanghai date "+%Y-%m-%d %H:%M:%S")
  GIT_REVISION=$(git rev-parse --short HEAD)
  GIT_BRANCH=$(git name-rev --name-only HEAD)

  version::get_version_vars
  add_ldflag "DisplayName" "${DisplayName}_${version}"
  add_ldflag "Description" "${Description}"
  add_ldflag "AppName"     "${APP_NAME}"
  add_ldflag "AppVersion"  "${version}"
  add_ldflag "BuildVersion" "${BUILD_VERSION}"
  add_ldflag "BuildTime"   "${BUILD_TIME}"
  add_ldflag "GitRevision" "${GIT_REVISION}"
  add_ldflag "GitBranch"   "${GIT_BRANCH}"
  add_ldflag "GitCommit"   "${GIT_COMMIT}"
  add_ldflag "GitTreeState" "${GIT_TREE_STATE}"
  add_ldflag "GitVersion"  "${GIT_VERSION}"
  add_ldflag "GitReleaseCommit" "${GIT_RELEASE_COMMIT}"
  echo "${ldflags[*]-}"
}

# 生成 checksums.txt（TVBoxWeb v1.6.0 升级器强制校验 SHA256）
function generateChecksums() {
  local pkgdir=$1
  local checksumsFile="${pkgdir}/checksums.txt"
  echo "生成 ${checksumsFile} ..."
  cd "${pkgdir}" || exit 1
  # 对 yezhen-* 文件（升级器实际下载的目标）和版本化文件都生成校验值
  if command -v sha256sum &>/dev/null; then
    sha256sum yezhen-* ${appname}_* > checksums.txt 2>/dev/null || true
  elif command -v shasum &>/dev/null; then
    shasum -a 256 yezhen-* ${appname}_* > checksums.txt 2>/dev/null || true
  else
    echo "警告：未找到 sha256sum 或 shasum，跳过 checksums.txt 生成" >&2
  fi
  cd - > /dev/null
  if [ -f "${checksumsFile}" ]; then
    echo "checksums.txt 已生成："
    cat "${checksumsFile}"
  fi
}

function githubActions() {
  builddir="./release"
  appdir="./cmd/app"
  disname="${appname}应用程序"
  describe="TVWeb - Video Box Web Service"
  echo "===> version: ${version}"
  rm -rf "${builddir}"
  buildAll "${builddir}" "${appname}" "${version}" "${appdir}" "${disname}" "${describe}"

  # 整理 packages 目录
  mkdir -p ./release/packages
  mv -fv ./release/${appname}* ./release/packages/
  mv -fv ./release/yezhen-*   ./release/packages/ 2>/dev/null || true

  # 生成 checksums.txt（放在 packages 目录，作为独立 Release Asset 上传）
  generateChecksums ./release/packages
}

function bootstrap() {
  if [ $# -ge 2 ] && [ -n "$2" ]; then
    version=$2
  fi
  case $1 in
    github) (githubActions) ;;
    *) echo "Usage: $0 github [version]" ;;
  esac
}

bootstrap "$1" "$2"
