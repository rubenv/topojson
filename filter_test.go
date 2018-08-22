package topojson

import (
	"testing"

	"github.com/cheekybits/is"
	geojson "github.com/paulmach/go.geojson"
)

func TestFilter(t *testing.T) {
	is := is.New(t)

	fc := geojson.NewFeatureCollection()
	fc.AddFeature(NewTestFeature("one", geojson.NewLineStringGeometry([][]float64{
		{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0},
	})))
	fc.AddFeature(NewTestFeature("two", geojson.NewLineStringGeometry([][]float64{
		{1, 0}, {2, 0}, {2, 1}, {1, 1}, {1, 0},
	})))
	fc.AddFeature(NewTestFeature("three", geojson.NewLineStringGeometry([][]float64{
		{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1},
	})))

	topo := NewTopology(fc, nil)
	is.NotNil(topo)

	al := len(topo.Arcs)
	is.True(al > 0)

	topo2 := topo.Filter([]string{"one", "two"})
	is.NotNil(topo2)

	al2 := len(topo2.Arcs)
	is.True(al > al2) // Arc has been eliminated

	expected := map[string][][]float64{
		"one": {{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
		"two": {{1, 0}, {2, 0}, {2, 1}, {1, 1}, {1, 0}},
	}

	fc2 := topo2.ToGeoJSON()
	is.NotNil(fc2)

	for _, feat := range fc2.Features {
		exp, ok := expected[feat.ID.(string)]
		is.True(ok)
		is.Equal(feat.Geometry.LineString, exp)
	}
}
