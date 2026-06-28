package files

import "path/filepath"

type Video struct {
	File
}

func NewVideoFromPath(path string) (*Video, error) {
	return &Video{File: File{path: path, name: filepath.Base(path)}}, nil
}

func NewVideoFromBytes(data []byte, name string) *Video {
	return &Video{File: File{data: data, name: name}}
}
