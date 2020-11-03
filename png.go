package goat

import (
	"bufio"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/png"
	"io"
)


func AsciiToPng(canvas Canvas, out io.Writer) error {
	initFontCache()
	dest := image.NewRGBA(image.Rect(0, 0, (canvas.Width+1)*8*PngScale, canvas.Height*16*PngScale+8+1))
	gc := draw2dimg.NewGraphicContext(dest)
	RenderAscii(canvas, gc)
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
