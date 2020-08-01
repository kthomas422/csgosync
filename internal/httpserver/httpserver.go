// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/httpserver/httpserver.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file contains the methods and functions for the http server.
*/

package httpserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kthomas422/csgosync/config"

	"github.com/kthomas422/csgosync/internal/logging"

	"github.com/kthomas422/csgosync/internal/filelist"
	"github.com/kthomas422/csgosync/internal/models"
)

type CsgoSync struct {
	L *logging.Log
	C *config.ServerConfig
}

func (cs *CsgoSync) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		files       map[string]string
		bytes       []byte
		jsonBody    []byte
		err         error
		remoteFiles models.FileHashMap
		resp        models.FileResponse
		ip          = GetRequestIp(r)
	)
	files, err = filelist.GenerateMap(cs.C.MapPath)
	if err != nil {
		cs.L.Err("couldn't load server maps: ", err)
		cs.L.Info(fmt.Sprintf("ip: %v | couldn't load server maps", ip))
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("{ \"Message\": \"Internal Error\"}"))
		if err != nil {
			cs.L.Err("failed to write back to client: ", err)
		}
		return
	}
	cs.L.Debug("header: ", r.Header)
	if pass := r.Header.Get("Pass"); pass != "" {
		if pass != cs.C.Pass {
			cs.L.Info(fmt.Sprintf("ip: %v | unauthorized: bad pass: %s", ip, r.Header.Get("Pass")))
			err = unAuth(w)
			if err != nil {
				cs.L.Err("failed to write back to client: ", err)
			}
			return
		}
	} else {
		cs.L.Info(fmt.Sprintf("ip: %v | unauthorized: no password", ip))
		err = unAuth(w)
		if err != nil {
			cs.L.Err("failed to write back to client: ", err)
		}
		return
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			cs.L.Err("failed to close response body: ", err)
		}
	}()
	switch r.Method { // switch to make it easier to add more methods later
	case http.MethodPost:
		cs.L.Info("request from: ", ip)
		bytes, err = ioutil.ReadAll(r.Body)
		if err != nil {
			cs.L.Err("failed to read request body: ", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, err = w.Write([]byte("{ \"Message\": \"Error reading request body\"}"))
			if err != nil {
				cs.L.Err("failed to write back to client: ", err)
			}
			return
		}
		err = json.Unmarshal(bytes, &remoteFiles)
		if err != nil {
			cs.L.Err("can't unmarshal json ", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, err = w.Write([]byte("{ \"Message\": \"Error parsing JSON\"}"))
			if err != nil {
				cs.L.Err("failed to write back to client: ", err)
			}
			return
		}

		resp.Files = filelist.CompareMaps(files, remoteFiles.Files)

		jsonBody, err = json.Marshal(resp)
		if err != nil {
			cs.L.Err("can't marshal json ", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte("{ \"Message\": \"Error creating response\"}"))
			if err != nil {
				cs.L.Err("failed to write back to client: ", err)
			}
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonBody)
		if err != nil {
			cs.L.Err("failed to write back to client: ", err)
		}
	default:
		cs.L.Info(fmt.Sprintf("ip: %v | hit non-existent endpoint: %s", ip, r.URL.String()))
		w.WriteHeader(http.StatusNotFound)
		_, err = w.Write([]byte("{ \"Message\": \"Not found\"}"))
		if err != nil {
			log.Println(err)
		}
	}
	cs.L.Info(fmt.Sprintf("ip: %v | successfully sent map delta (%d)", ip, len(resp.Files)))
}

func unAuth(w http.ResponseWriter) (err error) {
	w.WriteHeader(http.StatusUnauthorized)
	_, err = w.Write([]byte("{ \"Message\": \"Unauthorized\"}"))
	return
}

func GetRequestIp(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
