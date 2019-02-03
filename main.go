package main

import (
	"fmt"

	"github.com/mastery-la/colors/scraper"
)

func main() {
	urls := scraper.ExtractProductLinks("inputs/amsterdam72.html")

	for idx, url := range urls {
		fmt.Println(idx+1, url)
		sku := scraper.SKUFromProductLink(url)

		fmt.Println(sku.ColorName)
	}
}
