#!/bin/bash

function ensure_goxc_installation() {
  goxc_path=$(which goxc) ; rc=$?
  if [ $rc -eq 1 ] ; then
    echo "goxc not found, install..."
    go get github.com/laher/goxc
    goxc_path=$(which goxc)
  fi
}

VERSION=""
DEV_MODE=0

while getopts v:d: OPT; do
  case $OPT in
    v)
      VERSION=$OPTARG
      ;;
    d)
      echo "HELLO"
      DEV_MODE=1
      PRE_VERSION=$OPTARG
      ;;
  esac
done

if [ -z $VERSION ] ; then
  echo "$0 -v <version> [-d]"
  exit 1
fi

tag=$VERSION

if [ $DEV_MODE -eq 1 ] ; then
  ref=$(git rev-parse HEAD)
  dev_tag="dev-${ref:0:10}"
fi

if [ -n "${PRE_VERSION}" ] ; then
  dev_tag="${PRE_VERSION}"
fi

echo "tag: $tag"

git checkout dists
git rebase master

ensure_goxc_installation

goxc_flags="-pv $VERSION"
if [ -n "${dev_tag}" ] ; then
  goxc_flags="${goxc_flags} -pr $dev_tag"
fi
if uname -a | grep -iq Linux ; then
  goxc_flags="${goxc_flags} -bc linux,freebsd,openbsd"
fi
if uname -a | grep -iq Darwin ; then
  goxc_flags="${goxc_flags} -bc darwin"
fi
if uname -a | grep -iq Mingw ; then
  goxc_flags="${goxc_flags} -bc windows"
fi

goxc $goxc_flags

git checkout master
