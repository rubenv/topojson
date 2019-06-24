package topojson

import (
	"math"

	"github.com/paulmach/orb"
)

func (t *Topology) bounds() {
	t.BoundingBox = []float64{
		math.MaxFloat64,
		math.MaxFloat64,
		-math.MaxFloat64,
		-math.MaxFloat64,
	}

	for _, f := range t.input {
		t.boundGeometry(f.Geometry)
	}

}

func (t *Topology) boundGeometry(g orb.Geometry) {
	switch c := g.(type) {
	case orb.Point:
		t.BBox(c.Bound())
	case orb.MultiPoint:
		t.BBox(c.Bound())
	case orb.LineString:
		t.BBox(c.Bound())
	case orb.MultiLineString:
		t.BBox(c.Bound())
	case orb.Polygon:
		t.BBox(c.Bound())
	case orb.MultiPolygon:
		t.BBox(c.Bound())
	case orb.Collection:
		t.BBox(c.Bound())
		// for _, geo := range c {
		// 	t.boundGeometry(geo)
		// }
	}
}

func (t *Topology) BBox(b orb.Bound) {
	xx := []float64{b.Min[0], b.Max[0]}
	yy := []float64{b.Min[1], b.Max[1]}
	for _, x := range xx {
		if x < t.BoundingBox[0] {
			t.BoundingBox[0] = x
		}
		if x > t.BoundingBox[2] {
			t.BoundingBox[2] = x
		}
	}
	for _, y := range yy {
		if y < t.BoundingBox[1] {
			t.BoundingBox[1] = y
		}
		if y > t.BoundingBox[3] {
			t.BoundingBox[3] = y
		}
	}
}

func (t *Topology) boundPoint(p []float64) {
	x := p[0]
	y := p[1]

	if x < t.BoundingBox[0] {
		t.BoundingBox[0] = x
	}
	if x > t.BoundingBox[2] {
		t.BoundingBox[2] = x
	}
	if y < t.BoundingBox[1] {
		t.BoundingBox[1] = y
	}
	if y > t.BoundingBox[3] {
		t.BoundingBox[3] = y
	}
}

func (t *Topology) boundPoints(l [][]float64) {
	for _, p := range l {
		t.boundPoint(p)
	}
}

func (t *Topology) boundMultiPoints(ml [][][]float64) {
	for _, l := range ml {
		for _, p := range l {
			t.boundPoint(p)
		}
	}
}
