#!/usr/bin/env sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
out_dir="$root/plugins/canary-compact/bin"
package="./cmd/canary-compact-hook"

build_one() {
  goos="$1"
  goarch="$2"
  ext="$3"
  target_dir="$out_dir/$goos-$goarch"
  mkdir -p "$target_dir"

  output="$target_dir/canary-compact-hook$ext"
  echo "building $goos/$goarch -> $output"
  (
    cd "$root"
    CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
      go build -trimpath -buildvcs=false -ldflags "-s -w -buildid=" -o "$output" "$package"
  )

  if [ "$ext" = "" ]; then
    chmod 0755 "$output"
  fi
}

build_one darwin amd64 ""
build_one darwin arm64 ""
build_one linux amd64 ""
build_one linux arm64 ""
build_one windows amd64 ".exe"
chmod 0755 "$out_dir/canary-compact-hook"
