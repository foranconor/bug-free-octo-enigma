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

const (
	SCALE = 1000
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
	Treads  []Point
	Contour []Line
	Name    string
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

	font := canvas.NewFontFamily("fira")
	if err := font.LoadFontFile("./FiraCodeNerdFont-Regular.ttf", canvas.FontRegular); err != nil {
		log.Fatal(err)
	}

	face := font.Face(120, yellow, canvas.FontRegular)

	width := 1000.0
	stair := make([]Section, 0)
	stair = append(stair, Section{
		Kind:       "straight",
		Steps:      3,
		EndWidth:   width,
		StartWidth: width,
	})
	stair = append(stair, Section{
		Kind:       "right",
		Steps:      3,
		EndWidth:   width,
		StartWidth: width,
	})
	stair = append(stair, Section{
		Kind:       "straight",
		Steps:      10,
		EndWidth:   width,
		StartWidth: width,
	})

	cfg := make(map[string]float64)
	cfg["nosing"] = 25
	cfg["riser_thickness"] = 15
	cfg["tread_thickness"] = 21
	cfg["wedge_angle"] = 0.125
	cfg["riser_rebate"] = 5
	cfg["going"] = 255
	cfg["rise"] = 190
	cfg["bottom_horn"] = 50
	cfg["skirting"] = 30
	cfg["stringer_width"] = 235
	cfg["top_horn_height"] = 80
	cfg["top_horn_length"] = 80

	pretty.Println(cfg)

	c := canvas.New(500, 300)
	ctx := canvas.NewContext(c)
	ctx.SetStrokeWidth(2)

	o := Point{0, 0}
	o.Draw(ctx, blue, "(0, 0)", face)

	left, _ := Tips(stair, cfg)

	for i, p := range left {
		p.Draw(ctx, red, fmt.Sprintf("T%d", i+1), face)
		ts, rs := tr(p, cfg)
		drawPoints(ts, ctx, blue)
		drawPoints(rs, ctx, blue)

	}
	drawPoints(left, ctx, green)

	stringers := make([]Stringer, 0)
	_ = stringers

	cps := ControlPoints(left, stair)

	for i, s := range stair {
		switch s.Kind {
		case "straight":
			contour := StraightContour(stair, cps, i, cfg, ctx, face)
			//drawPoints(contour, ctx, yellow)
			_ = contour

		default:
			contour := OutsideWinderContour(stair, cps, i, cfg, ctx, face)
			drawPoints(contour, ctx, yellow)
		}
	}

	// var pth clip.Path
	//
	// ps := make([]Point, 0)
	// ps = append(ps, Point{0, 0})
	// ps = append(ps, Point{100, 0})
	// ps = append(ps, Point{100, 100})
	// ps = append(ps, Point{0, 100})
	//
	// for _, p := range ps {
	// 	pth = append(pth, clip.NewIntPointFromFloat(p.X, p.Y))
	// }

	//	off := clip.NewClipperOffset()
	//
	//	ts, rs := tr(Point{0, 0}, cfg)
	//	drawPoints(rs, ctx, blue)
	//	tps := toPath(ts)
	//	off.AddPath(tps, clip.JtSquare, clip.EtClosedPolygon)
	//	res := off.Execute(-10.0 * SCALE)
	//	pretty.Println(res)
	//	//drawPoints(ts, ctx, red)
	//	drawPoints(fromPath(res[0]), ctx, yellow)
	//	off.Clear()
	//	off.AddPath(res[0], clip.JtSquare, clip.EtClosedPolygon)
	//	res = off.Execute(-5 * SCALE)
	//	drawPoints(fromPath(res[0]), ctx, green)
	//
	//	cl := clip.NewClipper(clip.IoNone)
	//
	//	cl.AddPath(toPath(ts), clip.PtSubject, true)
	//	cl.AddPath(toPath(rs), clip.PtClip, true)
	//	res, suc := cl.Execute1(clip.CtUnion, clip.PftEvenOdd, clip.PftEvenOdd)
	//	pretty.Println(res, suc)
	//	thing := fromPath(res[0])
	//	drawPoints(thing, ctx, yellow)

	c.Fit(20)
	err := renderers.Write("testing.png", c, canvas.DPMM(3.2))
	if err != nil {
		log.Fatal(err)
	}
}

func MakeContour(lines []Line) []Point {
	ps := make([]Point, 0)
	for i := 0; i < len(lines); i++ {
		p, _ := intersection(lines[i], lines[(i+1)%len(lines)])
		ps = append(ps, p)
	}
	ps = append(ps, ps[0])
	return ps
}

