#!/bin/sh
mkdir client server
cp csgosync.yaml.example client/csgosync.yaml
cp csgosyncd.yaml.example server/csgosyncd.yaml

go build -o client/csgosync cmd/client/csgosync.go
go build -o server/csgosyncd cmd/server/csgosyncd.go
tar -czvf csgosync-linux-x64.tgz server client
rm client/csgosync server/csgosyncd

env GOOS=windows GOARCH=amd64 go build -o client/csgosync.exe cmd/client/csgosync.go
env GOOS=windows GOARCH=amd64 go build -o server/csgosyncd.exd cmd/server/csgosyncd.go
zip -r csgosync-windows-x64.zip server client
rm -rf client server

