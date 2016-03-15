package topojson

import (
	"encoding/json"

	"github.com/paulmach/go.geojson"
)

type Topology struct {
	Type      string     `json:"type"`
	Transform *Transform `json:"transform,omitempty"`

	BoundingBox []float64     `json:"bbox,omitempty"`
	Objects     []*Geometry   `json:"objects"`
	Arcs        [][][]float64 `json:"arcs"`

	// For internal use only
	input       []*inputGeometry
	coordinates [][]float64
	objects     map[string]*topologyObject
	lines       []*arc
	rings       []*arc
}

type Transform struct {
	Scale     [2]float64 `json:"scale"`
	Translate [2]float64 `json:"translate"`
}

func NewTopology(features *geojson.FeatureCollection) *Topology {
	topo := &Topology{
		input: nil, // TODO
	}

	topo.extract()
	topo.join()
	topo.cut()
	topo.dedup()

	return topo
}

// MarshalJSON converts the topology object into the proper JSON.
// It will handle the encoding of all the child geometries.
// Alternately one can call json.Marshal(t) directly for the same result.
func (t *Topology) MarshalJSON() ([]byte, error) {
	t.Type = "Topology"
	if t.Objects == nil {
		t.Objects = make([]*Geometry, 0) // TopoJSON requires the objects attribute to be at least []
	}
	if t.Arcs == nil {
		t.Arcs = make([][][]float64, 0) // TopoJSON requires the arcs attribute to be at least []
	}
	return json.Marshal(*t)
}

// Internal structs

type arc struct {
	Start int
	End   int
	Next  *arc
}

type point [2]float64

func newPoint(coords []float64) point {
	return point{coords[0], coords[1]}
}

func pointEquals(a, b []float64) bool {
	return a[0] == b[0] && a[1] == b[1]
}

type inputGeometry struct {
	id   string
	geom *geojson.Geometry
}

type topologyObject struct {
	Type geojson.GeometryType

	Geometries []*topologyObject // For geometry collections
	Arc        *arc              // For lines
	Arcs       []*arc            // For multi lines and polygons
	MultiArcs  [][]*arc          // For multi polygons
}
