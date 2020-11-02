package goat

import (
	"log"
	"io"
	"bufio"
	"image"
	"image/png"
	"runtime"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d"
)

func AsciiToPng(in io.Reader, out io.Writer) error {
	switch runtime.GOOS {
	case "linux":
		draw2d.SetFontFolder("/usr/share/fonts/truetype")
	default:
		log.Printf("Warning: I don't know how to looks for fonts on %s yet", runtime.GOOS)
	}
	dest := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	gc := draw2dimg.NewGraphicContext(dest)
	RenderAscii(in, gc)
	b := bufio.NewWriter(out)
	err := png.Encode(b, dest)
	if err != nil {
		return err
	}
	err = b.Flush()
	if err != nil {
		return err
	}
	return nil
}
