// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/cmd/client/csgosync.go
	Language:	Go 1.15
	Dev Env:	Linux 5.9

	This file creates the client binary for the csgo sync application.
*/

package main

import (
	"fmt"
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
		errs  []error
	)
	fmt.Println("csgo sync client")

	// Read in config
	viper.SetConfigFile("csgosync.yaml")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	// TODO: don't crash and burn on missing file... try from env to be 12 factor
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("failed to read config: ", err)
		config.Wait()
		os.Exit(1)
	}
	clientConfig := config.InitClientConfig()

	// Check certain config params are met
	if clientConfig.Uri == "" {
		err = clientConfig.GetUri()
		if err != nil {
			fmt.Println("failed to get uri", err)
			config.Wait()
			os.Exit(1)
		}
	}

	if clientConfig.Pass == "" {
		err = clientConfig.GetPass()
		if err != nil {
			fmt.Println("failed to get password", err)
			config.Wait()
			os.Exit(1)
		}
	}

	// Create the hash map of our files and send to server
	fmt.Println("generating hash map...")
	files.Files, errs = filelist.GenerateMap(clientConfig.MapPath)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println("error creating hash map:", err)
		}
		config.Wait()
		os.Exit(1)
	}

	fmt.Println("sending hashmap to server")
	resp, err := httpclient.SendServerHashes(clientConfig.Uri+"/csgosync", clientConfig.Pass, files)
	if err != nil {
		fmt.Println("failed to get files list from server ", err)
		fmt.Println(resp)
		config.Wait()
		os.Exit(1)
	}
	// Download the missing/different files from server (if any)
	if len(resp.Files) != 0 {
		fmt.Printf("downloading %d files from server...\n", len(resp.Files))
		httpclient.DownloadFiles(clientConfig.Uri, clientConfig.Pass, clientConfig.MapPath, resp.Files)
	} else {
		fmt.Println("nothing to do, already have server's maps")
	}
	config.Wait() // config package already handles user input, this prevents windturds from closing cmd
}
