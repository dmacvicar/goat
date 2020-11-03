package goat

import (
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dkit"
	"image/color"
	"math"
)

func RenderAscii(canvas Canvas, gc draw2d.GraphicContext) {
	gc.Translate(8, 16)

	for _, l := range canvas.Lines() {
		l.Draw(gc)
	}

	for _, t := range canvas.Triangles() {
		t.Draw(gc)
	}

	for _, c := range canvas.RoundedCorners() {
		c.Draw(gc)
	}

	for _, c := range canvas.Circles() {
		c.Draw(gc)
	}

	for _, b := range canvas.Bridges() {
		b.Draw(gc)
	}

	for _, t := range canvas.Text() {
		t.Draw(gc)
	}
}

// Draw a straight line as an SVG path.
func (l Line) Draw(gc draw2d.GraphicContext) {

	start := l.start.asPixel()
	stop := l.stop.asPixel()

	// For cases when a vertical line hits a perpendicular like this:
	//
	//   |          |
	//   |    or    v
	//  ---        ---
	//
	// We need to nudge the vertical line half a vertical cell in the
	// appropriate direction in order to meet up cleanly with the midline of
	// the cell next to it.

	// A diagonal segment all by itself needs to be shifted slightly to line
	// up with _ baselines:
	//     _
	//      \_
	//
	// TODO make this a method on Line to return accurate pixel
	if l.lonely {
		switch l.orientation {
		case NE:
			start.x -= 4
			stop.x -= 4
			start.y += 8
			stop.y += 8
		case SE:
			start.x -= 4
			stop.x -= 4
			start.y -= 8
			stop.y -= 8
		case S:
			start.y -= 8
			stop.y -= 8
		}

		// Half steps
		switch l.chop {
		case N:
			stop.y -= 8
		case S:
			start.y += 8
		}
	}

	if l.needsNudgingDown {
		stop.y += 8
		if l.horizontal() {
			start.y += 8
		}
	}

	if l.needsNudgingLeft {
		start.x -= 8
	}

	if l.needsNudgingRight {
		stop.x += 8
	}

	if l.needsTinyNudgingLeft {
		start.x -= 4
		if l.orientation == NE {
			start.y += 8
		} else if l.orientation == SE {
			start.y -= 8
		}
	}

	if l.needsTinyNudgingRight {
		stop.x += 4
		if l.orientation == NE {
			stop.y -= 8
		} else if l.orientation == SE {
			stop.y += 8
		}
	}

	gc.MoveTo(float64(start.x), float64(start.y))
	gc.LineTo(float64(stop.x), float64(stop.y))
	gc.Stroke()
}

