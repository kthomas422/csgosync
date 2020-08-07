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
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kthomas422/csgosync/internal/concurrency"

	"github.com/kthomas422/csgosync/internal/models"
)

const (
	timeOut                = 3600 // 1 hour to finish downloads, may need to increase?
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
		return nil, err
	}

	req.Header.Add("pass", pass)
	resp, err := httpClient.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respContents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return filesResp, errors.New("bad http status")
	}
	err = json.Unmarshal(respContents, &filesResp)
	//err = json.NewDecoder(resp.Body).Decode(&filesResp)
	if err != nil {
		return nil, err
	}

	return filesResp, nil
}

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
	if err != nil {
		<-concOH.HttpSem // release token
		return
	}
	<-concOH.HttpSem

	// create file
	concOH.FileSem <- concurrency.Token{}
	out, err := os.Create(filepath.Join(mapDir, file) + ".tmp")
	if err != nil {
		<-concOH.FileSem // release token
		fmt.Println("failed to created file: ", file)
		if out != nil {
			err = out.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
		return
	}

	if _, err = io.Copy(out, resp.Body); err != nil {
		<-concOH.FileSem // release token
		fmt.Println("failed to write to file")
		err = out.Close()
		if err != nil {
			fmt.Println("failed to close file")
		}
		err = out.Close()
		if err != nil {
			fmt.Println(err)
		}
		err = resp.Body.Close()
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	err = resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	err = out.Close()
	if err != nil {
		<-concOH.FileSem // release token
		fmt.Println("failed to close file")
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
