package main

import (
	"fmt"

	"scraper"
)

func main() {
	urls := scraper.ExtractProductLinks("amsterdam72.html")

	for idx, url := range urls {
		fmt.Println(idx+1, url)
		sku := scraper.SKUFromProductLink(url)

		fmt.Println(sku.ColorName)
	}
}
