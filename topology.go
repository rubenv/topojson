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
	opts        *TopologyOptions
	input       []*geojson.Feature
	coordinates [][]float64
	objects     []*topologyObject
	lines       []*arc
	rings       []*arc
	arcs        []*arc
	arcIndexes  map[arcEntry]int
	deletedArcs map[int]bool
	shiftArcs   map[int]int
}

type Transform struct {
	Scale     [2]float64 `json:"scale"`
	Translate [2]float64 `json:"translate"`
}

type TopologyOptions struct {
	// Pre-quantization precision
	PreQuantize float64

	// Post-quantization precision
	PostQuantize float64

	// Maximum simplification error, set to 0 to disable
	Simplify float64

	// ID property key
	IDProperty string
}

func NewTopology(fc *geojson.FeatureCollection, opts *TopologyOptions) *Topology {
	if opts == nil {
		opts = &TopologyOptions{
			PreQuantize:  0,
			PostQuantize: 0,
			Simplify:     0,
			IDProperty:   "id",
		}
	}

	topo := &Topology{
		input: fc.Features,
		opts:  opts,
	}

	topo.bounds()
	topo.preQuantize()
	topo.extract()
	topo.join()
	topo.cut()
	topo.dedup()
	topo.unpackArcs()
	topo.simplify()
	topo.unpackObjects()
	topo.removeEmpty()
	topo.postQuantize()
	topo.delta()

	// No longer needed
	topo.opts = nil

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
	return a != nil && b != nil && a[0] == b[0] && a[1] == b[1]
}

type topologyObject struct {
	ID         string
	Type       geojson.GeometryType
	Properties map[string]interface{}

	Geometries []*topologyObject // For geometry collections
	Arc        *arc              // For lines
	Arcs       []*arc            // For multi lines and polygons
	MultiArcs  [][]*arc          // For multi polygons
}

// TODO: Verify that point and MultiPoint pass through
