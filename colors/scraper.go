package colors

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/schollz/pluck/pluck"
)

type Scraper struct {
	Name          string
	PaintType     string
	Brand         string
	InputFile     string
	ExtractedURLs []string
	Results       []PaintSKU
}

func NewScraper(name string, paintType string, brand string, inputFile string) (*Scraper, error) {
	s := new(Scraper)

	newpath := filepath.Join(".", inputFile)
	_, err := os.Stat(newpath)
	if err != nil {
		return nil, err
	}

	s.Name = name
	s.PaintType = paintType
	s.Brand = brand
	s.InputFile = inputFile

	return s, nil
}

// PaintSKU holds information about a unique SKU of paint
// PaintSKUs include the form-factor that the paint is being sold.
// The same color may have multiple SKUs for different sizes of paint tubes.
type PaintSKU struct {
	ColorName      string
	Slug           string
	ProductURL     string
	ColorSwatchURL string
}

// Extract takes a path to an html file from dickblick.com and returns
// an array of PaintSKUs that were extracted by crawling the product links.
func (s *Scraper) Scrape() {
	s.extractProductURLs()
	s.fetchSKUs()
	s.sortAndDedupSKUs()
}

func (s *Scraper) extractProductURLs() {
	p, _ := pluck.New()

	p.Add(pluck.Config{
		Activators:  []string{"<", `class="itemskulink"`, "href", `"`},
		Deactivator: `"`,
		Limit:       -1,
	})

	p.PluckFile(s.InputFile)

	result := p.Result()
	paths := result["0"].([]string)

	baseURL := "https://www.dickblick.com"

	urls := mapOverSlice(paths, func(path string) string {
		return baseURL + path
	})

	s.ExtractedURLs = urls
}

func (s *Scraper) fetchSKUs() {
	c := make(chan PaintSKU, 100)
	var skus []PaintSKU

	for _, url := range s.ExtractedURLs {
		go fetchSKU(c, url)
	}

	for range s.ExtractedURLs {
		skus = append(skus, <-c)
	}

	s.Results = skus
}

func (s *Scraper) sortAndDedupSKUs() {
	skus := s.Results

	skus = sortAlpha(skus)
	skus = removeDuplicates(skus)

	s.Results = skus
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

	name := result["0"].(string)
	swatch := result["1"].(string)
	slug := strings.ToLower(strings.Replace(name, " ", "-", -1))

	skus <- PaintSKU{
		ColorName:      name,
		Slug:           slug,
		ProductURL:     url,
		ColorSwatchURL: swatch,
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
