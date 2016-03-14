package topojson

import "github.com/paulmach/go.geojson"

type inputGeometry struct {
	id   string
	geom *geojson.Geometry
}

func (t *Topology) extract(in []*inputGeometry) {
	t.objects = make(map[string]*topologyObject)

	for _, g := range in {
		t.objects[g.id] = t.extractGeometry(g.geom)
	}
}

func (t *Topology) extractGeometry(g *geojson.Geometry) *topologyObject {
	o := &topologyObject{
		Type: g.Type,
	}

	switch g.Type {
	case geojson.GeometryCollection:
		for _, geom := range g.Geometries {
			o.Geometries = append(o.Geometries, t.extractGeometry(geom))
		}
	case geojson.GeometryLineString:
		o.Arc = t.extractLine(g.LineString)
	case geojson.GeometryMultiLineString:
		o.Arcs = make([]*Arc, len(g.MultiLineString))
		for i, l := range g.MultiLineString {
			o.Arcs[i] = t.extractLine(l)
		}
	case geojson.GeometryPolygon:
		o.Arcs = make([]*Arc, len(g.Polygon))
		for i, r := range g.Polygon {
			o.Arcs[i] = t.extractRing(r)
		}
	}

	return o
}

func (t *Topology) extractLine(line [][]float64) *Arc {
	n := len(line)
	for i := 0; i < n; i++ {
		t.coordinates = append(t.coordinates, line[i])
	}

	index := len(t.coordinates) - 1
	arc := &Arc{Start: index - n + 1, End: index}
	t.lines = append(t.lines, arc)

	return arc
}

func (t *Topology) extractRing(ring [][]float64) *Arc {
	n := len(ring)
	for i := 0; i < n; i++ {
		t.coordinates = append(t.coordinates, ring[i])
	}

	index := len(t.coordinates) - 1
	arc := &Arc{Start: index - n + 1, End: index}
	t.rings = append(t.rings, arc)

	return arc
}
