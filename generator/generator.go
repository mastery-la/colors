package generator

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/go-playground/colors"
	"github.com/mastery-la/colors/downloader"
	"github.com/mastery-la/colors/scraper"
	"github.com/oliamb/cutter"
)

type Generator struct {
	Scraper    *scraper.Scraper
	Downloader *downloader.Downloader
	Result     Palette
}

type Palette struct {
	Name      string
	PaintType string
	Brand     string
	Paints    []Paint
}

type Paint struct {
	Name      string
	Slug      string
	SwatchURL string
	Colors    []Color
}

type Color struct {
	Hex       string
	SwatchURL string
}

func New(scraper *scraper.Scraper, downloader *downloader.Downloader) *Generator {
	g := new(Generator)

	g.Scraper = scraper
	g.Downloader = downloader

	return g
}

func (g *Generator) Generate() {
	p := new(Palette)

	p.Name = g.Scraper.Name
	p.PaintType = g.Scraper.PaintType
	p.Brand = g.Scraper.Brand

	paints := generatePaints(g.Scraper.Results, g.Downloader.OutputFolder)
	p.Paints = paints

	g.Result = *p
}

func generatePaints(skus []scraper.PaintSKU, directory string) []Paint {
	var hexString string
	var err error
	var paints []Paint

	for _, sku := range skus {
		var paint Paint
		primaryPath := filepath.Join(".", directory, sku.Slug+".jpg")

		paint.Name = sku.ColorName
		paint.Slug = sku.Slug
		paint.SwatchURL = primaryPath

		var colors []Color
		for i := range []int{0, 1, 2} {
			swatchPath := fmt.Sprintf("%s.%d.png", primaryPath, i+1)

			hexString, err = crop(primaryPath, swatchPath, image.Point{(i * 100) + 10, 0})
			if err != nil {
				log.Fatal(err)
			}

			color := Color{
				Hex:       hexString,
				SwatchURL: swatchPath,
			}

			colors = append(colors, color)
		}

		paint.Colors = colors
		paints = append(paints, paint)
	}

	return paints
}

func crop(inPath string, outPath string, tlPoint image.Point) (hexString string, err error) {
	fi, err := os.Open(inPath)
	if err != nil {
		return
	}
	defer fi.Close()
	img, err := jpeg.Decode(fi)
	if err != nil {
		return
	}

	cImg, err := cutter.Crop(img, cutter.Config{
		Height:  80,             // height in pixel or Y ratio(see Ratio Option below)
		Width:   80,             // width in pixel or X ratio
		Mode:    cutter.TopLeft, // Accepted Mode: TopLeft, Centered
		Anchor:  tlPoint,        // Position of the top left point
		Options: 0,              // Accepted Option: Ratio
	})
	if err != nil {
		return
	}

	fo, err := os.Create(outPath)
	if err != nil {
		return
	}
	defer fo.Close()

	err = png.Encode(fo, cImg)
	cnrgba, err := averageColorFromImage(outPath)
	rgb, err := colors.RGB(cnrgba.R, cnrgba.G, cnrgba.B)
	if err != nil {
		return
	}
	hexString = rgb.ToHEX().String()
	return
}

func averageColorFromImage(url string) (cnrgba color.NRGBA, err error) {
	fi, err := os.Open(url)
	if err != nil {
		return
	}
	defer fi.Close()
	im, err := png.Decode(fi)
	if err != nil {
		return
	}
	rgba := imageToRGBA(im)
	size := rgba.Bounds().Size()
	w, h := size.X, size.Y
	var r, g, b int
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := rgba.RGBAAt(x, y)
			r += int(c.R)
			g += int(c.G)
			b += int(c.B)
		}
	}
	r /= w * h
	g /= w * h
	b /= w * h
	cnrgba = color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
	return
}

func imageToRGBA(src image.Image) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Rect, src, image.ZP, draw.Src)
	return dst
}
