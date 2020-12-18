// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/filelist/filelist.go
	Language:	Go 1.15
	Dev Env:	Linux 5.9

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
	"os"
	"path/filepath"
	"sync"
)

// loadFiles loads the list of files in the directory
func loadFiles(dir string) (files []string, err error) {
	if len(dir) < 1 {
		return nil, errors.New("no directory passed in")
	}
	// need to add slash at the end if there isn't one
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}
	filesList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("couldn't read directory: %w", err)
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
		// concurrently hash a chunk of files
		go func(files []string) {
			var (
				err error
				f   *os.File
			)
			defer wg.Done()
			hasher := sha1.New()
			// iterate over files in chunk
			// TODO: this needs a lot of cleaning up
			for _, file := range files {
				f, err = os.Open(file)
				if err != nil {
					if f != nil {
						f.Close()
					}
					fmt.Printf("error: skipping file %s\n%s\n", file, err)
					continue
				}
				// consume the file in chunks (was way more fun to read the whole file at once but will
				// fill up the ram on a micro aws instance).
				buf := make([]byte, bufSize)
				for err != io.EOF {
					_, err = f.Read(buf)
					if err != nil && err != io.EOF {
						fmt.Println("error reading file: ", err)
						continue
					}
					hasher.Write(buf)
				}
				if err = f.Close(); err != nil {
					fmt.Println("error closing file: ", err)
				}
				// Lock the hashes so another thread can't write to them, write our hashes to it and unlock
				mux.Lock()
				hashes = append(hashes, hex.EncodeToString(hasher.Sum(nil)))
				mux.Unlock()
				hasher.Reset() // reset the hasher for the next file/hash
			}
		}(files[:fileChunk])

		files = files[fileChunk:]
	}
	// wait for all the spawned threads to hash their files
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
	// have a list of file and a list of hashes... cram them into a map
	for i, file := range files {
		maps[filepath.Base(file)] = hashes[i]
	}
	return maps, err
}

// Takes in 2 hashmaps and returns a slice of the filenames for the different ones
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
