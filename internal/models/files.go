package models

// Response contains the server response code and the list of files that are different
type FileResponse struct {
	Files []string `json:"files"`
}

// ClientFileHashMap contains the map of files with the value being the hash of the files
type FileHashMap struct {
	Files map[string]string `json:"files"`
}
