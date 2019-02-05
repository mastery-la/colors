package main

import (
	"fmt"
	"log"

	paint "github.com/mastery-la/colors/colors"
)

func main() {
	// Step 1: Scrape dickblick.com for paint SKUs
	s, err := paint.NewScraper(
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
	d, err := paint.NewDownloader(s, ".outputs")
	if err != nil {
		log.Fatal(err)
	}
	d.Download()

	// Step 3: Generate Palette from Paint SKUs and their swatches
	g := paint.NewGenerator(s, d)
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
