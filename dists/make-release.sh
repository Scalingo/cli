#!/bin/bash

[ "$DEBUG" = "1" ] && set -x
set -e
set -x

VERSION=""
BUILD_ONLY=false

while getopts :v:b OPT; do
  case $OPT in
    v)
      VERSION=$OPTARG
      ;;
    b)
      BUILD_ONLY=true
      ;;
  esac
done

if [ -z $VERSION ] ; then
  echo "Usage: $0 -v <version>" >&2
  exit 1
fi

bin_dir="bin/$VERSION"
mkdir -p $bin_dir

read -p "Which Rollbar token should be used in this release: " ROLLBAR_TOKEN
if [[ -z $ROLLBAR_TOKEN ]]; then
  echo "Rollbar token is mandatory" >&2
  exit 2
fi

function build_for() {
  local os=$1
  local archive_type=$2

  for arch in amd64 arm64 386 ; do
    if [ "$os" = "darwin" ] && [ "$arch" = "386" ] ; then
      continue
    fi

    pushd scalingo

    [ -e "./scalingo" ] && rm ./scalingo
    [ -e "./scalingo.exe" ] && rm ./scalingo.exe
    GOOS=$os GOARCH=$arch go build -ldflags " \
      -X main.buildstamp=$(date -u '+%Y-%m-%d_%I:%M:%S%p') \
      -X main.githash=$(git rev-parse HEAD) \
      -X main.VERSION=$VERSION \
      -X github.com/Scalingo/cli/config.RollbarToken=$ROLLBAR_TOKEN"

    release_dir="scalingo_${VERSION}_${os}_${arch}"
    archive_dir="$bin_dir/$release_dir"

    popd
    mkdir -p $archive_dir

    bin="scalingo/scalingo"
    if [ "$os" = "windows" ] ; then
      bin="scalingo/scalingo.exe"
    fi
    cp $bin $archive_dir
    cp README.md $archive_dir
    cp LICENSE $archive_dir

    pushd $bin_dir
    if [ "$archive_type" = "tarball" ] ; then
      tar czvf "${release_dir}.tar.gz" "$release_dir"
    else
      zip "${release_dir}.zip" $(find "${release_dir}")
    fi
    popd
  done
}

if [ "$BUILD_ONLY" = "false" ]; then
  current_version=$(cat VERSION)
  files=(.goxc.json README.md VERSION config/version.go)

  sed -i "s/To be Released/To be Released\n\n### ${VERSION}/g" CHANGELOG.md
  git add CHANGELOG.md

  for file in ${files[@]}; do
    sed -i "s/${current_version}/${VERSION}/g" $file
    git add $file
  done
  echo "Tagging version ${VERSION}"
  git commit -m "Bump ${VERSION}"
  git tag $VERSION
fi

if uname -a | grep -iq Linux ; then
  build_for "linux" "tarball"
  build_for "freebsd"
  build_for "openbsd"
  build_for "darwin"
  build_for "windows"
fi
if uname -a | grep -iq Darwin ; then
  build_for "darwin"
fi
if uname -a | grep -iq Mingw ; then
  build_for "windows"
fi
if uname -a | grep -iq Cygwin ; then
  build_for windows
fi

if [ "$BUILD_ONLY" = "false" ] ; then
  git push origin master
  git push origin $VERSION

  echo "Steps to create the release:"
  echo "- Go to https://github.com/Scalingo/cli/releases/new"
  echo "- Set the title to v${VERSION}"
  echo "- Set the branch to ${VERSION}"
  echo "- Add all zip and tar.gz from bin/${VERSION}"
  echo "- Set the content to:"
  sed -n "/$VERSION/,/$current_version/p" CHANGELOG.md | sed -e '1d;$d'
  echo -e "\n\nOnce done, restart the cli-download-service with"
  echo "scalingo --region osc-fr1 -a cli-download-service restart"
fi

exit 1
