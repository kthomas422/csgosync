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
func hashFiles(files []string) ([]string, []error) {
	const bufSize = 16777216 // 16.7MB
	var (
		hasher  = sha1.New()
		hashes  = make([]string, len(files))
		f       *os.File
		err     error
		fileErr error // file specific error
		errs    []error
	)
	// iterate over files and hash them
	// TODO: this needs a lot of cleaning up
	for i, file := range files {
		f, err = os.Open(file)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to open file: %s: %w", file, err))
			if f != nil {
				if err = f.Close(); err != nil {
					errs = append(errs, fmt.Errorf("failed to close file: %s, %w", file, err))
				}
			}
			continue
		}

		// consume the file in chunks (was way more fun to read the whole file at once but will
		// fill up the ram on a micro aws instance).
		buf := make([]byte, bufSize)
		for fileErr != io.EOF {
			_, fileErr = f.Read(buf)
			if fileErr != nil && fileErr != io.EOF {
				errs = append(errs, fmt.Errorf("error reading file: %s: %w", file, fileErr))
				hasher.Reset()
				continue
			}
			_, err = hasher.Write(buf)
			if err != nil {
				errs = append(errs, fmt.Errorf("could not put bytes in hasher: %s: %w", file, err))
			}
		}
		if err = f.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error closing file: %s: %w ", file, err))
		}
		hashes[i] = hex.EncodeToString(hasher.Sum(nil))
		hasher.Reset() // reset the hasher for the next file
	}

	return hashes, errs
}

// GenerateMap makes a map with the list of files from the directory and the file's hash
func GenerateMap(dir string) (map[string]string, []error) {
	var (
		maps = make(map[string]string)
		errs []error
	)
	files, err := loadFiles(dir)
	if len(files) == 0 {
		return maps, nil // return empty map since no files
	}
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to get list of files: %w", err))
		return nil, errs
	}
	hashes, hashErrs := hashFiles(files)
	if len(hashErrs) > 0 {
		errs = append(errs, hashErrs...)
		return nil, errs
	}
	// have a list of file and a list of hashes... cram them into a map
	for i, file := range files {
		maps[filepath.Base(file)] = hashes[i]
	}
	return maps, nil
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
