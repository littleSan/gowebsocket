#! /bin/bash

PWD="$(cd `dirname $0`; pwd)"
NAME="im-server"


##CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./$NAME -ldflags="-w -s"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./$NAME

echo "Success!"~