// Draw a solid triable as an SVG polygon element.
func (t Triangle) Draw(gc draw2d.GraphicContext) {
	// https://www.w3.org/TR/SVG/shapes.html#PolygonElement

	/*
		   	+-----+-----+
		    |    /|\    |
		    |   / | \   |
		  x +- / -+- \ -+
			| /   |   \ |
			|/    |    \|
		    +-----+-----+
		          y
	*/

	x, y := float64(t.start.asPixel().x), float64(t.start.asPixel().y)
	r := 0.0

	x0 := x + 8
	y0 := y
	x1 := x - 4
	y1 := y - 0.35*16
	x2 := x - 4
	y2 := y + 0.35*16

	switch t.orientation {
	case N:
		r = 270
		if t.needsNudging {
			x0 += 8
			x1 += 8
			x2 += 8
		}
	case NE:
		r = 300
		x0 += 4
		x1 += 4
		x2 += 4
		if t.needsNudging {
			x0 += 6
			x1 += 6
			x2 += 6
		}
	case NW:
		r = 240
		x0 += 4
		x1 += 4
		x2 += 4
		if t.needsNudging {
			x0 += 6
			x1 += 6
			x2 += 6
		}
	case W:
		r = 180
		if t.needsNudging {
			x0 -= 8
			x1 -= 8
			x2 -= 8
		}
	case E:
		r = 0
		if t.needsNudging {
			x0 -= 8
			x1 -= 8
			x2 -= 8
		}
	case S:
		r = 90
		if t.needsNudging {
			x0 += 8
			x1 += 8
			x2 += 8
		}
	case SW:
		r = 120
		x0 += 4
		x1 += 4
		x2 += 4
		if t.needsNudging {
			x0 += 6
			x1 += 6
			x2 += 6
		}
	case SE:
		r = 60
		x0 += 4
		x1 += 4
		x2 += 4
		if t.needsNudging {
			x0 += 6
			x1 += 6
			x2 += 6
		}
	}

	rad := r * math.Pi / 180
	// rotate the triangle around the center point
	// Translating to x,y, rotating and translating back
	x0r := float64((x0-x)*math.Cos(rad)-(y0-y)*math.Sin(rad)) + x
	y0r := float64((x0-x)*math.Sin(rad)+(y0-y)*math.Cos(rad)) + y

	x1r := float64((x1-x)*math.Cos(rad)-(y1-y)*math.Sin(rad)) + x
	y1r := float64((x1-x)*math.Sin(rad)+(y1-y)*math.Cos(rad)) + y

	x2r := float64((x2-x)*math.Cos(rad)-(y2-y)*math.Sin(rad)) + x
	y2r := float64((x2-x)*math.Sin(rad)+(y2-y)*math.Cos(rad)) + y

	gc.MoveTo(x0r, y0r)
	gc.LineTo(x1r, y1r)
	gc.LineTo(x2r, y2r)
	gc.Close()
	gc.SetFillColor(color.Black)
	gc.Fill()
	gc.Stroke()
}

// Draw a solid circle as an SVG circle element.
func (c *Circle) Draw(gc draw2d.GraphicContext) {
	fill := color.White

	if c.bold {
		fill = color.Black
	}

	pixel := c.start.asPixel()

	gc.MoveTo(float64(pixel.x), float64(pixel.y))
	draw2dkit.Circle(gc, float64(pixel.x), float64(pixel.y), float64(6))
	gc.SetStrokeColor(color.Black)
	gc.SetFillColor(fill)
	// FIXME for some reason using gc.FillStroke() generates a line inside the
	// circle
	gc.Fill()

	draw2dkit.Circle(gc, float64(pixel.x), float64(pixel.y), float64(6))
	gc.SetStrokeColor(color.Black)
	gc.Stroke()
}

// Draw a single text character as an SVG text element.
func (t Text) Draw(gc draw2d.GraphicContext) {
	p := t.start.asPixel()
	c := t.contents

	opacity := 0

	// Markdeep special-cases these character and treats them like a
	// checkerboard.
	switch c {
	case "▉":
		opacity = -64
	case "▓":
		opacity = 64
	case "▒":
		opacity = 128
	case "░":
		opacity = 191
	}

	if opacity != 0 {
		draw2dkit.Rectangle(gc, float64(p.x-4), float64(p.y-8), float64(p.x-4+8), float64(p.y-8+16))
		gc.FillStroke()
		return
	}

	gc.SetFontData(draw2d.FontData{Name: "Liberation Mono,monospace", Family: draw2d.FontFamilyMono})
	gc.SetFillColor(color.Black)
	gc.SetFontSize(13)
	gc.FillStringAt(c, float64(p.x), float64(p.y+4))
}

func dotProduct(ux, uy, vx, vy float64) float64 {
	return ux*vx + uy*vy
}

func magnitude(ux, uy float64) float64 {
	return math.Sqrt(math.Pow(ux, 2) + math.Pow(uy, 2))
}

// https://www.w3.org/TR/SVG11/implnote.html#ArcImplementationNotes
// (F.6.5.4)
func angleBetween(ux, uy, vx, vy float64) float64 {
	angle := math.Acos(dotProduct(ux, uy, vx, vy) / (magnitude(ux, uy) * magnitude(vx, vy)))
	if (ux*vy - uy*vx) < 0 {
		return -angle
	}
	return angle
}

