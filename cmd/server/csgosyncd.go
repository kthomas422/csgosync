// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:	csgosync/cmd/server/csgosyncd.go
	Language:	Go 1.15
	Dev Env:	Linux 5.9

	This file creates the server binary for the csgo sync application.
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kthomas422/csgosync/internal/csgolog"

	"github.com/kthomas422/csgosync/internal/filelist"

	"github.com/kthomas422/csgosync/config"

	"github.com/spf13/viper"

	"github.com/kthomas422/csgosync/internal/httpserver"
)

func main() {
	var cs httpserver.CsgoSync
	fmt.Println("CSGO Sync Server!")

	// load config
	viper.SetConfigFile("csgosyncd.yaml")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	// TODO: don't crash and burn on missing file... try from env to be 12 factor
	if err != nil {
		log.Fatal("failed to get config file: ", err)
	}
	cs.C = config.InitServerConfig()

	// init logger
	cs.L, err = csgolog.InitLogger(cs.C.LogFile)
	if err != nil {
		log.Fatalf("could not start logger: %v\n", err)
	}
	cs.L.Simple("CSGO Sync Server!")

	// Make sure we close the log file
	defer func() {
		if err := cs.L.Close(); err != nil {
			log.Printf("error closing log file: %v\n", err)
		}
	}()

	// Verify required config values are present
	if cs.C.Pass == "" {
		cs.L.Err("failed to get password: ", err)
		os.Exit(1)
	}
	if cs.C.Port == "" {
		cs.L.Simple("empty port number, defaulting to 8080")
		cs.C.Port = "8080"
	}

	if err := cs.L.Config(*cs.C); err != nil {
		log.Fatalf("could not write to logger: %v", err)
	}

	// Generate hash map (and regenerate every so often)
	// Timing how long it takes to get the map as well
	go func() {
		var errs []error
		for {
			cs.L.Simple("generating hash map")
			start := time.Now()
			cs.HashMap, errs = filelist.GenerateMap(cs.C.MapPath)
			if len(errs) > 0 {
				for _, err := range errs {
					cs.L.Err("failed getting hashmap", err)
				}
				os.Exit(1) // crash and burn since we can't make maps
			}
			elapsed := time.Since(start)
			cs.L.Simple(fmt.Sprintf("hash map generated in %v", elapsed))
			if err != nil {
				cs.L.Err("couldn't load server maps: ", err)
			}
			cs.L.Simple(fmt.Sprintf("files list: %v", cs.HashMap))
			time.Sleep(time.Hour * 24 * 7) // regenerate hash map every week TODO: set this as config
		}
	}()

	// Handler for serving map files
	// TODO add auth to file server
	http.Handle("/maps/", http.StripPrefix(
		"/maps/", http.FileServer(http.Dir(cs.C.MapPath))))

	// Handler for map hashes
	http.Handle("/csgosync", &cs)

	// Catchall handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cs.L.WebRequest(r)
		w.WriteHeader(http.StatusNotFound)
		_, err = w.Write([]byte("{ \"Message\": \"Not found\"}"))
		if err != nil {
			cs.L.Err("failed to write back to client: ", err)
		}
	})

	// Create web server and run it
	s := http.Server{
		ReadTimeout:       time.Second * 10,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Minute * 10, // hopefully files don't take longer than 10 minutes to download
		IdleTimeout:       time.Second * 30,
		Addr:              ":" + cs.C.Port,
	}

	cs.L.Err("server shutdown response:", s.ListenAndServe())
}
