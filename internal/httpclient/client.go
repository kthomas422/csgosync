package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

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

func SendServerHashes(uri string, body models.ClientFileHashMap) (*models.FileResponse, error) {
	var filesResp = new(models.FileResponse)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(filesResp)
	if err != nil {
		return nil, err
	}
	return filesResp, nil
}

func DownloadFiles(uri string, files []string) {
	var (
		wg            sync.WaitGroup
		webSemaphore  = make(chan struct{}, maxConcurrentDownloads)
		fileSemaphore = make(chan struct{}, maxOpenFiles)
	)
	for _, file := range files {
		go downloadFile(uri, file, webSemaphore, fileSemaphore, &wg)
	}
	wg.Wait()
}

// download the file from the url onto local drive with same name
func downloadFile(uri, file string, webS chan struct{}, fileS chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done() // Signal that download is done

	// get data
	webS <- struct{}{} // "take token"
	resp, err := httpClient.client.Get(uri + file)
	if err != nil {
		log.Println("error downloading file:", uri, file)
		log.Println(err)
		<-webS // release token
		return
	}
	defer resp.Body.Close()
	<-webS

	// create file
	fileS <- struct{}{}
	out, err := os.Create(file)
	if err != nil {
		log.Println("error creating file:", file)
		log.Println(err)
		<-fileS // release token
		return
	}
	defer out.Close()

	// write contents to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println("error writing to file:", file)
		log.Println(err)
		<-fileS // release token
		return
	}
	<-fileS // release token

	log.Println(" - SUCCESS\t", file)
	return
}
