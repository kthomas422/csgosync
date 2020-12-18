#!/bin/sh

LINUX_ENV="env CGO_ENABLED=0 GOOS=linux GOARCH=amd64"
WINTURDS_ENV="env CGO_ENABLED=0 GOOS=windows GOARCH=amd64"
BUILD_CMD=`go build -ldflags '-s -w'`

mkdir client server
cp csgosync.yaml.example client/csgosync.yaml
cp csgosyncd.yaml.example server/csgosyncd.yaml

env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o client/csgosync cmd/client/csgosync.go
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o server/csgosyncd cmd/server/csgosyncd.go
tar -czvf csgosync-linux-x64.tgz server client
rm client/csgosync server/csgosyncd

env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags '-s -w' -o client/csgosync.exe cmd/client/csgosync.go
env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags '-s -w' -o server/csgosyncd.exe cmd/server/csgosyncd.go
zip -r csgosync-windows-x64.zip server client
rm -rf client server

