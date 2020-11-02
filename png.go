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

func AsciiToPng(canvas Canvas, out io.Writer) error {
	switch runtime.GOOS {
	case "linux":
		draw2d.SetFontFolder("/usr/share/fonts/truetype")
	default:
		log.Printf("Warning: I don't know how to looks for fonts on %s yet", runtime.GOOS)
	}

		// 		canvas.Height*16+8+1, (canvas.Width+1)*8,

	dest := image.NewRGBA(image.Rect(0, 0, (canvas.Width+1)*8, canvas.Height*16+8+1))
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
