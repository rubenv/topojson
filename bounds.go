package topojson

import (
	"math"

	"github.com/paulmach/go.geojson"
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

func (t *Topology) boundGeometry(g *geojson.Geometry) {
	switch g.Type {
	case geojson.GeometryCollection:
		for _, geom := range g.Geometries {
			t.boundGeometry(geom)
		}
	case geojson.GeometryPoint:
		t.boundPoint(g.Point)
	case geojson.GeometryMultiPoint:
		t.boundPoints(g.MultiPoint)
	case geojson.GeometryLineString:
		t.boundPoints(g.LineString)
	case geojson.GeometryMultiLineString:
		t.boundMultiPoints(g.MultiLineString)
	case geojson.GeometryPolygon:
		t.boundMultiPoints(g.Polygon)
	case geojson.GeometryMultiPolygon:
		for _, poly := range g.MultiPolygon {
			t.boundMultiPoints(poly)
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
