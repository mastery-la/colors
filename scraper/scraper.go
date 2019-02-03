package scraper

import (
	"sort"

	"github.com/schollz/pluck/pluck"
)

// PaintSKU holds information about a unique SKU of paint
// PaintSKUs include the form-factor that the paint is being sold.
// The same color may have multiple SKUs for different sizes of paint tubes.
type PaintSKU struct {
	ColorName      string
	ProductURL     string
	ColorSwatchURL string
}

// ExtractPaintSKUs takes a path to an html file from dickblick.com and returns
// an array of PaintSKUs that were extracted by crawling the product links.
func ExtractPaintSKUs(fromFile string) []PaintSKU {
	urls := extractProductURLs(fromFile)
	skus := skusFromProductLinks(urls)
	skus = removeDuplicates(skus)
	skus = sortAlpha(skus)

	return skus
}

func extractProductURLs(fromFile string) []string {
	p, _ := pluck.New()

	p.Add(pluck.Config{
		Activators:  []string{"<", `class="itemskulink"`, "href", `"`},
		Deactivator: `"`,
		Limit:       -1,
	})

	p.PluckFile(fromFile)

	result := p.Result()
	paths := result["0"].([]string)

	baseURL := "https://www.dickblick.com"

	urls := mapOverSlice(paths, func(path string) string {
		return baseURL + path
	})

	return urls

}

func skusFromProductLinks(urls []string) []PaintSKU {
	c := make(chan PaintSKU, 100)
	var skus []PaintSKU

	for _, url := range urls {
		go fetchSKU(c, url)
	}

	for range urls {
		skus = append(skus, <-c)
	}

	return skus
}

func fetchSKU(skus chan<- PaintSKU, url string) {
	p, _ := pluck.New()

	// Extract Color Name
	p.Add(pluck.Config{
		Activators:  []string{"<", "h2", `class="skutitle"`, "—"},
		Deactivator: "</h2>",
		Limit:       -1,
	})

	// Extract Color Swatch URL
	p.Add(pluck.Config{
		Activators:  []string{"<", "div", `id="mainphotowrapper"`, "<", "img", "src", `"`},
		Deactivator: `"`,
		Limit:       -1,
	})

	// Extract Color Number
	p.Add(pluck.Config{
		Activators:  []string{"<", "td", `class="skuelement skuNo."`, ">"},
		Deactivator: `</td>`,
		Limit:       -1,
	})

	p.PluckURL(url)
	result := p.Result()

	skus <- PaintSKU{
		ColorName:      result["0"].(string),
		ProductURL:     url,
		ColorSwatchURL: result["1"].(string),
	}

}

func mapOverSlice(strings []string, f func(string) string) []string {
	mapped := make([]string, len(strings))
	for i, v := range strings {
		mapped[i] = f(v)
	}
	return mapped
}

func removeDuplicates(skus []PaintSKU) []PaintSKU {
	visitors := make(map[string]struct{}, len(skus))
	count := 0
	for _, v := range skus {
		if _, visited := visitors[v.ColorName]; visited {
			continue
		}
		visitors[v.ColorName] = struct{}{}
		skus[count] = v
		count++
	}
	return skus[:count]
}

func sortAlpha(skus []PaintSKU) []PaintSKU {
	sorted := make([]PaintSKU, len(skus))
	copy(sorted, skus)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ColorName < sorted[j].ColorName
	})

	return sorted
}
