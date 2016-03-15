package topojson

import "github.com/paulmach/go.geojson"

type Topology struct {
	input       []*inputGeometry
	coordinates [][]float64
	objects     map[string]*topologyObject
	lines       []*Arc
	rings       []*Arc
}

type Arc struct {
	Start int
	End   int
	Next  *Arc
}

type Point [2]float64

func NewPoint(coords []float64) Point {
	return Point{coords[0], coords[1]}
}

func PointEquals(a, b []float64) bool {
	return a[0] == b[0] && a[1] == b[1]
}

type inputGeometry struct {
	id   string
	geom *geojson.Geometry
}

type topologyObject struct {
	Type geojson.GeometryType

	Geometries []*topologyObject // For geometry collections
	Arc        *Arc              // For lines
	Arcs       []*Arc            // For multi lines and polygons
	MultiArcs  [][]*Arc          // For multi polygons
}
