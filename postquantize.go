package topojson

import geojson "github.com/paulmach/go.geojson"

func (t *Topology) postQuantize() {
	q0 := t.opts.PreQuantize
	q1 := t.opts.PostQuantize

	if q1 == 0 {
		return
	}

	var q *quantize

	if q0 != 0 {
		if q0 == q1 {
			return
		}

		k := q1 / q0

		q = newQuantize(0, 0, k, k)

		t.Transform.Scale[0] /= k
		t.Transform.Scale[1] /= k
	} else {
		x0 := t.BoundingBox[0]
		y0 := t.BoundingBox[1]
		x1 := t.BoundingBox[2]
		y1 := t.BoundingBox[3]

		kx := float64(1)
		if x1-x0 != 0 {
			kx = (q1 - 1) / (x1 - x0)
		}

		ky := float64(1)
		if y1-y0 != 0 {
			ky = (q1 - 1) / (y1 - y0)
		}

		q := newQuantize(-x0, -y0, kx, ky)
		t.Transform = q.Transform
	}

	for _, f := range t.input {
		t.postQuantizeGeometry(q, f.Geometry)
	}

	for i, arc := range t.Arcs {
		t.Arcs[i] = q.quantizeLine(arc, true)
	}
}

func (t *Topology) postQuantizeGeometry(q *quantize, g *geojson.Geometry) {
	switch g.Type {
	case geojson.GeometryCollection:
		for _, geom := range g.Geometries {
			t.postQuantizeGeometry(q, geom)
		}
	case geojson.GeometryPoint:
		g.Point = q.quantizePoint(g.Point)
	case geojson.GeometryMultiPoint:
		g.MultiPoint = q.quantizeLine(g.MultiPoint, false)
	}
}
