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
}

type ImageMapReader struct {
	entries map[string]ImageMapEntry
}

func NewImageMapReader() *ImageMapReader {
	return &ImageMapReader{}
}

// Read は、ファイルから設定を読み込む。
func (imr *ImageMapReader) Read(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var entries []ImageMapEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}

	// Identity -> データの形式のMAPで保存
	imr.entries = make(map[string]ImageMapEntry)
	for _, entry := range entries {
		imr.entries[entry.Identity] = entry
	}
	return nil
}

// Loadded は、Identityでentriesを検索し、返します。
func (imr *ImageMapReader) Loadded(identity string) (*ImageMapEntry, error) {
	data, ok := imr.entries[identity]
	if !ok {
		return nil, goerr.New("Identity is not found")
	}
	return &data, nil
}
