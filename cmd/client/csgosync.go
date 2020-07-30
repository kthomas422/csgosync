// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/cmd/client/csgosync.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file creates the client binary for the csgo sync application.
*/

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/kthomas422/csgosync/config"

	"github.com/kthomas422/csgosync/internal/filelist"
	"github.com/kthomas422/csgosync/internal/httpclient"
	"github.com/kthomas422/csgosync/internal/models"
)

func main() {
	var (
		files models.FileHashMap
		err   error
	)
	log.Println("csgo sync client")

	viper.SetConfigFile("csgosync.config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("failed to read config: ", err)
		os.Exit(1)
	}
	clientConfig := config.InitClientConfig()

	if clientConfig.Uri == "" {
		err = clientConfig.GetUri()
		if err != nil {
			fmt.Println("failed to get uri", err)
			os.Exit(1)
		}
	}

	if clientConfig.Pass == "" {
		err = clientConfig.GetPass()
		if err != nil {
			fmt.Println("failed to get password", err)
			os.Exit(1)
		}
	}

	fmt.Println("generating hash map...")
	files.Files, err = filelist.GenerateMap(clientConfig.MapPath)
	fmt.Println("sending hashmap to server")
	resp, err := httpclient.SendServerHashes(clientConfig.Uri+"/csgosync", clientConfig.Pass, files)
	if err != nil {
		log.Println("failed to get files list from server ", err)
		log.Println(resp)
		os.Exit(1)
	}
	if len(resp.Files) != 0 {
		fmt.Printf("downloading %d files from server...\n", len(resp.Files))
		httpclient.DownloadFiles(clientConfig.Uri, clientConfig.Pass, clientConfig.MapPath, resp.Files)
	} else {
		fmt.Println("nothing to do, already have server's maps")
	}
	config.Wait() // config package already handles user input, this prevents windturds from closing cmd
}