// Draws an arc using the SVG conventions of starting and end points as specified
// in:
// https://www.w3.org/TR/SVG11/implnote.html#ArcImplementationNotes
// For this particular use-case, we assume largeArcFlag and rotationAngle = 0
// which simplifies the equations
func svgArcTo(gc draw2d.GraphicContext, startX, startY, endX, endY, rx, ry float64, sweepFlag int) {
	// Xp means X'
	// https://www.w3.org/TR/SVG11/implnote.html#ArcImplementationNotes
	// (F.6.5.1) simplified because of rotation angle always 0
	startXp := float64(startX-endX) / 2.0
	startYp := float64(startY-endY) / 2.0

	// (F.6.5.2)
	factorTerm := math.Sqrt(
		(math.Pow(rx, 2)*math.Pow(ry, 2) - math.Pow(rx, 2)*math.Pow(startYp, 2) - math.Pow(ry, 2)*math.Pow(startXp, 2)) / (math.Pow(rx, 2)*math.Pow(startYp, 2) + math.Pow(ry, 2)*math.Pow(startXp, 2)))
	// largeArc flag is always 0 in this case
	if sweepFlag == 0 {
		factorTerm = -factorTerm
	}

	// (F.6.5.2)
	cxp := factorTerm * (rx * startYp) / ry
	cyp := -factorTerm * (ry * startXp) / rx

	// (F.6.5.3)
	// rotation angle is again zero, makes equation simpler
	cx := cxp + float64(startX+endX)/2.0
	cy := cyp + float64(startY+endY)/2.0

	// (F.6.5.5)
	theta1 := angleBetween(1, 0, (startXp-cxp)/rx, (startYp-cyp)/ry)
	// (F.6.5.6)
	deltaTheta := math.Mod(angleBetween((startXp-cxp)/rx, (startYp-cyp)/ry, (-startXp-cxp)/rx, (-startYp-cyp)/ry), 2*math.Pi)
	if sweepFlag == 0 && (deltaTheta > 0) {
		deltaTheta = deltaTheta - 2*math.Pi
	} else if sweepFlag == 1 && (deltaTheta < 0) {
		deltaTheta = deltaTheta + 2*math.Pi
	}

	gc.Save()
	gc.MoveTo(float64(startX), float64(startY))
	gc.ArcTo(cx, cy, rx, ry, theta1, deltaTheta)
	gc.Stroke()
	gc.Restore()
}

// Draw a rounded corner as an SVG elliptical arc element.
func (c *RoundedCorner) Draw(gc draw2d.GraphicContext) {
	// https://www.w3.org/TR/SVG/paths.html#PathDataEllipticalArcCommands

	x, y := c.start.asPixelXY()
	startX, startY, endX, endY, sweepFlag := 0, 0, 0, 0, 0

	switch c.orientation {
	case NW:
		startX = x + 8
		startY = y
		endX = x - 8
		endY = y + 16
	case NE:
		sweepFlag = 1
		startX = x - 8
		startY = y
		endX = x + 8
		endY = y + 16
	case SE:
		sweepFlag = 1
		startX = x + 8
		startY = y - 16
		endX = x - 8
		endY = y
	case SW:
		startX = x - 8
		startY = y - 16
		endX = x + 8
		endY = y
	}

	svgArcTo(gc, float64(startX), float64(startY), float64(endX), float64(endY), float64(16), float64(16), sweepFlag)
}

// Draw a bridge as an SVG elliptical arc element.
func (b Bridge) Draw(gc draw2d.GraphicContext) {
	x, y := b.start.asPixelXY()
	sweepFlag := 1

	if b.orientation == W {
		sweepFlag = 0
	}

	svgArcTo(gc, float64(x), float64(y-8), float64(x), float64(y+8), float64(9), float64(9), sweepFlag)
}
