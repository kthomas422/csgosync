// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/cmd/server/csgosyncd.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file creates the server binary for the csgo sync application.
*/

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kthomas422/csgosync/config"

	"github.com/spf13/viper"

	"github.com/kthomas422/csgosync/internal/logging"

	"github.com/kthomas422/csgosync/internal/httpserver"
)

// TODO: create "custom" http server with timeouts
func main() {
	var cs httpserver.CsgoSync
	cs.L = logging.Init(logging.DebugLvl, os.Stderr)
	cs.L.Info("Map Sync Server")
	viper.SetConfigFile("csgosyncd.config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		cs.L.Err("failed to get config file: ", err)
		os.Exit(1)
	}
	cs.C = config.InitServerConfig()

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

	http.Handle("/maps/", http.StripPrefix(
		"/maps/", http.FileServer(http.Dir(cs.C.MapPath))))
	http.Handle("/csgosync", &cs)
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
	cs.L.Err(http.ListenAndServe(":"+cs.C.Port, nil))
}
