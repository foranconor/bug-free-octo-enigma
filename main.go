package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	clip "github.com/ctessum/go.clipper"
	"github.com/kr/pretty"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

var (
	red    = color.RGBA{251, 128, 114, 255}
	green  = color.RGBA{141, 211, 199, 255}
	blue   = color.RGBA{190, 186, 218, 255}
	yellow = color.RGBA{255, 255, 179, 255}
)

type Point struct {
	X float64
	Y float64
}

type Line struct {
	Start Point
	End   Point
}

type Section struct {
	Kind       string
	Steps      int
	StartWidth float64
	EndWidth   float64
}

type Stringer struct {
	ControlPoints []Line
	BottomEnd     []Line
	TopEnd        []Line
	Section       Section
}

type Config map[string]float64

// Plan
// 1. Construct list of tread nosing points
// 2. filter to those that incur a change of direction
// 3. collect:
//	1, start and end lines
//	2, stringer top and bottom lines
//	3, join lines
// 4. make stringer objects
//
//
// TODO things still needing answering
// * corner point, how to split?

func main() {

	cfg := make(map[string]float64)

	cfg["nosing"] = 25
	cfg["riser_thickness"] = 15
	cfg["tread_thickness"] = 21
	cfg["wedge_angle"] = 0.125

	c := canvas.New(500, 300)
	ctx := canvas.NewContext(c)
	ctx.SetStrokeWidth(2)

	var pth clip.Path

	ps := make([]Point, 0)
	ps = append(ps, Point{0, 0})
	ps = append(ps, Point{100, 0})
	ps = append(ps, Point{100, 100})
	ps = append(ps, Point{0, 100})

	for _, p := range ps {
		pth = append(pth, clip.NewIntPointFromFloat(p.X, p.Y))
	}
	pretty.Println(pth)

	off := clip.NewClipperOffset()
	off.AddPath(pth, clip.JtRound, clip.EtClosedLine)

	qs := off.Execute(10)

	pretty.Println(qs)

	ts, _ := tr(Point{0, 0}, cfg)
	drawPoints(ts, ctx, yellow)

	c.Fit(20)
	err := renderers.Write("testing.png", c, canvas.DPMM(3.2))
	if err != nil {
		log.Fatal(err)
	}
}

func drawPoints(ps []Point, ctx *canvas.Context, c color.Color) {
	if len(ps) < 2 {
		return
	}
	ctx.SetStrokeColor(c)
	ctx.MoveTo(ps[0].X, ps[0].Y)
	for i := 1; i < len(ps); i++ {
		ctx.LineTo(ps[i].X, ps[i].Y)
	}
	ctx.Stroke()
}

func tr(tip Point, cfg Config) ([]Point, []Point) {
	const (
		TC = 300
		RC = 300
	)
	ts := make([]Point, 0)
	ts = append(ts, tip)
	ts = append(ts, Point{
		tip.X + TC + cfg["nosing"] + cfg["riser_thickness"],
		tip.Y,
	})
	ts = append(ts, Point{
		tip.X + TC + cfg["nosing"] + cfg["riser_thickness"],
		tip.Y - math.Tan(cfg["wedge_angle"])*TC - cfg["tread_thickness"],
	})
	ts = append(ts, Point{
		tip.X + cfg["nosing"] + cfg["riser_thickness"],
		tip.Y - cfg["tread_thickness"],
	})
	ts = append(ts, Point{
		tip.X,
		tip.Y - cfg["tread_thickness"],
	})
	ts = append(ts, tip)
	return ts, nil
}

func drawPath(ps clip.Path, ctx *canvas.Context, c color.Color) {
	if len(ps) < 2 {
		return
	}
	ctx.SetStrokeColor(c)
	ctx.MoveTo(float64(ps[0].X), float64(ps[0].Y))
	for i := 1; i < len(ps); i++ {
		ctx.LineTo(float64(ps[i].X), float64(ps[i].Y))
	}
	ctx.Stroke()
}

func (l *Line) Draw(ctx *canvas.Context, c color.Color) {
	ctx.SetStrokeColor(c)
	ctx.MoveTo(l.Start.X, l.Start.Y)
	ctx.LineTo(l.End.X, l.End.Y)
	ctx.Stroke()
}

func (p *Point) Draw(ctx *canvas.Context, c color.Color) {
	const SIZE = 20
	a := Line{Point{p.X - SIZE, p.Y}, Point{p.X + SIZE, p.Y}}
	b := Line{Point{p.X, p.Y - SIZE}, Point{p.X, p.Y + SIZE}}
	a.Draw(ctx, c)
	b.Draw(ctx, c)
}
func intersection(a, b Line) (Point, error) {
	i := a.Start.X - a.End.X
	j := b.Start.Y - b.End.Y
	k := a.Start.Y - a.End.Y
	l := b.Start.X - b.End.X
	denom := i*j - k*l
	if denom == 0 {
		return Point{0, 0}, fmt.Errorf("lines are either parallel or coincident")
	}
	m := a.Start.X*a.End.Y - a.Start.Y*a.End.X
	n := b.Start.X*b.End.Y - b.Start.Y*b.End.X
	xn := m*l - i*n
	yn := m*j - k*n
	return Point{xn / denom, yn / denom}, nil
}

func (l *Line) Offset(d float64) Line {
	scale := l.Length() / d

	sl := l.Scale(1 / scale)
	p := sl.Rotate(math.Pi / 2)
	q := p.Translate(l.ZeroStart())
	return Line{
		p.End,
		q.End,
	}
}

func (l *Line) Reverse() Line {
	return Line{
		l.End,
		l.Start,
	}
}

func (l *Line) Length() float64 {
	a := l.Start.X - l.End.X
	b := l.Start.Y - l.End.Y
	return math.Sqrt(a*a + b*b)
}

func (p *Point) Rotate(t float64) Point {
	x := p.X*math.Cos(t) - p.Y*math.Sin(t)
	y := p.X*math.Sin(t) + p.Y*math.Cos(t)
	return Point{x, y}
}

func (p *Point) Scale(s float64) Point {
	return Point{p.X * s, p.Y * s}
}

func (l *Line) Scale(s float64) Line {
	e := l.ZeroStart()
	f := e.Scale(s)
	return Line{
		l.Start,
		f.Translate(l.Start),
	}
}

func (l *Line) ZeroStart() Point {
	return l.End.Translate(Point{-1 * l.Start.X, -1 * l.Start.Y})
}

func (l *Line) Rotate(t float64) Line {
	e := l.ZeroStart()
	f := e.Rotate(t)
	return Line{
		l.Start,
		f.Translate(l.Start),
	}
}

func (p *Point) Translate(d Point) Point {
	return Point{p.X + d.X, p.Y + d.Y}
}

func (l *Line) Translate(d Point) Line {
	return Line{
		l.Start.Translate(d),
		l.End.Translate(d),
	}
}
