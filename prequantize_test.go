package topojson

import (
	"testing"

	"github.com/cheekybits/is"
	geojson "github.com/paulmach/go.geojson"
)

// Sets the quantization transform
func TestPreQuantize(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{}
	topo := &Topology{
		BoundingBox: []float64{0, 0, 1, 1},
		input:       in,
		opts: &TopologyOptions{
			PreQuantize:  1e4,
			PostQuantize: 1e4,
		},
	}

	topo.preQuantize()

	is.Equal(topo.Transform, &Transform{
		Scale:     [2]float64{float64(1) / 9999, float64(1) / 9999},
		Translate: [2]float64{0, 0},
	})
}

// Converts coordinates to fixed precision
func TestPreQuantizeConverts(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {0, 1}, {0, 0},
		})),
	}

	expected := [][]float64{
		{0, 0}, {9999, 0}, {0, 9999}, {0, 0},
	}

	topo := &Topology{
		input: in,
		opts: &TopologyOptions{
			PreQuantize:  1e4,
			PostQuantize: 1e4,
		},
	}

	topo.bounds()
	topo.preQuantize()

	is.Equal(topo.Transform, &Transform{
		Scale:     [2]float64{float64(1) / 9999, float64(1) / 9999},
		Translate: [2]float64{0, 0},
	})
	is.Equal(topo.input[0].Geometry.LineString, expected)
}

// Observes the quantization parameter
func TestPreQuantizeObserves(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {0, 1}, {0, 0},
		})),
	}

	expected := [][]float64{
		{0, 0}, {9, 0}, {0, 9}, {0, 0},
	}

	topo := &Topology{
		input: in,
		opts: &TopologyOptions{
			PreQuantize:  10,
			PostQuantize: 10,
		},
	}

	topo.bounds()
	topo.preQuantize()

	is.Equal(topo.input[0].Geometry.LineString, expected)
}

// Observes the bounding box
func TestPreQuantizeObservesBB(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {0, 1}, {0, 0},
		})),
	}
	topo := &Topology{
		BoundingBox: []float64{-1, -1, 2, 2},
		input:       in,
		opts: &TopologyOptions{
			PreQuantize:  10,
			PostQuantize: 10,
		},
	}

	topo.preQuantize()

	expected := [][]float64{
		{3, 3}, {6, 3}, {3, 6}, {3, 3},
	}
	is.Equal(topo.input[0].Geometry.LineString, expected)
}

// Applies to points as well as arcs
func TestPreQuantizeAppliesToPoints(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", geojson.NewMultiPointGeometry([][]float64{
			{0, 0}, {1, 0}, {0, 1}, {0, 0},
		}...)),
	}
	topo := &Topology{
		input: in,
		opts: &TopologyOptions{
			PreQuantize:  1e4,
			PostQuantize: 1e4,
		},
	}

	topo.bounds()
	topo.preQuantize()

	expected := [][]float64{
		{0, 0}, {9999, 0}, {0, 9999}, {0, 0},
	}

	is.Equal(topo.input[0].Geometry.MultiPoint, expected)
}

// Skips coincident points in lines
func TestPreQuantizeSkipsCoincidencesInLines(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {0.9, 0.9}, {1.1, 1.1}, {2, 2},
		})),
	}
	topo := &Topology{
		input: in,
		opts: &TopologyOptions{
			PreQuantize:  3,
			PostQuantize: 3,
		},
	}

	topo.bounds()
	topo.preQuantize()

	expected := [][]float64{
		{0, 0}, {1, 1}, {2, 2},
	}

	is.Equal(topo.input[0].Geometry.LineString, expected)
}

// Skips coincident points in polygons
func TestPreQuantizeSkipsCoincidencesInPolygons(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("polygon", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {0.9, 0.9}, {1.1, 1.1}, {2, 2}, {0, 0},
			},
		})),
	}
	topo := &Topology{
		input: in,
		opts: &TopologyOptions{
			PreQuantize:  3,
			PostQuantize: 3,
		},
	}

	topo.bounds()
	topo.preQuantize()

	expected := [][][]float64{
		{
			{0, 0}, {1, 1}, {2, 2}, {0, 0},
		},
	}

	is.Equal(topo.input[0].Geometry.Polygon, expected)
}

// Does not skip coincident points in points
func TestPreQuantizeDoesntSkipInPoints(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("multipoint", geojson.NewMultiPointGeometry([][]float64{
			{0, 0}, {0.9, 0.9}, {1.1, 1.1}, {2, 2}, {0, 0},
		}...)),
	}
	topo := &Topology{
		input: in,
		opts: &TopologyOptions{
			PreQuantize:  3,
			PostQuantize: 3,
		},
	}

	topo.bounds()
	topo.preQuantize()

	expected := [][]float64{
		{0, 0}, {1, 1}, {1, 1}, {2, 2}, {0, 0},
	}

	is.Equal(topo.input[0].Geometry.MultiPoint, expected)
}

// Includes closing point in degenerate lines
func TestPreQuantizeIncludesClosingLine(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", geojson.NewLineStringGeometry([][]float64{
			{1, 1}, {1, 1}, {1, 1},
		})),
	}
	topo := &Topology{
		BoundingBox: []float64{0, 0, 2, 2},
		input:       in,
		opts: &TopologyOptions{
			PreQuantize:  3,
			PostQuantize: 3,
		},
	}

	topo.preQuantize()

	expected := [][]float64{
		{1, 1}, {1, 1},
	}

	is.Equal(topo.input[0].Geometry.LineString, expected)
}

// Includes closing point in degenerate polygons
func TestPreQuantizeIncludesClosingPolygon(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("polygon", geojson.NewPolygonGeometry([][][]float64{
			{
				{0.9, 1}, {1.1, 1}, {1.01, 1}, {0.9, 1},
			},
		})),
	}
	topo := &Topology{
		BoundingBox: []float64{0, 0, 2, 2},
		input:       in,
		opts: &TopologyOptions{
			PreQuantize:  3,
			PostQuantize: 3,
		},
	}

	topo.preQuantize()

	expected := [][][]float64{
		{
			{1, 1}, {1, 1}, {1, 1}, {1, 1},
		},
	}

	is.Equal(topo.input[0].Geometry.Polygon, expected)
}
