package topojson

import (
	"testing"
	"github.com/paulmach/go.geojson"
	"github.com/cheekybits/is"
)

func TestPointFeature(t *testing.T) {
	is := is.New(t)

	fc := geojson.NewFeatureCollection()
	f := geojson.NewPointFeature([]float64{0, 0})
	f.SetProperty("id", "point")
	fc.AddFeature(f)

	topo := NewTopology(fc, nil)

	is.Equal([]float64{0, 0}, topo.Objects["point"].Point)
}

func TestMultiPointFeature(t *testing.T) {
	is := is.New(t)

	fc := geojson.NewFeatureCollection()
	f := geojson.NewMultiPointFeature([]float64{0, 0}, []float64{1, 1})
	f.SetProperty("id", "multipoint")
	fc.AddFeature(f)

	topo := NewTopology(fc, nil)

	is.Equal([][]float64{{0, 0},{1, 1}}, topo.Objects["multipoint"].MultiPoint)
}