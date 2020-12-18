// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/httpclient/client.go
	Language:	Go 1.15
	Dev Env:	Linux 5.9

	This file contains the structs and methods for the http client for the csgo sync application.
*/

package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kthomas422/csgosync/internal/concurrency"

	"github.com/kthomas422/csgosync/internal/models"
)

const (
	timeOut                = 3600 // 1 hour to finish downloads, may need to increase?
	maxConcurrentDownloads = 64   // Limit download to 64 files at a time
	maxOpenFiles           = 64   // Limit to 64 files open at once
)

// Wrapper for builtin http client
var httpClient struct {
	client *http.Client
}

// Create http client on startup
func init() {
	httpClient.client = &http.Client{
		Timeout: time.Duration(timeOut) * time.Second,
	}
}

// SendServerHashes sends the server the hashmap and returns a list of files that were missing or different.
func SendServerHashes(uri, pass string, body models.FileHashMap) (*models.FileResponse, error) {
	var filesResp = new(models.FileResponse)

	// prepare request body
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("pass", pass)

	// send request
	resp, err := httpClient.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// read response
	defer resp.Body.Close()
	respContents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return filesResp, fmt.Errorf("bad http status: %s", resp.Status)
	}
	err = json.Unmarshal(respContents, &filesResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return filesResp, nil
}

// Downloads the files from the server that we need
func DownloadFiles(uri, pass, mapDir string, files []string) {
	var (
		concOH = concurrency.InitOH(maxConcurrentDownloads, maxOpenFiles)
	)
	for _, file := range files {
		concOH.Wg.Add(1)
		go downloadFile(uri, file, mapDir, concOH)
	}
	concOH.Wg.Wait()
}

// download the file from the url onto local drive with same name
// TODO: find a pretty way to print progress bar
// TODO: cleanup the 10000000 error branches
func downloadFile(uri, file, mapDir string, concOH *concurrency.OverHead) {
	defer concOH.Wg.Done() // Signal that download is done

	// get data
	concOH.HttpSem <- concurrency.Token{} // "take token"
	resp, err := httpClient.client.Get(uri + "/maps/" + file)
	<-concOH.HttpSem
	if err != nil {
		fmt.Printf("failed to download file: %v\n", err)
		return
	}

	// create tmp file
	concOH.FileSem <- concurrency.Token{}
	out, err := os.Create(filepath.Join(mapDir, file) + ".tmp")
	if err != nil {
		<-concOH.FileSem // release token
		fmt.Printf("failed to created file: %s, error: %v\n", file, err)
		if out != nil {
			err = out.Close()
			if err != nil {
				fmt.Printf("failed to close file: %v", err)
			}
		}
		return
	}

	// Copy the file from the server into our tmp file
	if _, err = io.Copy(out, resp.Body); err != nil {
		<-concOH.FileSem // release token
		fmt.Printf("failed to write to file: %s, error: %v\n", file, err)
		err = out.Close()
		if err != nil {
			fmt.Printf("failed to close file: %s, error: %v\n", file, err)
		}
		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("unable to close server response body: %v\n", err)
		}
		return
	}
	err = resp.Body.Close()
	if err != nil {
		fmt.Printf("unable to close server response body: %v\n", err)
	}

	err = out.Close()
	if err != nil {
		<-concOH.FileSem // release token
		fmt.Printf("failed to close file: %s error: %v", file, err)
	}

	// wrote to tmp file in case it failed... now rename to the "real" name
	if err = os.Rename(
		filepath.Join(mapDir, file)+".tmp", filepath.Join(mapDir, file)); err != nil {
		<-concOH.FileSem // release token
		fmt.Println("failed to remove tmp file")
		return
	}
	<-concOH.FileSem // release token
	fmt.Printf("file: %s downloaded\n", file)
}
