package files

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var AllowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".bmp":  true,
}

type BaseFile interface {
	Read() ([]byte, error)
	Size() (int64, error)
	Name() string
}

type File struct {
	path string
	name string
	data []byte
	url  string
}

func NewFileFromPath(path string) (*File, error) {
	name := filepath.Base(path)
	return &File{path: path, name: name}, nil
}

func NewFileFromBytes(data []byte, name string) *File {
	return &File{data: data, name: name}
}

func (f *File) Read() ([]byte, error) {
	if f.data != nil {
		return f.data, nil
	}
	if f.path != "" {
		return os.ReadFile(f.path)
	}
	if f.url != "" {
		resp, err := http.Get(f.url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return io.ReadAll(resp.Body)
	}
	return nil, fmt.Errorf("no data source for file")
}

func (f *File) Size() (int64, error) {
	if f.data != nil {
		return int64(len(f.data)), nil
	}
	if f.path != "" {
		info, err := os.Stat(f.path)
		if err != nil {
			return 0, err
		}
		return info.Size(), nil
	}
	return 0, fmt.Errorf("cannot determine file size")
}

func (f *File) Name() string {
	return f.name
}


