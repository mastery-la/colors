package downloader

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mastery-la/colors/scraper"
)

// Downloader holds neccessary information to use throughout downloading
type Downloader struct {
	Scraper         *scraper.Scraper
	OutputFolder    string
	Results         []string
	downloadChannel chan string
}

// New creates a new instance of Downloader
func New(scraper *scraper.Scraper, outputFolder string) (*Downloader, error) {
	d := new(Downloader)

	newpath := filepath.Join(".", outputFolder)
	err := os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	d.OutputFolder = newpath
	d.Scraper = scraper
	d.downloadChannel = make(chan string, 100)

	return d, nil
}

func (d *Downloader) Download() {
	for _, sku := range d.Scraper.Results {
		go d.SaveImage(sku)
	}

	for range d.Scraper.Results {
		d.Results = append(d.Results, <-d.downloadChannel)
	}
}

// SaveImage downloads the image from the given URL and saves to file
func (d *Downloader) SaveImage(fromSKU scraper.PaintSKU) {
	path := filepath.Join(d.OutputFolder, fromSKU.Slug+".jpg")
	img, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer img.Close()

	resp, err := http.Get(fromSKU.ColorSwatchURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(img, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	d.downloadChannel <- path
}
