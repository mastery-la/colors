package scraper

import "github.com/schollz/pluck/pluck"

// PaintSKU holds information about a unique SKU of paint
// PaintSKUs include the form-factor that the paint is being sold.
// The same color may have multiple SKUs for different sizes of paint tubes.
type PaintSKU struct {
	ItemNumber     string
	ProductURL     string
	ColorNumber    string
	ColorName      string
	ColorSwatchURL string
}

// SKUFromProductLink returns a PaintSKU that it extracts from the product
// at the specified link. Assumes that url is a dickblick.com paint product.
func SKUFromProductLink(url string) PaintSKU {
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

	p.PluckURL(url)
	result := p.Result()

	return PaintSKU{
		ItemNumber:     "12312",
		ProductURL:     url,
		ColorNumber:    "string",
		ColorName:      result["0"].(string),
		ColorSwatchURL: result["1"].(string),
	}
}

// ExtractProductLinks Given a path to an html file extracted from dickblick.com,
// returns an array of URLs to the products.
func ExtractProductLinks(fromFile string) []string {
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

	urls := Map(paths, func(path string) string {
		return baseURL + path
	})

	return urls
}

// Map over a slice of strings, modifying each string with the provided function.
func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
