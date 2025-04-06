package denv

import (
	"reflect"
	"testing"
)

func TestImageMapReader_Read(t *testing.T) {
	imr := ImageMapReader{}

	err := imr.Read("../../resources/image_map_test.json")
	if err != nil {
		t.Errorf("Error reading JSON data: %v", err)
	}

	expectedEntries := map[string]ImageMapEntry{
		"ubuntu": {
			Identity:   "ubuntu",
			ImageURI:   "ubuntu:noble-20250127",
			Option:     "-p 8080:8080",
			Entrypoint: "",
			Cmd:        "",
			Shell:      "zsh",
		},
		"python": {
			Identity:   "python",
			ImageURI:   "python:3.10.12-slim-bullseye",
			Option:     "-p 8080:8080",
			Entrypoint: "",
			Cmd:        "",
			Shell:      "bash",
		},
	}

	if !reflect.DeepEqual(imr.entries, expectedEntries) {
		t.Errorf("Parsed entries do not match expected entries")
	}
}

func TestImageMapReader_Loadded(t *testing.T) {
	imr := ImageMapReader{}
	imr.entries = map[string]ImageMapEntry{
		"ubuntu": {
			Identity:   "ubuntu",
			ImageURI:   "takumi2786/ubuntu:22.04-jammy-amd64",
			Option:     "-p 8080:8080",
			Entrypoint: "",
			Cmd:        "",
		},
		"python": {
			Identity:   "python",
			ImageURI:   "python:3.10.12-slim-bullseye",
			Option:     "-p 8080:8080",
			Entrypoint: "",
			Cmd:        "",
		},
	}

	// Existing identity
	data, err := imr.Loadded("ubuntu")
	if err != nil {
		t.Errorf("Unexpected error for existing identity: %v", err)
	}
	if data == nil || data.Identity != "ubuntu" {
		t.Errorf("Unexpected data for existing identity")
	}

	// Non-existing identity
	data, err = imr.Loadded("invalid")
	if err == nil {
		t.Errorf("Expected error for non-existing identity but got nil")
	}
	if data != nil {
		t.Errorf("Unexpected data for non-existing identity")
	}
}
