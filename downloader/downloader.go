package downloader

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Downloader holds neccessary information to use throughout downloading
type Downloader struct {
	OutputFolder string
}

// New creates a new instance of Downloader
func New(outputFolder string) (*Downloader, error) {
	p := new(Downloader)

	newpath := filepath.Join(".", outputFolder)
	os.MkdirAll(newpath, os.ModePerm)

	p.OutputFolder = newpath

	return p, nil
}

// SaveImage downloads the image from the given URL and saves to file
func (d *Downloader) SaveImage(fromURL string, toFile string) error {
	path := filepath.Join(d.OutputFolder, toFile)
	img, err := os.Create(path)
	if err != nil {
		return err
	}
	defer img.Close()

	resp, err := http.Get(fromURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(img, resp.Body)
	return err
}
