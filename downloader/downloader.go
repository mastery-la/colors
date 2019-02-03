package downloader

import (
	"os"
	"path/filepath"
)

type Downloader struct {
	outputFolder string
}

func New(outputFolder string) (*Downloader, error) {
	p := new(Downloader)

	newpath := filepath.Join(".", outputFolder)
	os.MkdirAll(newpath, os.ModePerm)

	p.outputFolder = newpath

	return p, nil
}

// func (d *Downloader) SaveImage(url string) (err error) {

// }
