#!/bin/bash

[[ -z `command -v gox` ]] && go get github.com/mitchellh/gox

gox -os="linux" -output="result/`basename $(pwd)`_{{.OS}}_{{.Arch}}" -ldflags="-s -w -X 'main.Version=`git tag|awk 'END {print}'`' -X 'main.BuildDate=`date "+%Y%m%d-%H%M"`' -X 'main.GoVersion=`go version|awk '{print $3,$4}'`' -X 'main.GitVersion=`git rev-parse HEAD`'" .
