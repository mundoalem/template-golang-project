#!/usr/bin/env sh

# #############################################################################
# GENERAL
# #############################################################################

PROGNAME=$(basename "$PWD")

# #############################################################################
# DIRECTORIES
# #############################################################################

BUILDDIR="build"
CMDDIR="cmd"
INTERNALDIR="internal"
PKGDIR="pkg"
VENDORDIR="vendor"

# #############################################################################
# GO
# #############################################################################

export GO111MODULE="on"
export GO15VENDOREXPERIMENT="1"

# #############################################################################
# FUNCTIONS
# #############################################################################

_build() {
  _clean
  go build -o "$BUILDDIR/$PROGNAME" "$CMDDIR/$PROGNAME/main.go"
  chmod +x "$BUILDDIR/$PROGNAME"
}

_dep() {
  test $# -eq 1 || return 0

  dep="$1"

  which "$dep" >/dev/null 2>&1

  if [ $? -ne 0 ]; then
    echo "error: $dep not found in your \$PATH"
    return 1
  fi

  return 0
}

_clean() {
  rm -f "$BUILDDIR/$PROGNAME"*
  rm -f coverage.txt
}

_lint() {
  go fmt "./$CMDDIR/..."
  go fmt "./$INTERNALDIR/..."
  go fmt "./$PKGDIR/..."
}

_lock() {
  go mod vendor
  go mod tidy
}

_release() {
  test -f "$BUILDDIR/$PROGNAME" || _build

  version=$(
    cat "cmd/$PROGNAME/"*.go |
    sed -n 's/^const\ Version\ =\ \"\(.*\)\"$/\1/p' |
    head -1
  )

  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  arch=$(uname -m | tr '[:upper:]' '[:lower:]')

  artifact="${PROGNAME}-${version}-${os}-${arch}.tar.gz"
  tarball="$BUILDDIR/$artifact"
    
  test -f "$BUILDDIR/$PROGNAME" || _build
  tar czf "$tarball" "$BUILDDIR/$PROGNAME"

  if [ "$CI" == "true" ]; then
    echo "::set-output name=version::v$version"
  else
    echo "$PROGNAME v$version"
  fi
}

_reset() {
  rm -rf ${VENDORDIR:?}/*
}

_run() {
  _build
  "$BUILDDIR/$PROGNAME" "$@"
}

_scan() {
  _dep "snyk" || exit 1
  snyk test --fail-on=upgradable
}

_test() {
  go test -v -count=1 -race -coverprofile=coverage.txt -covermode=atomic ./...
}

_main() {
  option="$1"

  shift

  _dep "go" || exit 1

  case $option in
    build)
      _build
      ;;
    clean)
      _clean
      ;;
    lint)
      _lint
      ;;
    lock)
      _lock
      ;;
    release)
      _release
      ;;
    reset)
      _reset
      ;;
    run)
      _run "$@"
      ;;
    scan)
      _scan
      ;;
    test)
      _test
      ;;
    *)
      echo "$PROGNAME: illegal option -- $option"
  esac
}

# #############################################################################

_main "$@"
