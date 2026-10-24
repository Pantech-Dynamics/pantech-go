#!/usr/bin/env sh
# Copy the public API spec into testdata/public-api.yaml.
#
#   scripts/sync-spec.sh                 download the published spec
#   scripts/sync-spec.sh <path-or-url>   use a local file or another URL
set -eu

DEFAULT_URL="https://docs.pantechdynamics.com/openapi/public-api.yaml"
SRC="${1:-$DEFAULT_URL}"
DEST="$(dirname "$0")/../testdata/public-api.yaml"

case "$SRC" in
  http://* | https://*)
    curl --fail --silent --show-error --location --retry 3 -o "$DEST" "$SRC"
    ;;
  *)
    cp "$SRC" "$DEST"
    ;;
esac
echo "Wrote testdata/public-api.yaml from $SRC"
