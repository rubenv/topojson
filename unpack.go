package topojson

import "github.com/paulmach/go.geojson"

type arcEntry struct {
	Start int
	End   int
}

func (t *Topology) unpack() {
	arcIndexes := make(map[arcEntry]int)

	// Unpack arcs
	for i, a := range t.arcs {
		arcIndexes[arcEntry{a.Start, a.End}] = i
		t.Arcs = append(t.Arcs, t.coordinates[a.Start:a.End+1])
	}
	t.arcs = nil
	t.coordinates = nil

	// Unpack objects
	for _, o := range t.objects {
		t.Objects = append(t.Objects, t.unpackObject(arcIndexes, o))
	}
	t.objects = nil
}

func (t *Topology) unpackObject(arcs map[arcEntry]int, o *topologyObject) *Geometry {
	obj := &Geometry{
		ID:         o.ID,
		Type:       o.Type,
		Properties: o.Properties,
	}

	switch o.Type {
	case geojson.GeometryCollection:
		for _, geom := range o.Geometries {
			obj.Geometries = append(obj.Geometries, t.unpackObject(arcs, geom))
		}
	case geojson.GeometryLineString:
		obj.LineString = lookupArc(arcs, o.Arc)
	case geojson.GeometryMultiLineString:
		obj.MultiLineString = lookupArcs(arcs, o.Arcs)
	case geojson.GeometryPolygon:
		obj.Polygon = lookupArcs(arcs, o.Arcs)
	case geojson.GeometryMultiPolygon:
		obj.MultiPolygon = lookupMultiArcs(arcs, o.MultiArcs)
	}

	return obj
}

func lookupArc(arcs map[arcEntry]int, a *arc) []int {
	result := make([]int, 0)

	for a != nil {
		if a.Start < a.End {
			index := arcs[arcEntry{a.Start, a.End}]
			result = append(result, index)
		} else {
			index := arcs[arcEntry{a.End, a.Start}]
			result = append(result, -index)
		}
		a = a.Next
	}

	return result
}

func lookupArcs(arcs map[arcEntry]int, a []*arc) [][]int {
	result := make([][]int, 0)
	for _, arc := range a {
		result = append(result, lookupArc(arcs, arc))
	}
	return result
}

func lookupMultiArcs(arcs map[arcEntry]int, a [][]*arc) [][][]int {
	result := make([][][]int, 0)
	for _, s := range a {
		result = append(result, lookupArcs(arcs, s))
	}
	return result
}
