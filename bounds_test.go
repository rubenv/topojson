package topojson

import (
	"testing"

	"github.com/cheekybits/is"
	geojson "github.com/paulmach/go.geojson"
)

func TestBoundingBox(t *testing.T) {
	is := is.New(t)

	in := []*geojson.Feature{
		NewTestFeature("foo", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})),
		NewTestFeature("bar", geojson.NewLineStringGeometry([][]float64{
			{-1, 0}, {1, 0}, {-2, 3},
		})),
	}

	topo := &Topology{input: in}
	topo.bounds()

	is.Equal(topo.BoundingBox, []float64{-2, 0, 2, 3})
}
