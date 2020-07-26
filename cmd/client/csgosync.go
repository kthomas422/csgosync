// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/cmd/client/csgosync.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file creates the client binary for the csgo sync application.
*/

package main

import (
	"log"
	"os"

	"github.com/kthomas422/csgosync/internal/auth"
	"github.com/kthomas422/csgosync/internal/constants"
	"github.com/kthomas422/csgosync/internal/filelist"
	"github.com/kthomas422/csgosync/internal/httpclient"
	"github.com/kthomas422/csgosync/internal/models"
)

func main() {
	var (
		files models.FileHashMap
	)
	log.Println("csgo sync client")
	err := auth.GetUri()
	if err != nil {
		log.Fatal("failed to get uri", err)
	}

	err = auth.GetPass()
	if err != nil {
		log.Fatal("failed to get password", err)
	}

	log.Println("Generating file hash map")
	files.Files, err = filelist.GenerateMap(constants.ClientMapDir)
	log.Println("Sending hashes to server")
	resp, err := httpclient.SendServerHashes(auth.Uri()+"/csgosync", auth.Password(), files)
	if err != nil {
		log.Println("failed to get files list from server ", err)
		log.Println(resp)
		os.Exit(1)
	}
	log.Println("downloading files")
	log.Println(resp.Files)
	httpclient.DownloadFiles(auth.Uri(), auth.Password(), resp.Files)
}
