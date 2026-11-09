#!/bin/bash
module=$(grep "module" go.mod | cut -d ' ' -f 2)
options=("linux:amd64" "linux:arm64" "windows:amd64" "darwin:arm64" "darwin:amd64")
version=$(git tag -l "v[0-99]*.[0-99]*.[0-99]*" --sort=-creatordate | head -n 1)
versionDir="$module/pkg"
appname=$(basename $module)
githubUser="xxl6097"
githubRepo="TVWeb"

function buildgo() {
  echo "====s=$1,$2,$3,$4,$5,$6"
  builddir=$1
  appname=$2
  version=$3
  appdir=$4
  os=$5
  arch=$6
  extra=$7
  dstFilePath=${builddir}/${appname}_${version}_${os}_${arch}
  flags='';
  if [ "${os}" = "linux" ] && [ "${arch}" = "arm" ] && [ "${extra}" != "" ] ; then
    if [ "${extra}" = "7" ]; then
      flags=GOARM=7;
      dstFilePath=${builddir}/${appname}_${version}_${os}_${arch}hf
    elif [ "${extra}" = "5" ]; then
      flags=GOARM=5;
      dstFilePath=${builddir}/${appname}_${version}_${os}_${arch}
    fi;
  elif [ "${os}" = "windows" ] ; then
    dstFilePath=${builddir}/${appname}_${version}_${os}_${arch}.exe
  fi;
  printf "build：GOOS=%-7s GOARCH=%-8s ==> %s\n" ${os} ${arch} ${dstFilePath}

  filename=$(basename "$dstFilePath")
  binName="-X '${versionDir}.BinName=${filename}'"
  githubUserFlag="-X '${versionDir}.GithubUser=${githubUser}'"
  githubRepoFlag="-X '${versionDir}.GithubRepo=${githubRepo}'"
  
  env CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} ${flags} go build -trimpath -ldflags "$ldflags $binName $githubUserFlag $githubRepoFlag -linkmode internal" -o ${dstFilePath} ${appdir}
}

function buildAll() {
  builddir=$1
  appname=$2
  version=$3
  appdir=$4
  disname=$5
  describe=$6
  ldflags=$(buildLdflags $appname $disname $describe)
  for arch in "${options[@]}"; do
      IFS=":" read -r os arch extra <<< "$arch"
      buildgo $builddir $appname $version $appdir $os $arch $extra
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
        if [[ "${DASHES_IN_VERSION}" == "---" ]] ; then
            GIT_VERSION=$(echo "${GIT_VERSION}" | sed "s/-\([0-9]\{1,\}\)-g\([0-9a-f]\{14\}\)$/.\1\-\2/")
        elif [[ "${DASHES_IN_VERSION}" == "--" ]] ; then
            GIT_VERSION=$(echo "${GIT_VERSION}" | sed "s/-g\([0-9a-f]\{14\}\)$/-\1/")
        fi
        if [[ "${GIT_TREE_STATE}" == "dirty" ]]; then
            GIT_VERSION+="-dirty"
        fi

        if [[ "${GIT_VERSION}" =~ ^v([0-9]+)\.([0-9]+)(\.[0-9]+)?([-].*)?([+].*)?$ ]]; then
            GIT_MAJOR=${BASH_REMATCH[1]}
            GIT_MINOR=${BASH_REMATCH[2]}
            GIT_MINRR=${BASH_REMATCH[3]}
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
      ldflags+=(
          "-X '${versionDir}.${key}=${val}'"
      )
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
  add_ldflag "AppName" "${APP_NAME}"
  add_ldflag "AppVersion" "${version}"
  add_ldflag "BuildVersion" "${BUILD_VERSION}"
  add_ldflag "BuildTime" "${BUILD_TIME}"
  add_ldflag "GitRevision" "${GIT_REVISION}"
  add_ldflag "GitBranch" "${GIT_BRANCH}"
  add_ldflag "GitCommit" "${GIT_COMMIT}"
  add_ldflag "GitTreeState" "${GIT_TREE_STATE}"
  add_ldflag "GitVersion" "${GIT_VERSION}"
  add_ldflag "GitReleaseCommit" "${GIT_RELEASE_COMMIT}"
  echo "${ldflags[*]-}"
}

function githubActions() {
  builddir="./release"
  appdir="./cmd/app"
  disname="${appname}应用程序"
  describe="TVWeb - Video Box Web Service"
  echo "===>version:${version}"
  rm -rf ${builddir}
  buildAll $builddir $appname "$version" $appdir $disname $describe
  mkdir -p ./release/packages
  mv -fv ./release/${appname}* ./release/packages
}

function bootstrap() {
  if [ $# -ge 2 ] && [ -n "$2" ]; then
    version=$2
  fi
  case $1 in
  github) (githubActions) ;;
    *) echo "Usage: $0 github [version]"  ;;
  esac
}

bootstrap $1 $2
