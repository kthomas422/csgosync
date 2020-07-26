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

	"github.com/kthomas422/csgosync/internal/auth"

	"github.com/kthomas422/csgosync/internal/constants"
	"github.com/kthomas422/csgosync/internal/filelist"
	"github.com/kthomas422/csgosync/internal/models"
)

// TODO: add port config via command line
// TODO: config file? (for paths, port and password?)
// TODO: add custom logging package
func main() {
	log.Println("csgo sync server")
	files, err := filelist.GenerateMap(constants.ServerMapDir)
	if err != nil {
		log.Fatal("couldn't load server maps: ", err)
	}

	err = auth.GetPass()
	if err != nil {
		log.Fatal("failed to get password: ", err)
	}

	http.Handle("/maps", http.FileServer(http.Dir(constants.ServerMapDir)))

	http.HandleFunc("/csgosync", func(w http.ResponseWriter, r *http.Request) {
		var (
			bytes       []byte
			jsonBody    []byte
			err         error
			remoteFiles models.FileHashMap
			resp        models.FileResponse
		)
		log.Println("header: ", r.Header)
		if pass := r.Header.Get("Pass"); pass != "" {
			if pass != auth.Password() {
				log.Println("unauthorized bad pass: ", r.Header.Get("pass"))
				unAuth(w)
				return
			}
		} else {
			log.Println("unauthorized no password")
			unAuth(w)
			return
		}
		defer r.Body.Close()
		switch r.Method { // switch to make it easier to add more methods later
		case http.MethodPost:
			log.Println("post")
			bytes, err = ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusUnprocessableEntity)
				_, err = w.Write([]byte("{ \"Message\": \"Error reading request body\"}"))
				if err != nil {
					log.Println()
				}
				return
			}
			err = json.Unmarshal(bytes, &remoteFiles)
			if err != nil {
				log.Println("can't unmarshal json ", err)
				w.WriteHeader(http.StatusUnprocessableEntity)
				_, err = w.Write([]byte("{ \"Message\": \"Error parsing JSON\"}"))
				if err != nil {
					log.Println()
				}
				return
			}

			resp.Files = filelist.CompareMaps(files, remoteFiles.Files)

			jsonBody, err = json.Marshal(resp)
			if err != nil {
				log.Println("can't marshal json ", err)
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(w.Write([]byte("{ \"Message\": \"Error creating response\"}")))
			}
			log.Println("success!")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(jsonBody)
			if err != nil {
				log.Println(err)
			}
		default:
			log.Println("not found")
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write([]byte("{ \"Message\": \"Not found\"}"))
			if err != nil {
				log.Println(err)
			}
		}
	})

	log.Println("Serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func unAuth(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write([]byte("{ \"Message\": \"Unauthorized\"}"))
	if err != nil {
		log.Println(err)
	}
}
