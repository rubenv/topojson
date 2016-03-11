package topojson

import (
	"reflect"
	"testing"

	"github.com/cheekybits/is"
	"github.com/paulmach/go.geojson"
)

// See https://github.com/mbostock/topojson/blob/master/test/topology/extract-test.js

// extract copies coordinates sequentially into a buffer
func TestCopiesCoordinates(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"foo", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})},
		{"bar", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})},
	}

	expected := [][]float64{
		{0, 0}, {1, 0}, {2, 0},
		{0, 0}, {1, 0}, {2, 0},
	}

	topo := &Topology{}
	topo.extract(in)
	is.Equal(len(topo.coordinates), len(expected))
	for k, v := range topo.coordinates {
		is.Equal(v, expected[k])
	}
}

// extract includes closing coordinates in polygons
func TestClosingCoordinates(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"foo", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0}, {0, 0},
		})},
	}

	expected := [][]float64{
		{0, 0}, {1, 0}, {2, 0}, {0, 0},
	}

	topo := &Topology{}
	topo.extract(in)
	is.Equal(len(topo.coordinates), len(expected))
	for k, v := range topo.coordinates {
		is.Equal(v, expected[k])
	}
}

// extract represents lines as contiguous slices of the coordinate buffer
func TestLineSlices(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"foo", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})},
		{"bar", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})},
	}

	topo := &Topology{}
	topo.extract(in)

	foo := topo.objects["foo"]
	is.Equal(foo.Type, geojson.GeometryLineString)
	is.True(reflect.DeepEqual(foo.Arc, Arc{0, 2}))

	bar := topo.objects["bar"]
	is.Equal(bar.Type, geojson.GeometryLineString)
	is.True(reflect.DeepEqual(bar.Arc, Arc{3, 5}))
}

// extract exposes the constructed lines and rings in the order of construction
func TestExtractRingsOrder(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"line", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})},
		{"multiline", geojson.NewMultiLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})},
		{"polygon", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {2, 0}, {0, 0},
			},
		})},
	}

	topo := &Topology{}
	topo.extract(in)

	is.True(reflect.DeepEqual(topo.lines, []Arc{
		{0, 2},
		{3, 5},
	}))
	is.True(reflect.DeepEqual(topo.rings, []Arc{
		{6, 9},
	}))
}

// extract supports nested geometry collections
func TestExtractNested(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"foo", geojson.NewCollectionGeometry(geojson.NewCollectionGeometry(geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {0, 1},
		})))},
	}

	topo := &Topology{}
	topo.extract(in)

	foo := topo.objects["foo"]
	is.Equal(foo.Type, geojson.GeometryCollection)

	geometries := foo.Geometries
	is.Equal(len(geometries), 1)
	is.Equal(geometries[0].Type, geojson.GeometryCollection)

	geometries = foo.Geometries[0].Geometries
	is.Equal(len(geometries), 1)
	is.Equal(geometries[0].Type, geojson.GeometryLineString)
	is.True(reflect.DeepEqual(geometries[0].Arc, Arc{0, 1}))
}
