package goat

import (
	"bytes"
	"encoding/xml"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dsvg"
	"io"
)

// FIXME
// this probably needs to be fixed in draw2d
// when drawing text, XML is not escaped and generated
// SVG is invalid if contains < > &
// We wrap the GC only for the SVG case to escape the text here
type FixDraw2DEscapeXml struct {
	draw2d.GraphicContext
}

func (gc *FixDraw2DEscapeXml) FillStringAt(text string, x, y float64) (width float64) {
	buf := new(bytes.Buffer)
	xml.EscapeText(buf, []byte(text))
	return gc.GraphicContext.FillStringAt(buf.String(), x, y)
}

func AsciiToSvg(canvas Canvas, out io.Writer) error {
	initFontCache()
	svg := draw2dsvg.NewSvg()
	svg.FontMode = draw2dsvg.SysFontMode
	var gc draw2d.GraphicContext
	gc = draw2dsvg.NewGraphicContext(svg)
	gc = &FixDraw2DEscapeXml{gc}

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
