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

  ask() {
    info "$* [Y/n] "
    while true; do
        read answer
        case $(echo "$answer" | tr "[A-Z]" "[a-z]") in
        y|yes|"" ) return 0;;
        n|no     ) return 1;;
        esac
        info "Invalid input, please enter 'y' or 'n': "
    done
  }

  usage() {
    echo "Installs Scalingo client."
    echo
    echo "Options:"
    echo "  -h, --help             displays help and exits"
    echo "  -i, --install-dir DIR  Scalingo client installation directory, creating it if"
    echo "                         necessary (defaults to /usr/local/bin or /opt/homebrew/bin for Apple Silicon)"
    echo "  -y, --yes              overwrites previously installed Scalingo client"
    echo
  }

  if [ "x$DEBUG" = "xtrue" ] ; then
    set -x
  fi

  uname -a | grep -qi 'Linux' ; is_linux=$?
  uname -a | grep -qi 'Darwin' ; is_darwin=$?

  os=$(uname -s | tr '[A-Z]' '[a-z]')
  ext='tar.gz'
  if [ "$os" != "linux" ] && [ "$os" != "darwin" ]; then
    echo "Unsupported OS: $(uname -s)"
    exit 1
  fi

  arch=$(uname -m)
  case $arch in
    x86_64)
      arch=amd64
      ;;
    i686)
      arch=386
      ;;
    aarch64)
      arch=arm64
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

  version=$(curl --silent https://cli-dl.scalingo.com/version | tr -d ' \t\n')
  if [ -z "$version" ]; then
    echo "-----> Fail to get the version of the CLI" >&2
    echo "You probably have an old version of curl. Please check your curl version and update accordingly." >&2
    exit 1
  fi

  dirname="scalingo_${version}_${os}_${arch}"
  archive_name="${dirname}.${ext}"
  url=https://github.com/Scalingo/cli/releases/download/${version}/${archive_name}

  status "Downloading Scalingo client...  "
  curl --silent --fail --location --output ${tmpdir}/${archive_name} ${url}
  if [ ! -f ${tmpdir}/${archive_name} ]; then
    echo "" >&2
    echo "-----> Fail to download the CLI archive" >&2
    exit 1
  fi
  echo "DONE"
  status "Extracting...   "
  tar -C "${tmpdir}" -x -f "${tmpdir}/${archive_name}"

  exe_path=${tmpdir}/${dirname}/scalingo
  if [ ! -f "$exe_path" ]; then
    echo "" >&2
    echo "-----> Fail to extract the CLI archive" >&2
    exit 1
  fi
  echo "DONE"

  default_target_dir="/usr/local/bin"
  if [ "$os" == "darwin" ] && [ "$arch" == "arm64" ]; then
    default_target_dir="/opt/homebrew/bin"
  fi
  target_dir="${target_dir:-$default_target_dir}"
  target="$target_dir/scalingo"

  if [ -x "$target" -a -z "$yes_to_overwrite" ] ; then
    export DISABLE_UPDATE_CHECKER=true
    new_version=$($exe_path --version | cut -d' ' -f4)
    old_version=$("$target" --version | cut -d' ' -f4)
    warn "Scalingo client is already installed (version ${old_version})\n"

    if ! ask "Do you want to replace it with version ${new_version}?"  ; then
      status "Aborting...\n"
      exit -1
    fi
  fi

  status "Install Scalingo client to $target_dir\n"
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
    error "Fail to install Scalingo client (return $rc)\n"
  else
    status "Installation completed, the command 'scalingo' is available.\n"
    status "Here's what's new in this version:\n\n$(scalingo changelog)\n"
  fi
}

# Avoid error if download failure
main "$@"
