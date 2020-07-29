// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/cmd/server/csgosyncd.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7y

	This file creates the server binary for the csgo sync application.
*/

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kthomas422/csgosync/internal/logging"

	"github.com/kthomas422/csgosync/internal/httpserver"

	"github.com/kthomas422/csgosync/internal/auth"

	"github.com/kthomas422/csgosync/internal/constants"
)

// TODO: config file (for paths, port and password)
// TODO: create "custom" http server with timeouts
func main() {
	var cs httpserver.CsgoSync
	cs.L = logging.Init(logging.DebugLvl, os.Stderr)
	cs.L.Info("Map Sync Server")

	err := auth.GetPass()
	if err != nil {
		cs.L.Err("failed to get password: ", err)
		os.Exit(1)
	}

	http.Handle("/maps/", http.StripPrefix(
		"/maps/", http.FileServer(http.Dir(constants.ServerMapDir))))
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

	cs.L.Info("Serving on port 8080")
	cs.L.Err(http.ListenAndServe(":8080", nil))
}
