// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/httpclient/client.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file contains the structs and methods for the http client for the csgo sync application.
*/

package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kthomas422/csgosync/internal/constants"

	"github.com/kthomas422/csgosync/internal/concurrency"

	"github.com/kthomas422/csgosync/internal/models"
)

const (
	timeOut                = 60
	maxConcurrentDownloads = 64
	maxOpenFiles           = 64
)

var httpClient struct {
	client *http.Client
}

func init() {
	httpClient.client = &http.Client{
		Timeout: time.Duration(timeOut) * time.Second,
	}
}

func SendServerHashes(uri, pass string, body models.FileHashMap) (*models.FileResponse, error) {
	var filesResp = new(models.FileResponse)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("fail to create new request")
		return nil, err
	}

	req.Header.Add("pass", pass)
	resp, err := httpClient.client.Do(req)
	if err != nil {
		log.Println("failed to do request")
		return nil, err
	}
	defer resp.Body.Close()
	respContents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("failed to read response")
		os.Exit(1)
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("bad server response ", resp.Status)
		fmt.Println(string(respContents))
		return filesResp, errors.New("bad http status")
	}
	err = json.Unmarshal(respContents, &filesResp)
	//err = json.NewDecoder(resp.Body).Decode(&filesResp)
	if err != nil {
		log.Println("failed to decode json")
		return nil, err
	}

	return filesResp, nil
}

func DownloadFiles(uri, pass string, files []string) {
	var (
		concOH = concurrency.InitOH(maxConcurrentDownloads, maxOpenFiles)
	)
	for _, file := range files {
		concOH.Wg.Add(1)
		go downloadFile(uri, file, concOH)
	}
	concOH.Wg.Wait()
}

// download the file from the url onto local drive with same name
func downloadFile(uri, file string, concOH *concurrency.OverHead) {
	defer concOH.Wg.Done() // Signal that download is done

	fmt.Println("inside download file")
	log.Println("["+uri+"]", "<"+file+">")

	// get data
	concOH.HttpSem <- concurrency.Token{} // "take token"
	resp, err := httpClient.client.Get(uri + "/maps/" + file)
	if err != nil {
		log.Println("error downloading file:", uri, file)
		log.Println(err)
		<-concOH.HttpSem // release token
		return
	}
	defer resp.Body.Close()
	<-concOH.HttpSem

	// create file
	concOH.FileSem <- concurrency.Token{}
	out, err := os.Create(filepath.Join(constants.ClientMapDir, file))
	if err != nil {
		log.Println("error creating file:", file)
		log.Println(err)
		<-concOH.FileSem // release token
		return
	}
	defer out.Close()

	// write contents to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println("error writing to file:", file)
		log.Println(err)
		<-concOH.FileSem // release token
		return
	}
	<-concOH.FileSem // release token

	log.Println(" - SUCCESS\t", file)
	return
}
