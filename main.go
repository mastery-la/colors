package main

import (
	"fmt"
	"strings"

	"github.com/mastery-la/colors/downloader"
	"github.com/mastery-la/colors/scraper"
)

func main() {
	skus := scraper.ExtractPaintSKUs("inputs/amsterdam72.html")
	_, err := downloader.New("outputs")
	if err != nil {
		return
	}

	for _, sku := range skus {
		fmt.Println(sku.ColorName)
		fmt.Println(strings.Replace(sku.ColorName, " ", "", -1) + ".jpg")
		fmt.Println(sku.ColorSwatchURL)
		fmt.Print("\n")
	}

}
