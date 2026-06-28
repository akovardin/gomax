package files

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Photo struct {
	File
}

func NewPhotoFromPath(path string) (*Photo, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if !AllowedImageExtensions[ext] {
		return nil, fmt.Errorf("unsupported photo extension: %s", ext)
	}
	return &Photo{File: File{path: path, name: filepath.Base(path)}}, nil
}

func NewPhotoFromBytes(data []byte, name string) *Photo {
	return &Photo{File: File{data: data, name: name}}
}

func NewPhotoFromURL(url string, name string) *Photo {
	return &Photo{File: File{url: url, name: name}}
}