func OutsideWinderContour(stair []Section, pitches [][]Point, index int, config Config, ctx *canvas.Context, face *canvas.FontFace) []Point {
	lines := make([]Line, 0)

	pp := make([]Line, 0)
	for i := 1; i < len(pitches[index]); i++ {
		pitch := Line{pitches[index][i-1], pitches[index][i]}
		pp = append(pp, pitch)
	}
	pp = append(pp, Line{
		pitches[index][len(pitches[index])-1],
		pitches[index+1][0],
	})
	tops := make([]Line, 0)
	bottoms := make([]Line, 0)
	for _, l := range pp {
		t := l.Offset(config["skirting"])
		b := l.Offset(config["skirting"] - config["stringer_width"])
		tops = append(tops, t)
		bottoms = append(bottoms, b)
	}
	if index == 0 {
		lines = append(lines, BottomLines()...)
	} else {
		pretty.Println("here")
		tops[0].Draw(ctx, red)
		bottoms[0].Draw(ctx, red)
		drawPoints(pitches[index-1], ctx, yellow)
		pp[0].Start.Draw(ctx, yellow, "pp0", face)
		lines = append(lines, JoinBelow(pitches[index-1], pp[0].Start, tops[0], bottoms[0], config, true))
	}
	lines = append(lines, tops...)

	if index == len(stair)-1 {
		lines = append(lines, TopLines(stair, config)...)
	} else {
		lines = append(lines, JoinAbove(pitches[index+1], pp[len(pp)-1].End, tops[len(tops)-1], bottoms[len(bottoms)-1], config))
	}
	lines = append(lines, bottoms...)
	pretty.Println(lines)
	for _, l := range lines {
		l.Draw(ctx, blue)
	}
	contour := MakeContour(lines)
	return contour
}

func StraightContour(stair []Section, pitches [][]Point, index int, config Config, ctx *canvas.Context, face *canvas.FontFace) []Point {
	lines := make([]Line, 0)
	p := Line{pitches[index][0], pitches[index][1]}
	top, bottom := p.StringTopBottom(config)
	if index == 0 {
		// lines that make up the bottom of the stair
		lines = append(lines, BottomLines()...)
	} else {
		p.Start.Draw(ctx, yellow, "pstart", face)
		top.Draw(ctx, red)
		bottom.Draw(ctx, red)
		drawPoints(pitches[index-1], ctx, yellow)
		lines = append(lines, JoinBelow(pitches[index-1], p.Start, top, bottom, config, false))
	}
	// top line of stringer
	lines = append(lines, top)
	if index == len(stair)-1 {
		// top of the stair
		lines = append(lines, TopLines(stair, config)...)
	} else {
		// join with the section above
		ap := Line{p.End, pitches[index+1][1]}
		atop, abot := ap.StringTopBottom(config)
		a, _ := intersection(top, atop)
		b, _ := intersection(bottom, abot)
		lines = append(lines, Line{a, b})
	}
	lines = append(lines, bottom)
	contour := MakeContour(lines)
	return contour
}

func (l *Line) StringTopBottom(config Config) (Line, Line) {
	top := l.Offset(config["skirting"])
	bottom := l.Offset(config["skirting"] - config["stringer_width"])
	return top, bottom
}

func BottomLines() []Line {
	ground := Line{
		Point{0, 0}, Point{100, 0},
	}
	front := Line{
		Point{0, 0}, Point{0, 100},
	}
	ls := make([]Line, 2)
	ls[0] = ground
	ls[1] = front
	return ls
}

func JoinBelow(pitch []Point, start Point, top, bottom Line, config Config, winder bool) Line {
	pretty.Println("pitch", pitch, "start", start, "top", top, "bottom", bottom)
	// join with the stringer below
	offset := -1
	if winder {
		offset = -2
	}
	bp := Line{pitch[len(pitch)+offset], start}
	bt, bb := bp.StringTopBottom(config)
	a, _ := intersection(top, bt)
	b, _ := intersection(bottom, bb)
	return Line{a, b}
}

func JoinAbove(pitch []Point, end Point, top, bottom Line, config Config) Line {
	ap := Line{end, pitch[len(pitch)-1]}
	atop, abot := ap.StringTopBottom(config)
	a, _ := intersection(top, atop)
	b, _ := intersection(bottom, abot)
	return Line{a, b}
}

func TopLines(stair []Section, config Config) []Line {
	tr := TotalRise(stair, config)
	tg := OutsideGoing(stair, config) + config["bottom_horn"] + config["riser_thickness"] + config["nosing"]
	topFloor := Line{Point{0, tr}, Point{tg, tr}}
	joist := Line{Point{tg, tr}, Point{tg, 0}}
	hornTop := topFloor.Offset(config["top_horn_height"])
	hornBack := joist.Offset(config["top_horn_length"])
	ls := make([]Line, 4)
	ls[0] = hornTop
	ls[1] = hornBack
	ls[2] = topFloor
	ls[3] = joist
	return ls
}

