#!/bin/bash

# Override the pushd and popd commands to make them silent
pushd () {
  command pushd "$@" > /dev/null
}

popd () {
  command popd "$@" > /dev/null
}

[ "$DEBUG" = "1" ] && set -x
set -e

version=""
parallel_jobs="-1"
while getopts :v:j: OPT; do
  case $OPT in
    j)
      parallel_jobs=$OPTARG
      ;;
    v)
      version=$OPTARG
      ;;
  esac
done

if [ -z $version ] ; then
  echo "Usage: $0 -v <version>" >&2
  exit 1
fi

# Move to the directory where lie this script
_cur_dir=$(cd $(dirname $0) && pwd)
cd $_cur_dir

pushd ${_cur_dir}/../scalingo

echo "Remove previously compiled binaries"
rm -f ${_cur_dir}/../scalingo/scalingo*

bin_dir="${_cur_dir}/../bin/$version"

echo "Compiling the binaries for all architectures"
operating_systems="linux freebsd openbsd darwin windows"
architectures="amd64 arm64 386"
ignore_osarch="!darwin/386"
gox -parallel $parallel_jobs -os="${operating_systems}" -arch="${architectures}" -osarch="${ignore_osarch}" \
    -ldflags "-X main.buildstamp=$(date -u '+%Y-%m-%d_%I:%M:%S%p') \
      -X main.githash=$(git rev-parse HEAD) \
      -X main.version=$version"

echo ""
echo "Binaries compilation finished, archiving them"

for binary in ./scalingo*; do
  echo "--> Create the archive for $(basename $binary)"

  # binary contains a string like "./scalingo_windows_386.exe"
  os=$(basename $binary | cut -d "." -f 1 | cut -d "_" -f 2)
  arch=$(basename $binary | cut -d "." -f 1 | cut -d "_" -f 3)
  release_dir="scalingo_${version}_${os}_${arch}"
  archive_dir="$bin_dir/$release_dir"
  mkdir -p $archive_dir

  bin="scalingo"
  if [ "$os" = "windows" ] ; then
    bin="scalingo.exe"
  fi

  mv $binary ${archive_dir}/${bin}
  cp ${_cur_dir}/../README.md $archive_dir
  cp ${_cur_dir}/../LICENSE $archive_dir

  pushd $bin_dir
  if [ "$os" = "linux" ] ; then
    tar czf "${release_dir}.tar.gz" "$release_dir"
  else
    zip "${release_dir}.zip" $(find "${release_dir}") > /dev/null
  fi
  popd
done

popd
