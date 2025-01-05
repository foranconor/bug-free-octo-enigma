package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"threeTest/config"
	"threeTest/geo"

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

type Section struct {
	Kind       string
	Steps      int
	StartWidth float64
	EndWidth   float64
}

type Stringer struct {
	Treads  []geo.Point
	Rebates []geo.Point
	Contour []geo.Point
	Name    string
}

func main() {

	font := canvas.NewFontFamily("fira")
	if err := font.LoadFontFile("./FiraCode-Regular.ttf", canvas.FontRegular); err != nil {
		log.Fatal(err)
	}

	face := font.Face(12, 0, yellow, canvas.FontRegular)

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

	cfg := config.LoadConfig("something")

	c := canvas.New(500, 300)
	ctx := canvas.NewContext(c)
	ctx.SetStrokeWidth(2)

	o := geo.Point{X: 0, Y: 0}
	o.Draw(ctx, blue, "(0, 0)", face)

	left, _ := Tips(stair, cfg)

	// for i, p := range left {
	// 	//p.Draw(ctx, red, fmt.Sprintf("T%d", i+1), face)
	// 	ts, rs := tr(p, cfg)
	// 	//drawPoints(ts, ctx, blue)
	// 	//drawPoints(rs, ctx, blue)
	//
	// }
	//drawPoints(left, ctx, green)

	stringers := make([]Stringer, 0)
	_ = stringers

	cps := ControlPoints(left, stair)

	num := 1
	part := 1

	rebates := make(clip.Paths, 0)
	for _, t := range left {
		ts, rs := tr(t, cfg)
		tp := toPath(ts)
		rp := toPath(rs)
		rebates = append(rebates, tp, rp)
	}

	un := clip.NewClipper(clip.IoNone)
	un.AddPaths(rebates, clip.PtClip, true)
	unionedRebates, suc := un.Execute1(clip.CtUnion, clip.PftNonZero, clip.PftNonZero)
	pretty.Println(suc)

	for i, s := range stair {
		switch s.Kind {
		case "straight":
			contour := StraightContour(stair, cps, i, cfg, ctx, face)
			drawPoints(contour, ctx, yellow)
			name := fmt.Sprintf("%dLP%d", num, part)
			pretty.Println(name)
			part = part + 1
			path := toPath(contour)
			off := clip.NewClipperOffset()
			off.AddPath(path, clip.JtSquare, clip.EtClosedPolygon)
			res := off.Execute((cfg["trenching_radius"] + cfg["machining_extra"]) * 1000)
			//k := fromPath(res[0])
			//drawPoints(k, ctx, green)
			cl := clip.NewClipper(clip.IoNone)
			cl.AddPath(res[0], clip.PtClip, true)
			cl.AddPaths(unionedRebates, clip.PtSubject, true)
			intRes, _ := cl.Execute1(clip.CtIntersection, clip.PftEvenOdd, clip.PftEvenOdd)
			for _, j := range intRes {

				rem := fromPath(j)
				drawPoints(rem, ctx, red)
			}
		default:
			contour := OutsideWinderContour(stair, cps, i, cfg, ctx, face)
			contour[0].Draw(ctx, green, "start", face)
			contour[1].Draw(ctx, green, "next", face)

			drawPoints(contour, ctx, yellow)
			// split winder on turn
			sx := cps[i][0].X + s.EndWidth
			split := geo.Line{
				Start: geo.Point{X: sx, Y: 0},
				End:   geo.Point{X: sx, Y: 4000},
			}
			split.Draw(ctx, red)
			lower := make([]geo.Point, 0)
			upper := make([]geo.Point, 0)
			onLower := true
			_ = lower
			// go around the stringer and find the split
			for i, c := range contour {
				line := geo.Line{
					Start: c,
					End:   contour[(i+1)%len(contour)],
				}
				if onLower {
					lower = append(lower, c)
					if sx > line.Start.X && sx < line.End.X {
						// this line straddles to split, get the intersection
						x, _ := geo.Intersection(split, line)
						lower = append(lower, x)
						lower = append(lower, geo.Point{
							X: x.X + cfg["stringer_thickness"],
							Y: x.Y,
						})
						upper = append(upper, x)
						onLower = false
					}
				} else {
					upper = append(upper, c)
					if sx < line.Start.X && sx > line.End.X {
						x, _ := geo.Intersection(split, line)
						upper = append(upper, x)
						upper = append(upper, upper[0])
						lower = append(lower, geo.Point{
							X: x.X + cfg["stringer_thickness"],
							Y: x.Y,
						})
						lower = append(lower, x)
						onLower = true
					}
				}
			}
			drawPoints(lower, ctx, green)
			drawPoints(upper, ctx, blue)
			lowerName := fmt.Sprintf("%dLP%d", num, part)
			num = num + 1
			part = 1
			upperName := fmt.Sprintf("%dLP%d", num, part)
			part = 2
			pretty.Println(lowerName)
			pretty.Println(upperName)

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

func MakeContour(lines []geo.Line) []geo.Point {
	ps := make([]geo.Point, 0)
	for i := 0; i < len(lines); i++ {
		p, _ := geo.Intersection(lines[i], lines[(i+1)%len(lines)])
		ps = append(ps, p)
	}
	ps = append(ps, ps[0])
	return ps
}

func OutsideWinderContour(stair []Section, pitches [][]geo.Point, index int, config config.Params, ctx *canvas.Context, face *canvas.FontFace) []geo.Point {
	lines := make([]geo.Line, 0)

	pp := make([]geo.Line, 0)
	for i := 1; i < len(pitches[index]); i++ {
		pitch := geo.Line{
			Start: pitches[index][i-1],
			End:   pitches[index][i],
		}
		pp = append(pp, pitch)
	}
	pp = append(pp, geo.Line{
		Start: pitches[index][len(pitches[index])-1],
		End:   pitches[index+1][0],
	})
	tops := make([]geo.Line, 0)
	bottoms := make([]geo.Line, 0)
	for _, l := range pp {
		t := l.Offset(config["skirting"])
		b := l.Offset(config["skirting"] - config["stringer_width"])
		tops = append(tops, t)
		bottoms = append(bottoms, b)
	}
	if index == 0 {
		lines = append(lines, BottomLines()...)
	} else {
		lines = append(lines, JoinBelow(pitches[index-1], pp[0].Start, tops[0], bottoms[0], config, true))
	}
	lines = append(lines, tops...)

	if index == len(stair)-1 {
		lines = append(lines, TopLines(stair, config)...)
	} else {
		lines = append(lines, JoinAbove(pitches[index+1], pp[len(pp)-1].End, tops[len(tops)-1], bottoms[len(bottoms)-1], config))
	}
	reversedBottoms := make([]geo.Line, 0)
	for i := len(bottoms) - 1; i >= 0; i-- {
		rl := geo.Line{
			Start: bottoms[i].End,
			End:   bottoms[i].Start,
		}
		reversedBottoms = append(reversedBottoms, rl)
	}
	lines = append(lines, reversedBottoms...)

	contour := MakeContour(lines)
	return contour
}

func StraightContour(stair []Section, pitches [][]geo.Point, index int, config config.Params, ctx *canvas.Context, face *canvas.FontFace) []geo.Point {
	lines := make([]geo.Line, 0)
	p := geo.Line{
		Start: pitches[index][0],
		End:   pitches[index][1],
	}
	top, bottom := StringTopBottom(p, config)
	if index == 0 {
		// lines that make up the bottom of the stair
		lines = append(lines, BottomLines()...)
	} else {
		lines = append(lines, JoinBelow(pitches[index-1], p.Start, top, bottom, config, false))
	}
	// top line of stringer
	lines = append(lines, top)
	if index == len(stair)-1 {
		// top of the stair
		lines = append(lines, TopLines(stair, config)...)
	} else {
		// join with the section above
		ap := geo.Line{
			Start: p.End,
			End:   pitches[index+1][1],
		}
		atop, abot := StringTopBottom(ap, config)
		a, _ := geo.Intersection(top, atop)
		b, _ := geo.Intersection(bottom, abot)
		lines = append(lines, geo.Line{
			Start: a,
			End:   b,
		})
	}
	lines = append(lines, bottom)
	contour := MakeContour(lines)
	return contour
}

func StringTopBottom(pitch geo.Line, config config.Params) (geo.Line, geo.Line) {
	top := pitch.Offset(config["skirting"])
	bottom := pitch.Offset(config["skirting"] - config["stringer_width"])
	return top, bottom
}

func BottomLines() []geo.Line {
	ground := geo.Line{
		Start: geo.Point{X: 0, Y: 0},
		End:   geo.Point{X: 100, Y: 0},
	}
	front := geo.Line{
		Start: geo.Point{X: 0, Y: 0},
		End:   geo.Point{X: 0, Y: 100},
	}
	ls := make([]geo.Line, 2)
	ls[0] = ground
	ls[1] = front
	return ls
}

func JoinBelow(pitch []geo.Point, start geo.Point, top, bottom geo.Line, config config.Params, winder bool) geo.Line {
	// join with the stringer below
	offset := -1
	if winder {
		offset = -2
	}
	bp := geo.Line{
		Start: pitch[len(pitch)+offset],
		End:   start,
	}
	bt, bb := StringTopBottom(bp, config)
	a, _ := geo.Intersection(top, bt)
	b, _ := geo.Intersection(bottom, bb)
	return geo.Line{
		Start: a,
		End:   b,
	}
}

func JoinAbove(pitch []geo.Point, end geo.Point, top, bottom geo.Line, config config.Params) geo.Line {
	ap := geo.Line{
		Start: end,
		End:   pitch[len(pitch)-1],
	}
	atop, abot := StringTopBottom(ap, config)
	a, _ := geo.Intersection(top, atop)
	b, _ := geo.Intersection(bottom, abot)
	return geo.Line{
		Start: a,
		End:   b,
	}
}

func TopLines(stair []Section, config config.Params) []geo.Line {
	tr := TotalRise(stair, config)
	tg := OutsideGoing(stair, config) + config["bottom_horn"] + config["riser_thickness"] + config["nosing"]
	topFloor := geo.Line{
		Start: geo.Point{X: 0, Y: tr},
		End:   geo.Point{X: tg, Y: tr},
	}
	joist := geo.Line{
		Start: geo.Point{X: tg, Y: tr},
		End:   geo.Point{X: tg, Y: 0},
	}
	hornTop := topFloor.Offset(config["top_horn_height"])
	hornBack := joist.Offset(config["top_horn_length"])
	ls := make([]geo.Line, 4)
	ls[0] = hornTop
	ls[1] = hornBack
	ls[2] = topFloor
	ls[3] = joist
	return ls
}

func TotalRise(stair []Section, config config.Params) float64 {
	rise := 0.0
	for _, s := range stair {
		rise = rise + float64(s.Steps)*config["rise"]
	}
	return rise + config["rise"]
}

func OutsideGoing(stair []Section, config config.Params) float64 {
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

func ControlPoints(ts []geo.Point, stair []Section) [][]geo.Point {
	cps := make([][]geo.Point, 0)
	for i := 0; i < len(stair); i++ {
		end := 0
		for j := i; j >= 0; j-- {
			end = end + stair[j].Steps
		}
		start := end - stair[i].Steps
		ps := make([]geo.Point, 0)
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

func drawPoints(ps []geo.Point, ctx *canvas.Context, c color.Color) {
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

func toPath(ps []geo.Point) clip.Path {
	qs := make(clip.Path, len(ps))
	for i, p := range ps {
		qs[i] = clip.NewIntPointFromFloat(p.X*SCALE, p.Y*SCALE)
	}
	return qs
}

func fromPath(ps clip.Path) []geo.Point {
	qs := make([]geo.Point, len(ps))
	for i, p := range ps {
		dp := p.ToDoublePoint()
		qs[i] = geo.Point{X: dp.X / SCALE, Y: dp.Y / SCALE}
	}
	qs = append(qs, qs[0])
	return qs
}

func tr(tip geo.Point, cfg config.Params) ([]geo.Point, []geo.Point) {
	const (
		TC = 1000
		RC = 300
	)
	ts := make([]geo.Point, 0)
	ts = append(ts, tip)
	ts = append(ts, geo.Point{
		X: tip.X + TC + cfg["nosing"] + cfg["riser_thickness"],
		Y: tip.Y,
	})
	ts = append(ts, geo.Point{
		X: tip.X + TC + cfg["nosing"] + cfg["riser_thickness"],
		Y: tip.Y - math.Tan(cfg["wedge_angle"])*TC - cfg["tread_thickness"],
	})
	ts = append(ts, geo.Point{
		X: tip.X + cfg["nosing"] + cfg["riser_thickness"],
		Y: tip.Y - cfg["tread_thickness"],
	})
	ts = append(ts, geo.Point{
		X: tip.X,
		Y: tip.Y - cfg["tread_thickness"],
	})
	ts = append(ts, tip)
	rs := make([]geo.Point, 0)
	rTip := geo.Point{
		X: tip.X + cfg["nosing"],
		Y: tip.Y - cfg["tread_thickness"] + cfg["riser_rebate"],
	}
	rs = append(rs, rTip)
	rs = append(rs, geo.Point{
		X: rTip.X + cfg["riser_thickness"],
		Y: rTip.Y,
	})
	rs = append(rs, geo.Point{
		X: rTip.X + cfg["riser_thickness"] + RC*math.Tan(cfg["wedge_angle"]),
		Y: rTip.Y - cfg["tread_thickness"] - RC,
	})
	rs = append(rs, geo.Point{
		X: rTip.X,
		Y: rTip.Y - cfg["tread_thickness"] - RC,
	})
	rs = append(rs, rTip)

	return ts, rs
}

func Tips(sections []Section, config config.Params) ([]geo.Point, []geo.Point) {
	ps := make([]geo.Point, 0)
	at := geo.Point{X: config["bottom_horn"], Y: config["rise"]}
	ps = append(ps, at)

	// each section is responsible for amking the step that steps into the next section
	for _, s := range sections {
		if s.Kind == "straight" {
			for i := 1; i <= s.Steps; i++ {
				ps = append(ps, geo.Point{
					X: at.X + config["going"]*float64(i),
					Y: at.Y + config["rise"]*float64(i),
				})
			}
			at = ps[len(ps)-1]
		} else {
			angle := math.Pi / 2 / float64(s.Steps)
			if s.Steps == 3 {
				// step 1
				ps = append(ps, geo.Point{
					X: at.X + math.Tan(angle)*s.StartWidth,
					Y: at.Y + config["rise"],
				})
				ps = append(ps, geo.Point{
					X: at.X + s.EndWidth + s.StartWidth - s.StartWidth*math.Tan(angle),
					Y: at.Y + 2*config["rise"],
				})
				ps = append(ps, geo.Point{
					X: at.X + s.EndWidth + s.StartWidth,
					Y: at.Y + 3*config["rise"],
				})
				at = ps[len(ps)-1]
			}
		}
	}
	return ps, nil
}
