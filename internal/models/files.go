// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/models/files.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file contains the json models for the csgo sync application.
*/

package models

// Response contains the server response code and the list of files that are different
type FileResponse struct {
	Files []string `json:"files"`
}

// ClientFileHashMap contains the map of files with the value being the hash of the files
type FileHashMap struct {
	Files map[string]string `json:"files"`
}
