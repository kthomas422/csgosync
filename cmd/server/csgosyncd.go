// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/cmd/server/csgosyncd.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file creates the server binary for the csgo sync application.
*/

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kthomas422/csgosync/internal/filelist"
	"github.com/kthomas422/csgosync/internal/models"
	"github.com/kthomas422/csgosync/internal/myconst"
)

// TODO: add port config via command line
// TODO: config file? (for paths, port and password?)
func main() {
	log.Println("csgo sync server")
	files, err := filelist.GenerateMap(myconst.SteamDir)
	if err != nil {
		log.Fatal("couldn't load server maps", err)
	}

	http.Handle("/maps", http.FileServer(http.Dir(myconst.SteamDir)))

	http.HandleFunc("/csgosync", func(w http.ResponseWriter, r *http.Request) {
		var (
			bytes       []byte
			jsonBody    []byte
			err         error
			remoteFiles models.FileHashMap
			resp        models.FileResponse
		)
		defer r.Body.Close()
		switch r.Method { // switch to make it easier to add more methods later
		case http.MethodPost:
			log.Println("post")
			bytes, err = ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("{ \"Message\": \"Error reading request body\"}"))
				return
			}
			err = json.Unmarshal(bytes, &remoteFiles)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("{ \"Message\": \"Error parsing JSON\"}"))
				return
			}

			resp.Files = filelist.CompareMaps(files, remoteFiles.Files)

			jsonBody, err = json.Marshal(resp)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("{ \"Message\": \"Error creating response\"}"))
			}
			w.WriteHeader(http.StatusOK)
			w.Write(jsonBody)
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{ \"Message\": \"Not found\"}"))
		}
	})

	http.ListenAndServe(":8080", nil)
}
