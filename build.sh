#!/bin/sh
set -eux

export AWS_REGION=us-east-1

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

for func in lambdas/*; do
  cd "$func"
  go build -ldflags="-s -w -buildid=" -trimpath -o bootstrap
  cd -
done

stackit up \
  --stack-name whodunnit-lambdas \
  --template cfn.yml
