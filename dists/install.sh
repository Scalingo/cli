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

  usage() {
    echo "Installs scalingo client."
    echo
    echo "Options:"
    echo "  -h, --help             displays help and exits"
    echo "  -i, --install-dir DIR  scalingo client installation directory, creating it if"
    echo "                         necessary (defaults to /usr/local/bin)"
    echo "  -y, --yes              overwrites previously installed scalingo client"
    echo
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

  while [ "$#" -gt "0" ]
  do
    key="$1"

    case $key in
      -h|--help)
      usage
      exit
      ;;
      -i|--install-dir)
      target_dir="$2"
      shift
      shift
      if [ -e "$target_dir" ] && [ ! -d "$target_dir" ] ; then
        error "target directory '$target_dir' exists but is not a directory\n"
        exit 1
      fi
      ;;
      -y|--yes)
      yes_to_overwrite=1
      shift
      ;;
      *)
      usage
      error "unknown argument $1\n"
      exit 1
      ;;
    esac
  done

  tmpdir=$(mktemp -d /tmp/scalingo_cli_XXX)
  trap "clean_install ${tmpdir}" EXIT
  version=$(curl -s https://cli-dl.scalingo.com/version | tr -d ' \t\n')
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
  target_dir="${target_dir:-/usr/local/bin}"
  target="$target_dir/scalingo"

  if [ -x "$target" -a -z "$yes_to_overwrite" ] ; then
    new_version=$($exe_path -v | cut -d' ' -f4)
    old_version=$("$target" -v | cut -d' ' -f4)
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

  status "Install scalingo client to $target_dir\n"
  if [ ! -w "$target_dir" ] ; then
    sudo=sudo
    info "sudo required...\n"
  fi

  if [ ! -e "$target_dir" ] ; then
    info "$target_dir does not exist, creating...\n"
    if [ -w "$(basename "$target_dir")" ] ; then
      mkdir -p "$target_dir"
    else
      $sudo mkdir -p "$target_dir"
    fi
  fi

  $sudo mv $exe_path "$target" ; rc=$?

  if [ $rc -ne 0 ] ; then
    error "Fail to install scalingo client (return $rc)\n"
  else
    status "Installation completed, the command 'scalingo' is available.\n"
  fi
}

# Avoid error if download failure
main "$@"
