// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:	csgosync/cmd/server/csgosyncd.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file creates the server binary for the csgo sync application.
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kthomas422/csgosync/internal/logging"

	"github.com/kthomas422/csgosync/config"

	"github.com/spf13/viper"

	"github.com/kthomas422/csgosync/internal/httpserver"
)

// TODO: create "custom" http server with timeouts
func main() {
	var cs httpserver.CsgoSync
	fmt.Println("CSGO Sync Server!")

	// load config
	viper.SetConfigFile("csgosyncd.yaml")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		cs.L.Err("failed to get config file: ", err)
		os.Exit(1)
	}
	cs.C = config.InitServerConfig()

	// init logger
	cs.LogFile, err = logging.OpenLogFile(cs.C.LogFile)
	if err != nil {
		cs.LogFile = os.Stderr
	}
	cs.L = logging.Init(cs.C.LogLvl, cs.LogFile)
	cs.L.Info("Map Sync Server")
	defer func() {
		if cs.LogFile != os.Stdout || cs.LogFile != os.Stderr {
			if err := cs.LogFile.Close(); err != nil {
				log.Println("failed to close log file: ", err)
			}
		}
	}()

	if cs.C.Pass == "" {
		err = cs.C.GetPass()
		if err != nil {
			cs.L.Err("failed to get password: ", err)
			os.Exit(1)
		}
	}
	if cs.C.Port == "" {
		cs.L.Info("empty port number, defaulting to 8080")
		cs.C.Port = "8080"
	}

	cs.L.Debug(fmt.Sprintf("server config: %#v", cs.C))

	// Handler for serving map files
	// TODO add auth to file server
	http.Handle("/maps/", http.StripPrefix(
		"/maps/", http.FileServer(http.Dir(cs.C.MapPath))))

	// Handler for map hashes
	http.Handle("/csgosync", &cs)

	// Catchall handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ip := httpserver.GetRequestIp(r)
		cs.L.Info(fmt.Sprintf("ip: %v | hit non-existent endpoint: %s", ip, r.URL.String()))
		w.WriteHeader(http.StatusNotFound)
		_, err = w.Write([]byte("{ \"Message\": \"Not found\"}"))
		if err != nil {
			cs.L.Err("failed to write back to client: ", err)
		}

	})

	cs.L.Info("Serving on port ", cs.C.Port)
	s := http.Server{
		ReadTimeout:       time.Second * 10,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Minute * 10, // hopefully files don't take longer than 10 minutes to download
		IdleTimeout:       time.Second * 30,
		Addr:              ":" + cs.C.Port,
	}
	cs.L.Debug(fmt.Sprintf("ReadTimeout: %d", s.ReadTimeout))
	cs.L.Debug(fmt.Sprintf("ReadHeaderTimeout: %d", s.ReadHeaderTimeout))
	cs.L.Debug(fmt.Sprintf("WriteTimeout: %d", s.WriteTimeout))
	cs.L.Debug(fmt.Sprintf("IdleTimeout: %d", s.IdleTimeout))
	cs.L.Err(s.ListenAndServe())
}
