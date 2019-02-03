package main

import (
	"github.com/mastery-la/colors/downloader"
	"github.com/mastery-la/colors/generator"
	"github.com/mastery-la/colors/scraper"
)

func main() {
	skus := scraper.ExtractPaintSKUs(".inputs/amsterdam72.html")
	d := downloader.New(".outputs")

	for _, sku := range skus {
		d.SaveImage(sku.ColorSwatchURL, sku.Slug+".jpg")
	}

	generator.GenerateColorPalette(skus, ".outputs")
}
