package topojson

import (
	"testing"

	"github.com/cheekybits/is"
	"github.com/paulmach/go.geojson"
)

// See https://github.com/mbostock/topojson/blob/master/test/topology/join-test.js

// join the returned hashmap has true for junction points
func TestHasJunctions(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"cba", geojson.NewLineStringGeometry([][]float64{
			{2, 0}, {1, 0}, {0, 0},
		})},
		{"ab", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0},
		})},
	}

	topo := &Topology{}
	topo.extract(in)
	junctions := topo.join()

	is.True(junctions.Has([]float64{2, 0}))
	is.True(junctions.Has([]float64{0, 0}))
}

// join the returned hashmap has undefined for non-junction points
func TestNonJunctions(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"cba", geojson.NewLineStringGeometry([][]float64{
			{2, 0}, {1, 0}, {0, 0},
		})},
		{"ab", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {2, 0},
		})},
	}

	topo := &Topology{}
	topo.extract(in)
	junctions := topo.join()

	is.False(junctions.Has([]float64{1, 0}))
}

// join exact duplicate lines ABC & ABC have junctions at their end points
func TestJoinDuplicate(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"abc", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})},
		{"abc2", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})},
	}

	topo := &Topology{}
	topo.extract(in)
	junctions := topo.join()

	is.Equal(len(junctions), 2)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{2, 0}))
}

// join reversed duplicate lines ABC & CBA have junctions at their end points"
func TestJoinReversedDuplicate(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"abc", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})},
		{"cba", geojson.NewLineStringGeometry([][]float64{
			{2, 0}, {1, 0}, {0, 0},
		})},
	}

	topo := &Topology{}
	topo.extract(in)
	junctions := topo.join()

	is.Equal(len(junctions), 2)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{2, 0}))
}

// join exact duplicate rings ABCA & ABCA have no junctions
func TestJoinDuplicateRings(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"abca", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {2, 0}, {0, 0},
			},
		})},
		{"abca2", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {2, 0}, {0, 0},
			},
		})},
	}

	topo := &Topology{}
	topo.extract(in)
	junctions := topo.join()

	is.Equal(len(junctions), 0)
}

// join reversed duplicate rings ACBA & ABCA have no junctions
func TestJoinReversedDuplicateRings(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"abca", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {2, 0}, {0, 0},
			},
		})},
		{"acba", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {2, 0}, {1, 0}, {0, 0},
			},
		})},
	}

	topo := &Topology{}
	topo.extract(in)
	junctions := topo.join()

	is.Equal(len(junctions), 0)
}

// join rotated duplicate rings BCAB & ABCA have no junctions
func TestJoinRotatedDuplicateRings(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"abca", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {2, 0}, {0, 0},
			},
		})},
		{"bcab", geojson.NewPolygonGeometry([][][]float64{
			{
				{1, 0}, {2, 0}, {0, 0}, {1, 0},
			},
		})},
	}

	topo := &Topology{}
	topo.extract(in)
	junctions := topo.join()

	is.Equal(len(junctions), 0)
}

// join ring ABCA & line ABCA have a junction at A
func TestJoinRingLine(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"abcaLine", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0}, {0, 0},
		})},
		{"abcaPolygon", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {2, 0}, {0, 0},
			},
		})},
	}

	topo := &Topology{}
	topo.extract(in)
	junctions := topo.join()

	is.Equal(len(junctions), 1)
	is.True(junctions.Has([]float64{0, 0}))
}

// join ring BCAB & line ABCA have a junction at A
func TestJoinLineRingReversed(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"abcaLine", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0}, {0, 0},
		})},
		{"bcabPolygon", geojson.NewPolygonGeometry([][][]float64{
			{
				{1, 0}, {2, 0}, {0, 0}, {1, 0},
			},
		})},
	}

	topo := &Topology{}
	topo.extract(in)
	junctions := topo.join()

	is.Equal(len(junctions), 1)
	is.True(junctions.Has([]float64{0, 0}))
}

// join ring ABCA & line BCAB have a junction at B
func TestJoinRingLineReversed(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"bcabLine", geojson.NewLineStringGeometry([][]float64{
			{1, 0}, {2, 0}, {0, 0}, {1, 0},
		})},
		{"abcaPolygon", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {2, 0}, {0, 0},
			},
		})},
	}

	topo := &Topology{}
	topo.extract(in)
	junctions := topo.join()

	is.Equal(len(junctions), 1)
	is.True(junctions.Has([]float64{1, 0}))
}

// join when an old arc ABC extends a new arc AB, there is a junction at B
func TestJoinOldArcExtends(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"abc", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})},
		{"ab", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0},
		})},
	}

	topo := &Topology{}
	topo.extract(in)
	junctions := topo.join()

	is.Equal(len(junctions), 3)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{2, 0}))
}

// join when a reversed old arc CBA extends a new arc AB, there is a junction at B
func TestJoinOldArcExtendsReversed(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"cba", geojson.NewLineStringGeometry([][]float64{
			{2, 0}, {1, 0}, {0, 0},
		})},
		{"ab", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0},
		})},
	}

	topo := &Topology{}
	topo.extract(in)
	junctions := topo.join()

	is.Equal(len(junctions), 3)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{1, 0}))
	is.True(junctions.Has([]float64{2, 0}))
}

// join when a new arc ADE shares its start with an old arc ABC, there is a junction at A
func TestJoinNewArcSharesStart(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"ade", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 0}, {2, 0},
		})},
		{"abc", geojson.NewLineStringGeometry([][]float64{
			{0, 0}, {1, 1}, {2, 1},
		})},
	}

	topo := &Topology{}
	topo.extract(in)
	junctions := topo.join()

	is.Equal(len(junctions), 3)
	is.True(junctions.Has([]float64{0, 0}))
	is.True(junctions.Has([]float64{2, 0}))
	is.True(junctions.Has([]float64{2, 1}))
}

// join ring ABA has no junctions
func TestJoinRingNoJunctions(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"aba", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {1, 0}, {0, 0},
			},
		})},
	}

	topo := &Topology{}
	topo.extract(in)
	junctions := topo.join()

	is.Equal(len(junctions), 0)
}

// join ring AA has no junctions
func TestJoinRingAANoJunctions(t *testing.T) {
	is := is.New(t)

	in := []*inputGeometry{
		{"aa", geojson.NewPolygonGeometry([][][]float64{
			{
				{0, 0}, {0, 0},
			},
		})},
	}

	topo := &Topology{}
	topo.extract(in)
	junctions := topo.join()

	is.Equal(len(junctions), 0)
}
