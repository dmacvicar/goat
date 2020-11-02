package goat

import (
	"encoding/xml"
	"github.com/llgcode/draw2d/draw2dsvg"
	"io"
)

func AsciiToSvg(canvas Canvas, out io.Writer) error {
	svg := draw2dsvg.NewSvg()
	svg.FontMode = draw2dsvg.SysFontMode
	gc := draw2dsvg.NewGraphicContext(svg)
	RenderAscii(canvas, gc)
	_, err := out.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	encoder := xml.NewEncoder(out)
	encoder.Indent("", "\t")
	err = encoder.Encode(svg)
	if err != nil {
		return err
	}
	return nil
}
