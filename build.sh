#!/bin/bash

[[ -z `command -v gox` ]] && go get github.com/mitchellh/gox

VERSION=`git describe --tags $(git rev-list --tags --max-count=1)`
NOW=`TZ=Asia/Shanghai date "+%Y%m%d-%H%M"`
GO_VERSION=`go version|awk '{print $3,$4}'`
GIT_VERSION=`git rev-parse HEAD`

gox -os="linux" -output="result/`basename $(pwd)`_{{.OS}}_{{.Arch}}" -ldflags="-s -w -X 'main.Version=$VERSION' -X 'main.BuildDate=$NOW' -X 'main.GoVersion=$GO_VERSION' -X 'main.GitVersion=$GIT_VERSION'" .
