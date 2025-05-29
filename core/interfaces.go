package gapBotApi

import (
	"github.com/amirimatin/gapBotApi/v2/params"
	"io"
	"os"
)

type Chattable interface {
	params() (params.Params, error)
	method() string
}

type Fileable interface {
	Chattable
	file() RequestFile
}

type RequestFile struct {
	Name string
	Type string
	Data RequestFileData
}

type RequestFileData interface {
	NeedsUpload() bool
	UploadData() (string, io.Reader, error)
	SendData() string
}

type FileReader struct {
	Name   string
	Reader io.Reader
}

func (fr FileReader) NeedsUpload() bool                      { return true }
func (fr FileReader) UploadData() (string, io.Reader, error) { return fr.Name, fr.Reader, nil }
func (fr FileReader) SendData() string                       { panic("FileReader must be uploaded") }

type FilePath string

func (fp FilePath) NeedsUpload() bool { return true }
func (fp FilePath) UploadData() (string, io.Reader, error) {
	f, err := os.Open(string(fp))
	if err != nil {
		return "", nil, err
	}
	return f.Name(), f, nil
}
func (fp FilePath) SendData() string { panic("FilePath must be uploaded") }
