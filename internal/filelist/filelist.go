// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/filelist/filelist.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file contains the methods for file listing/hashing for the csgo sync application.
*/

package filelist

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// loadFiles loads the list of files in the directory
func loadFiles(dir string) (files []string, err error) {
	if len(dir) < 1 {
		return nil, errors.New("no directory passed in")
	}
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}
	filesList, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("couldn't load any files")
		return nil, err
	}
	for _, file := range filesList {
		if file.Mode().IsRegular() && file.Size() > 0 {
			files = append(files, dir+file.Name())
		}
	}
	return files, nil
}

// hashFiles takes a list of files and computes the hashes of the files
func hashFiles(files []string) (hashes []string, err error) {
	const bufSize = 16777216 // 16.7MB
	var (
		wg        = new(sync.WaitGroup)
		mux       = new(sync.Mutex)
		fileChunk = 50 // files to hash per thread
	)
	for len(files) > 0 {
		if len(files) < fileChunk {
			fileChunk = len(files)
		}
		wg.Add(1)
		// TODO send errors to main thread to handle
		go func(files []string) {
			var (
				err error
				f   *os.File
			)
			defer wg.Done()
			hasher := sha1.New()
			for _, file := range files {
				f, err = os.Open(file)
				if err != nil {
					if f != nil {
						f.Close()
					}
					fmt.Printf("error: skipping file %s\n%s\n", file, err)
					continue
				}
				buf := make([]byte, bufSize)
				for err != io.EOF {
					_, err = f.Read(buf)
					if err != nil && err != io.EOF {
						fmt.Println("error reading file: ", err)
						continue
					}
					hasher.Write(buf)
				}
				err = f.Close()
				if err != nil {
					fmt.Println("error closing file: ", err)
				}
				mux.Lock()
				hashes = append(hashes, hex.EncodeToString(hasher.Sum(nil)))
				mux.Unlock()
				hasher.Reset()
			}
		}(files[:fileChunk])

		files = files[fileChunk:]
	}
	wg.Wait()
	return hashes, nil
}

// GenerateMap makes a map with the list of files from the directory and the file's hash
func GenerateMap(dir string) (map[string]string, error) {
	maps := make(map[string]string)
	files, err := loadFiles(dir)
	if len(files) == 0 {
		return maps, nil // return empty map since no files
	}
	if err != nil {
		return nil, err
	}
	hashes, err := hashFiles(files)
	if err != nil {
		return nil, err
	}
	for i, file := range files {
		maps[filepath.Base(file)] = hashes[i]
	}
	return maps, err
}

func CompareMaps(serverFiles, clientFiles map[string]string) (delta []string) {
	for serverFile, serverHash := range serverFiles {
		if _, ok := clientFiles[serverFile]; !ok { // check if the client has the file
			delta = append(delta, serverFile)
		} else if serverHash != clientFiles[serverFile] { // check if same file
			delta = append(delta, serverFile)
		}
	}
	return delta
}
