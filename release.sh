#!/bin/bash

mkdir build/
for os in darwin linux windows
do
  echo "Building for OS=${os} and arch=amd64"
  env GOOS=${os} GOARCH=amd64 go build -o "build/inout_tester_${os}_amd64" -i . 
done
