package filelist

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"log"
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
	hasher := sha256.New()
	for _, file := range files {
		f, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		hasher.Write(f)
		hashes = append(hashes, hex.EncodeToString(hasher.Sum(nil)))
		hasher.Reset()
	}
	return hashes, nil
}

// GenerateMap makes a map with the list of files from the directory and the file's hash
func GenerateMap(dir string) (map[string]string, error) {
	maps := make(map[string]string)
	files, err := loadFiles(dir)
	if err != nil {
		return nil, err
	}
	hashes, err := hashFiles(files)
	if err != nil {
		return nil, err
	}
	for i, file := range files {
		maps[file] = hashes[i]
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
