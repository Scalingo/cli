#!/bin/bash

function sha256Check() {
  cmd=()
  if command -v shasum &> /dev/null; then
    cmd=(shasum -a 256)
  elif command -v sha256sum &> /dev/null; then
    cmd=(sha256sum)
  else
    error "shasum commands could not be found, please install it first"
    exit 1
  fi
  "${cmd[@]}" "$@"
}

main() {
  status() {
    echo -en "-----> $*"
  }

  info() {
    echo -en "       $*"
  }

  error() {
    echo -en " !     $*" >&2
  }

  warn() {
    echo -en " /!\\   $*"
  }

  clean_install() {
    tmpdir=$1

    rm -r "$tmpdir"
    # If installed through one line install, remove script
    if [ "$0" = "install" ] ; then
      rm "$0"
    fi
  }

  ask() {
    info "$* [Y/n] "
    while true; do
        read answer
        case $(echo "$answer" | tr '[:upper:]' '[:lower:]') in
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

  if [ -n "$DEBUG" ] ; then
    set -x
  fi

  uname -a | grep -qi 'Linux' ; is_linux=$?
  uname -a | grep -qi 'Darwin' ; is_darwin=$?

  os=$(uname -s | tr '[:upper:]' '[:lower:]')
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

  # Use tr's short parameter to be compatible with MacOS: https://ss64.com/osx/tr.html
  version=$(curl --silent https://cli-dl.scalingo.com/version | tr -d ' \t\n')
  if [ -z "$version" ]; then
    error "Fail to get the version of the CLI\n"
    error "You probably have an old version of curl. Please check your curl version and update accordingly.\n"
    exit 1
  fi

  dirname="scalingo_${version}_${os}_${arch}"
  archive_name="${dirname}.${ext}"
  url="https://github.com/Scalingo/cli/releases/download/${version}/${archive_name}"

  status "Downloading Scalingo client...  "
  curl --silent --fail --location --output "${tmpdir}/${archive_name}" "$url"
  if [ ! -f "${tmpdir}/${archive_name}" ]; then
    echo ""
    error "Fail to download the CLI archive\n"
    exit 1
  fi
  echo "DONE"

  status "Verifying the checksum...  "
  checksums_url="https://github.com/Scalingo/cli/releases/download/${version}/checksums.txt"
  # Use cut's short parameter to be compatible with MacOS: https://ss64.com/osx/cut.html
  checksum_computed=$(sha256Check "${tmpdir}/${archive_name}" | cut -d" " -f1)
  checksum_expected=$(curl --silent --location "$checksums_url" | grep "$archive_name" | cut -d" " -f1)
  if [[ "$checksum_computed" != "$checksum_expected" ]]; then
    echo "INVALID"
    error "Checksums don't match.\n"
    error "You may want to retry to install the Scalingo CLI. If the problem persists, please contact our support team.\n"
    exit 1
  fi
  echo "VALID"

  status "Extracting...   "
  tar -C "${tmpdir}" -x -f "${tmpdir}/${archive_name}"

  exe_path=${tmpdir}/${dirname}/scalingo
  if [ ! -f "$exe_path" ]; then
    echo ""
    error "Fail to extract the CLI archive\n"
    exit 1
  fi
  echo "DONE"

  default_target_dir="/usr/local/bin"
  if [ "$os" == "darwin" ] && [ "$arch" == "arm64" ]; then
    default_target_dir="/opt/homebrew/bin"
  fi
  target_dir="${target_dir:-$default_target_dir}"
  target="$target_dir/scalingo"

  if [ -x "$target" ] && [ -z "$yes_to_overwrite" ] ; then
    export DISABLE_UPDATE_CHECKER=true
    # Use cut's short parameter to be compatible with MacOS: https://ss64.com/osx/cut.html
    new_version=$($exe_path --version | cut -d' ' -f4)
    old_version=$("$target" --version | cut -d' ' -f4)
    warn "Scalingo client is already installed (version ${old_version})\n"

    if ! ask "Do you want to replace it with version ${new_version}?"  ; then
      status "Aborting...\n"
      exit 1
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

  $sudo mv "$exe_path" "$target" ; rc=$?

  if [ $rc -ne 0 ] ; then
    error "Fail to install Scalingo client (return $rc)\n"
  else
    status "Installation completed, the command 'scalingo' is available.\n"
    status "Here's what's new in this version:\n\n$(scalingo changelog)\n"
  fi
}

# Avoid error if download failure
main "$@"
