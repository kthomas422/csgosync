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
	"strings"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/kthomas422/csgosync/internal/concurrency"

	"github.com/kthomas422/csgosync/internal/models"
)

const (
	timeOut                = 3600 // may take an hour to finish download
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
func downloadFile(uri, file, mapDir string, concOH *concurrency.OverHead) {
	var counter *WriteCounter
	defer concOH.Wg.Done() // Signal that download is done

	// get data
	concOH.HttpSem <- concurrency.Token{} // "take token"
	resp, err := httpClient.client.Get(uri + "/maps/" + file)
	if err != nil {
		<-concOH.HttpSem // release token
		return
	}
	defer resp.Body.Close()
	<-concOH.HttpSem

	// create file
	concOH.FileSem <- concurrency.Token{}
	out, err := os.Create(filepath.Join(mapDir, file) + ".tmp")
	if err != nil {
		<-concOH.FileSem // release token
		fmt.Println("failed to created file: ", file)
		return
	}
	defer out.Close()

	// 	https://golangcode.com/download-a-file-with-progress/
	// Create our progress reporter and pass it to be used alongside our writer
	if len(file) > 16 {
		counter = &WriteCounter{FileName: file[:16]} // truncate file name to 16 chars
	} else {
		counter = &WriteCounter{FileName: file} // truncate file name to 16 chars
	}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		<-concOH.FileSem // release token
		fmt.Println("failed to write to file")
		err = out.Close()
		if err != nil {
			fmt.Println("failed to close file")
		}
		return
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	// Close the file without defer so it can happen before Rename()
	err = out.Close()
	if err != nil {
		<-concOH.FileSem // release token
		fmt.Println("failed to close file")
	}

	// write to tmp file in case it failed... now rename to the "real" name
	if err = os.Rename(
		filepath.Join(mapDir, file)+".tmp", filepath.Join(mapDir, file)); err != nil {
		<-concOH.FileSem // release token
		fmt.Println("failed to remove tmp file")
		return
	}

	<-concOH.FileSem // release token
	return
}

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total    uint64
	FileName string
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// TODO: bug where filenames don't really match progress bar
func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 80))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading %s %s complete", wc.FileName, humanize.Bytes(wc.Total))
}
