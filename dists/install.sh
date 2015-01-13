#!/bin/bash

main() {
  status() {
    echo -en "-----> $*"
  }

  info() {
    echo -en "       $*"
  }

  error() {
    echo -en " !     $*"
  }
  
  warn() {
    echo -en " /!\\   $*"
  }

  clean_install() {
    tmpdir=$1

    rm -r $tmpdir
    # If installed through one line install, remove script
    if [ "x$0" = "xinstall" ] ; then
      rm "$0"
    fi
  }

  if [ "x$DEBUG" = "xtrue" ] ; then
    set -x
  fi

  uname -a | grep -qi 'Linux' ; is_linux=$?
  uname -a | grep -qi 'Darwin' ; is_darwin=$?

  os=$(uname -s | tr '[A-Z]' '[a-z]')
  ext=zip
  case $os in
    linux)
      ext='tar.gz'
      ;;
    darwin)
      ;;
    *)
      echo "Unsupported OS: $(uname -s)"
      exit 1
      ;;
  esac

  arch=$(uname -m)
  case $arch in
    x86_64)
      arch=amd64
      ;;
    i686)
      arch=386
      ;;
  esac

  tmpdir=$(mktemp -d /tmp/scalingo_cli_XXX)
  trap "clean_install ${tmpdir}" EXIT
  version=$(curl -s http://cli-dl.scalingo.io/version | tr -d ' \t\n')
  dirname="scalingo_${version}_${os}_${arch}"
  archive_name="${dirname}.${ext}"
  url=https://github.com/Scalingo/cli/releases/download/${version}/${archive_name}

  status "Downloading Scalingo client...  " 
  curl -s -L -o ${tmpdir}/${archive_name} ${url}
  echo "DONE"
  status "Extracting...   "
  case $ext in
    zip)
      unzip -d "${tmpdir}" "${tmpdir}/${archive_name}"
      ;;
    tar.gz)
      tar -C "${tmpdir}" -x -f "${tmpdir}/${archive_name}"
      ;;
  esac
  echo "DONE"

  exe_path=${tmpdir}/${dirname}/scalingo
  target=/usr/local/bin/scalingo

  if [ -x $target ] ; then
    new_version=$($exe_path -v | cut -d' ' -f4)
    old_version=$($target -v | cut -d' ' -f4)
    warn "Scalingo client is already installed (version ${old_version})\n"
    info "Do you want to replace it with version ${new_version}? [Y/n] "

    read input
    [ -z $input ] && input='Y'
    while echo $input | grep -qvE '[YyNn]' ; do
      info "Invalid input, please enter 'y' or 'n': "
      read input
    done

    if [ "x$input" = "xn" ] ; then
      status "Aborting...\n"
      exit -1
    fi
  fi

  status "Install scalingo client to /usr/local/bin\n"
  if [ ! -w /usr/local/bin ] ; then
    sudo=sudo
    info "sudo required...\n"
  fi

  $sudo mv $exe_path $target ; rc=$?

  if [ $rc -ne 0 ] ; then
    error "Fail to install scalingo client (return $rc)\n"
  else
    status "Installation completed, the command 'scalingo' is available.\n"
  fi
}

# Avoid error if download failure
main
