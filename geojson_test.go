package topojson

import (
	"testing"

	"github.com/cheekybits/is"
	geojson "github.com/paulmach/go.geojson"
)

func TestGeoJSON(t *testing.T) {
	is := is.New(t)

	poly := geojson.NewPolygonFeature([][][]float64{
		{
			{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0},
		},
	})
	poly.ID = "poly"
	poly.SetProperty("id", "poly")

	fc := geojson.NewFeatureCollection()
	fc.AddFeature(poly)

	topo := NewTopology(fc, nil)
	is.NotNil(topo)
	is.Equal(len(topo.Objects), 1)
	is.Equal(len(topo.Arcs), 1)

	fc2 := topo.ToGeoJSON()
	is.NotNil(fc2)
	is.Equal(fc, fc2)
}
