#!/bin/bash

cd `dirname $0`
GOOS=linux GOARCH=amd64 go build -o hello
zip handler.zip ./hello