package goat

import (
	"io"
	"encoding/xml"
	 "github.com/llgcode/draw2d/draw2dsvg"
)

func ASCIItoSVG(in io.Reader, out io.Writer) {
	svg := draw2dsvg.NewSvg()
	svg.FontMode = draw2dsvg.SysFontMode
	gc := draw2dsvg.NewGraphicContext(svg)
	RenderASCII(in, gc)
	out.Write([]byte(xml.Header))
	encoder := xml.NewEncoder(out)
	encoder.Indent("", "\t")
	_ = encoder.Encode(svg)
}
