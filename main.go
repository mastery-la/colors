package main

import (
	"fmt"
	"log"

	"github.com/mastery-la/colors/downloader"
	"github.com/mastery-la/colors/generator"
	"github.com/mastery-la/colors/scraper"
)

func main() {
	// Step 1: Scrape dickblick.com for paint SKUs
	s, err := scraper.New(
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
	d, err := downloader.New(s, ".outputs")
	if err != nil {
		log.Fatal(err)
	}
	d.Download()

	// Step 3: Generate Palette from Paint SKUs and their swatches
	g := generator.New(s, d)
	g.Generate()

	for _, paint := range g.Result.Paints {
		fmt.Println(paint.Name)
		fmt.Println(paint.SwatchURL)
		fmt.Println("-------------")
		for index, color := range paint.Colors {
			fmt.Printf("%d. (%s) %s\n", index+1, color.Hex, color.SwatchURL)
		}
		fmt.Print("\n\n\n")
	}
}
