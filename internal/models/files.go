package models

// Response contains the server response code and the list of files that are different
type FileResponse struct {
	Status  int      `json:"status"`
	FileMap []string `json:"files"`
}

// ClientFileHashMap contains the map of files with the value being the hash of the files
type ClientFileHashMap struct {
	Files map[string]string `json:"files"`
}

// ServerFileHashMap contains a map with the key being the filename and the value being
// a struct with the file hash and the contents of the file (compressed).
type ServerFileHashMap struct {
	Files map[string]struct {
		hash     string
		contents []byte
	}
}
