package goat

import (
	"fmt"
	"github.com/go-fonts/liberation/liberationmonoregular"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"log"
)

type goatFontCache map[string]*truetype.Font

func (fc goatFontCache) Store(fd draw2d.FontData, font *truetype.Font) {
	fc[fd.Name] = font
}

func (fc goatFontCache) Load(fd draw2d.FontData) (*truetype.Font, error) {
	font, stored := fc[fd.Name]
	if !stored {
		return nil, fmt.Errorf("font %s is not stored in font cache", fd.Name)
	}
	return font, nil
}

func initFontCache() {
	fontCache := goatFontCache{}
	// add font to cache
	gofont, err := truetype.Parse(liberationmonoregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	fontCache.Store(draw2d.FontData{Name: "Liberation Mono,monospace"}, gofont)

	draw2d.SetFontCache(fontCache)
}
