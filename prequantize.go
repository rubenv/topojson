package topojson

import geojson "github.com/paulmach/go.geojson"

func (t *Topology) preQuantize() {
	if t.opts.PreQuantize == 0 {
		return
	}
	if t.opts.PostQuantize == 0 {
		t.opts.PostQuantize = t.opts.PreQuantize
	}

	q0 := t.opts.PreQuantize
	q1 := t.opts.PostQuantize

	x0 := t.BoundingBox[0]
	y0 := t.BoundingBox[1]
	x1 := t.BoundingBox[2]
	y1 := t.BoundingBox[3]

	kx := float64(1)
	if x1-x0 != 0 {
		kx = (q1 - 1) / (x1 - x0) * q0 / q1
	}

	ky := float64(1)
	if y1-y0 != 0 {
		ky = (q1 - 1) / (y1 - y0) * q0 / q1
	}

	q := newQuantize(-x0, -y0, kx, ky)

	for _, f := range t.input {
		t.preQuantizeGeometry(q, f.Geometry)
	}

	t.Transform = q.Transform
}

func (t *Topology) preQuantizeGeometry(q *quantize, g *geojson.Geometry) {
	switch g.Type {
	case geojson.GeometryCollection:
		for _, geom := range g.Geometries {
			t.preQuantizeGeometry(q, geom)
		}
	case geojson.GeometryPoint:
		g.Point = q.quantizePoint(g.Point)
	case geojson.GeometryMultiPoint:
		g.MultiPoint = q.quantizeLine(g.MultiPoint, false)
	case geojson.GeometryLineString:
		g.LineString = q.quantizeLine(g.LineString, true)
	case geojson.GeometryMultiLineString:
		g.MultiLineString = q.quantizeMultiLine(g.MultiLineString, true)
	case geojson.GeometryPolygon:
		g.Polygon = q.quantizeMultiLine(g.Polygon, true)
	case geojson.GeometryMultiPolygon:
		for i, poly := range g.MultiPolygon {
			g.MultiPolygon[i] = q.quantizeMultiLine(poly, true)
		}
	}
}
