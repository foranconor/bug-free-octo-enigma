package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

const (
	NOSING               = 25
	BOTTOM_HORN          = 50
	TOP_HORN             = 80
	SKIRTING             = 35
	TOP_TREAD            = 45
	TREAD_THICKNESS      = 25
	RISER_THICKNESS      = 12
	WEDGE_ANGLE          = 0.125
	OVERSHOOT            = 12
	RISER_REBATE         = 5
	C                    = 100 // constant used for trig
	NOSING_TOP_RADIUS    = 10
	NOSING_BOTTOM_RADIUS = 10
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

func main() {
	c := canvas.New(500, 300)
	ctx := canvas.NewContext(c)
	ctx.SetStrokeWidth(2)

	steps := 5
	going := 280.0
	rise := 182.0
	width := 280.0

	ps := contour(steps, going, rise, width)
	drawPoints(ps, ctx, yellow)
	rs := rebates(steps, going, rise, width)
	for _, r := range rs {
		drawPoints(r, ctx, red)
	}

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

func (l *Line) Draw(ctx *canvas.Context) {
	ctx.MoveTo(l.Start.X, l.Start.Y)
	ctx.LineTo(l.End.X, l.End.Y)
	ctx.Stroke()
}

func (p *Point) Draw(ctx *canvas.Context, c color.Color) {
	ctx.SetStrokeColor(c)

	const SIZE = 30
	a := Line{Point{p.X - SIZE, p.Y}, Point{p.X + SIZE, p.Y}}
	b := Line{Point{p.X, p.Y - SIZE}, Point{p.X, p.Y + SIZE}}
	a.Draw(ctx)
	b.Draw(ctx)
}

func rebates(n int, g, r, w float64) [][]Point {
	cs := make([][]Point, 0)
	for i := 0; i < n; i++ {
		cs = append(cs, treadRebate(i, g, r, w))
		cs = append(cs, riserRebate(i, g, r, w))
	}
	return cs
}

func riserRebate(n int, g, r, w float64) []Point {
	ps := make([]Point, 0)
	at := Point{BOTTOM_HORN + NOSING + RISER_THICKNESS + float64(n)*g, r - TREAD_THICKNESS + float64(n)*r}
	end := at
	ps = append(ps, at)
	at = Point{at.X, at.Y + RISER_REBATE}
	ps = append(ps, at)
	at = Point{at.X - RISER_THICKNESS, at.Y}
	ps = append(ps, at)

	riserFront := Line{at, Point{at.X, at.Y - r}}
	pl := pitchLine(g, r)
	extent := pl.Offset(SKIRTING - w - OVERSHOOT)
	at, _ = intersection(extent, riserFront)
	ps = append(ps, at)
	rebateBack := Line{
		end,
		Point{
			end.X + C*math.Tan(WEDGE_ANGLE),
			end.Y - C,
		},
	}
	at, _ = intersection(extent, rebateBack)
	ps = append(ps, at)
	ps = append(ps, end)

	return ps
}

func treadRebate(n int, g, r, w float64) []Point {
	ps := make([]Point, 0)
	at := Point{BOTTOM_HORN + NOSING + RISER_THICKNESS + float64(n)*g, r - TREAD_THICKNESS + float64(n)*r}
	end := at
	ps = append(ps, at)
	at = Point{at.X - NOSING - RISER_THICKNESS, at.Y}
	ps = append(ps, at)
	at = Point{at.X, at.Y + TREAD_THICKNESS}
	ps = append(ps, at)

	treadTop := Line{at, Point{at.X + g, at.Y}}

	pl := pitchLine(g, r)
	extent := pl.Offset(SKIRTING - w - OVERSHOOT)

	at, _ = intersection(extent, treadTop)

	ps = append(ps, at)
	rebateBottom := Line{
		end,
		Point{
			end.X + C,
			end.Y - C*math.Tan(WEDGE_ANGLE),
		},
	}

	at, _ = intersection(extent, rebateBottom)
	ps = append(ps, at)
	ps = append(ps, end)

	return ps
}

// starting point is at the back of the riser under the tread
func nosing(p Point) []Point {
	ps := make([]Point, 0)
	brc := Point{p.X - RISER_THICKNESS - NOSING + NOSING_BOTTOM_RADIUS, p.Y + NOSING_BOTTOM_RADIUS}
	trc := Point{p.X - RISER_THICKNESS - NOSING + NOSING_TOP_RADIUS, p.Y + TREAD_THICKNESS - NOSING_TOP_RADIUS}
	_ = brc
	_ = trc

	if NOSING_BOTTOM_RADIUS == 0 {

	}

	return nil
}

func contour(n int, g, r, w float64) []Point {
	var (
		totalRise        = float64(n) * r
		totalGoing       = float64(n-1)*g + BOTTOM_HORN + NOSING + RISER_THICKNESS
		q          Point = Point{0, 0}
	)
	// lines we care about
	pitch := pitchLine(g, r)
	stringerTop := pitch.Offset(SKIRTING)
	top := Line{
		Point{0, totalRise},
		Point{totalGoing, totalRise},
	}
	topHorn := top.Offset(TOP_HORN)
	joist := Line{
		Point{totalGoing, 0},
		Point{totalGoing, totalRise},
	}
	farHorn := joist.Offset(-TOP_HORN)
	stringerBottom := stringerTop.Offset(-w)
	ground := Line{
		Point{0, 0},
		Point{totalGoing, 0},
	}

	// the points that make the contour
	ps := make([]Point, 0)
	ps = append(ps, q)
	q, _ = intersection(Line{q, Point{0, 1}}, stringerTop)
	ps = append(ps, q)
	q, _ = intersection(stringerTop, topHorn)
	ps = append(ps, q)
	q, _ = intersection(farHorn, topHorn)
	ps = append(ps, q)
	q, _ = intersection(farHorn, top)
	ps = append(ps, q)
	ps = append(ps, Point{totalGoing, totalRise})
	q, _ = intersection(joist, stringerBottom)
	ps = append(ps, q)
	q, _ = intersection(stringerBottom, ground)
	ps = append(ps, q)
	ps = append(ps, Point{0, 0})

	return ps
}

func pitchLine(g, r float64) Line {
	p := Line{
		Point{0, 0},
		Point{g, r},
	}
	return p.Translate(Point{BOTTOM_HORN, r})
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
