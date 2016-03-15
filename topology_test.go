package topojson

import (
	"testing"

	"github.com/cheekybits/is"
	"github.com/paulmach/go.geojson"
)

func TestTopology(t *testing.T) {
	is := is.New(t)

	poly := geojson.NewPolygonFeature([][][]float64{
		{
			{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0},
		},
	})
	poly.SetProperty("id", "poly")

	fc := geojson.NewFeatureCollection()
	fc.AddFeature(poly)

	topo := NewTopology(fc, &TopologyOptions{
		Quantize: -1,
		Simplify: 0,
	})
	is.NotNil(topo)
	is.Equal(len(topo.Objects), 1)
	is.Equal(len(topo.Arcs), 1)
}
