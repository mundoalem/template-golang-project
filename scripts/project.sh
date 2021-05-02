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

_clean() {
  rm -f "$BUILDDIR/$PROGNAME"*
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
  version=$(
    cat cmd/golang-project-template/*.go |
    sed -n 's/^const\ Version\ =\ \"\(.*\)\"$/\1/p' |
    head -1
  )

  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  arch=$(uname -m | tr '[:upper:]' '[:lower:]')

  artifact="${PROGNAME}-${version}-${os}-${arch}.tar.gz"
  tarball="$BUILDDIR/$artifact"
    
  test -f "$BUILDDIR/$PROGNAME" || _build
  tar czf "$tarball" "$BUILDDIR/$PROGNAME"

  echo "::set-output name=version::v$version"
}

_reset() {
  rm -rf ${VENDORDIR:?}/*
}

_run() {
  _build
  "$BUILDDIR/$PROGNAME" "$@"
}

_test() {
  go test -v -count=1 -race ./...
}

_main() {
  option="$1"

  shift

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
    test)
      _test
      ;;
    *)
      echo "$PROGNAME: illegal option -- $option"
  esac
}

# #############################################################################

_main "$@"
