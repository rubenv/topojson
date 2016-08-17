package topojson

import "math"

type quantize struct {
	Transform *Transform

	dx, dy, kx, ky float64
}

func newQuantize(dx, dy, kx, ky float64) *quantize {
	return &quantize{
		dx: dx,
		dy: dy,
		kx: kx,
		ky: ky,

		Transform: &Transform{
			Scale:     [2]float64{1 / kx, 1 / ky},
			Translate: [2]float64{-dx, -dy},
		},
	}
}

func (q *quantize) quantizePoint(p []float64) []float64 {
	x := round((p[0] + q.dx) * q.kx)
	y := round((p[1] + q.dy) * q.ky)
	return []float64{x, y}
}

func (q *quantize) quantizeLine(in [][]float64, skipEqual bool) [][]float64 {
	out := make([][]float64, 0, len(in))

	var last []float64

	for _, p := range in {
		pt := q.quantizePoint(p)
		if !pointEquals(pt, last) || !skipEqual {
			out = append(out, pt)
			last = pt
		}
	}

	if len(out) < 2 {
		out = append(out, out[0])
	}

	return out
}

func (q *quantize) quantizeMultiLine(in [][][]float64, skipEqual bool) [][][]float64 {
	out := make([][][]float64, len(in))
	for i, line := range in {
		line = q.quantizeLine(line, skipEqual)
		for len(line) < 4 {
			line = append(line, line[0])
		}
		out[i] = line
	}

	return out
}

func round(v float64) float64 {
	if v < 0 {
		return math.Ceil(v - 0.5)
	} else {
		return math.Floor(v + 0.5)
	}
}
