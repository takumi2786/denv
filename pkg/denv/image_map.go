package denv

import (
	"encoding/json"
	"os"

	"github.com/m-mizutani/goerr/v2"
)

type ImageMapEntry struct {
	Identity   string `json:"identity"`
	ImageURI   string `json:"image_uri"`
	Option     string `json:"option"`
	Entrypoint string `json:"entrypoint"`
	Cmd        string `json:"cmd"`
	Shell      string `json:"shell"`
}

type ImageMapReader struct {
	entries map[string]ImageMapEntry
}

func NewImageMapReader() *ImageMapReader {
	return &ImageMapReader{}
}

// Read read file and parse it as list of ImageMapEntry
func (imr *ImageMapReader) Read(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var entries []ImageMapEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}

	imr.entries = make(map[string]ImageMapEntry)
	for _, entry := range entries {
		imr.entries[entry.Identity] = entry
	}
	return nil
}

// Loaded search ImageMapEntry by identity
func (imr *ImageMapReader) Loadded(identity string) (*ImageMapEntry, error) {
	data, ok := imr.entries[identity]
	if !ok {
		return nil, goerr.New("Identity is not found")
	}
	return &data, nil
}
