package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mastery-la/colors/colors"
)

func main() {
	// Step 1: Scrape dickblick.com for paint SKUs
	s, err := colors.NewScraper(
		"Amsterdam Standard Series Acrylics",
		"acrylic",
		"Amsterdam",
		".inputs/amsterdam72.html",
	)
	if err != nil {
		log.Fatal(err)
	}
	s.Scrape()

	// Step 2: Download Color Swatches
	d, err := colors.NewDownloader(s, "static/images")
	if err != nil {
		log.Fatal(err)
	}
	d.Download()

	// Step 3: Generate Palette from Paint SKUs and their swatches
	g := colors.NewGenerator(s, d)
	g.Generate()

	json, err := json.Marshal(g.Result.Paints)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(json))
}
