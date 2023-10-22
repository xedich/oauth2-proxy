#!/usr/bin/env bash
set -euxo pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# make oauth2-proxy importable for others to use as a library
mkdir self
mv static self
mv oauthproxy*.go self
mv version.go self

COMMAND=sed
if [[ "$OSTYPE" == "darwin"* ]]; then
    if (! command -v gsed); then
        brew install gnu-sed
    fi
    COMMAND=gsed
fi

find self -name '*.go' -exec "$COMMAND" -i 's|package main|package self|' {} \;
"$COMMAND" -i 's|VERSION,|self.VERSION,|' main.go
"$COMMAND" -i 's| NewOAuthProxy| self.NewOAuthProxy|' main.go
"$COMMAND" -i 's|import (|import ( "github.com/oauth2-proxy/oauth2-proxy/v7/self"|' main.go
go fmt main.go

# add extra function for oauth2-proxy to be used by nps
cp "$SCRIPT_DIR/oauthproxy_extra.go.res" "self/oauthproxy_extra.go"