func TotalRise(stair []Section, config Config) float64 {
	rise := 0.0
	for _, s := range stair {
		rise = rise + float64(s.Steps)*config["rise"]
	}
	return rise + config["rise"]
}

func OutsideGoing(stair []Section, config Config) float64 {
	going := 0.0
	for _, s := range stair {
		if s.Kind == "straight" {
			going = going + float64(s.Steps)*config["going"]
		} else {
			going = going + s.StartWidth + s.EndWidth
		}
	}
	return going
}

func ControlPoints(ts []Point, stair []Section) [][]Point {
	cps := make([][]Point, 0)
	for i := 0; i < len(stair); i++ {
		end := 0
		for j := i; j >= 0; j-- {
			end = end + stair[j].Steps
		}
		start := end - stair[i].Steps
		ps := make([]Point, 0)
		if stair[i].Kind == "straight" {
			ps = append(ps, ts[start])
			ps = append(ps, ts[end])
		} else {
			ps = append(ps, ts[start:end]...)
		}
		cps = append(cps, ps)
	}
	return cps
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

func toPath(ps []Point) clip.Path {
	qs := make(clip.Path, len(ps))
	for i, p := range ps {
		qs[i] = clip.NewIntPointFromFloat(p.X*SCALE, p.Y*SCALE)
	}
	return qs
}

func fromPath(ps clip.Path) []Point {
	qs := make([]Point, len(ps))
	for i, p := range ps {
		dp := p.ToDoublePoint()
		qs[i] = Point{dp.X / SCALE, dp.Y / SCALE}
	}
	return qs
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
	rs := make([]Point, 0)
	rTip := Point{
		tip.X + cfg["nosing"],
		tip.Y - cfg["tread_thickness"] + cfg["riser_rebate"],
	}
	rs = append(rs, rTip)
	rs = append(rs, Point{
		rTip.X + cfg["riser_thickness"],
		rTip.Y,
	})
	rs = append(rs, Point{
		rTip.X + cfg["riser_thickness"] + RC*math.Tan(cfg["wedge_angle"]),
		rTip.Y - cfg["tread_thickness"] - RC,
	})
	rs = append(rs, Point{
		rTip.X,
		rTip.Y - cfg["tread_thickness"] - RC,
	})
	rs = append(rs, rTip)

	return ts, rs
}

func Tips(sections []Section, config Config) ([]Point, []Point) {
	ps := make([]Point, 0)
	at := Point{config["bottom_horn"], config["rise"]}
	ps = append(ps, at)

	// each section is responsible for amking the step that steps into the next section
	for _, s := range sections {
		if s.Kind == "straight" {
			for i := 1; i <= s.Steps; i++ {
				ps = append(ps, Point{
					at.X + config["going"]*float64(i),
					at.Y + config["rise"]*float64(i),
				})
			}
			at = ps[len(ps)-1]
		} else {
			angle := math.Pi / 2 / float64(s.Steps)
			if s.Steps == 3 {
				// step 1
				ps = append(ps, Point{
					at.X + math.Tan(angle)*s.StartWidth,
					at.Y + config["rise"],
				})
				ps = append(ps, Point{
					at.X + s.EndWidth + s.StartWidth - s.StartWidth*math.Tan(angle),
					at.Y + 2*config["rise"],
				})
				ps = append(ps, Point{
					at.X + s.EndWidth + s.StartWidth,
					at.Y + 3*config["rise"],
				})
				at = ps[len(ps)-1]
			}
		}
	}
	return ps, nil
}

func (l *Line) Draw(ctx *canvas.Context, c color.Color) {
	ctx.SetStrokeColor(c)
	ctx.MoveTo(l.Start.X, l.Start.Y)
	ctx.LineTo(l.End.X, l.End.Y)
	ctx.Stroke()
}

func (p *Point) Draw(ctx *canvas.Context, c color.Color, text string, face *canvas.FontFace) {
	const (
		SIZE   = 20
		OFFSET = 10
	)

	a := Line{Point{p.X - SIZE, p.Y}, Point{p.X + SIZE, p.Y}}
	b := Line{Point{p.X, p.Y - SIZE}, Point{p.X, p.Y + SIZE}}
	ctx.DrawText(p.X+OFFSET, p.Y-OFFSET, canvas.NewTextBox(face, text, 0, 0, canvas.Left, canvas.Top, 0, 0))
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
