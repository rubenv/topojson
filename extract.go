package topojson

import "github.com/paulmach/go.geojson"

func (t *Topology) extract() {
	t.objects = make([]*topologyObject, 0, len(t.input))

	for _, g := range t.input {
		t.objects = append(t.objects, t.extractFeature(g))
	}
	t.input = nil // no longer needed
}

func (t *Topology) extractFeature(f *geojson.Feature) *topologyObject {
	g := f.Geometry
	o := t.extractGeometry(g)

	idProp := "id"
	if t.opts != nil && t.opts.IDProperty != "" {
		idProp = t.opts.IDProperty
	}

	id, err := f.PropertyString(idProp)
	if err == nil {
		o.ID = id
	}

	o.Properties = f.Properties

	return o
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
		o.Arcs = make([]*arc, len(g.MultiLineString))
		for i, l := range g.MultiLineString {
			o.Arcs[i] = t.extractLine(l)
		}
	case geojson.GeometryPolygon:
		o.Arcs = make([]*arc, len(g.Polygon))
		for i, r := range g.Polygon {
			o.Arcs[i] = t.extractRing(r)
		}
	}

	return o
}

func (t *Topology) extractLine(line [][]float64) *arc {
	n := len(line)
	for i := 0; i < n; i++ {
		t.coordinates = append(t.coordinates, line[i])
	}

	index := len(t.coordinates) - 1
	arc := &arc{Start: index - n + 1, End: index}
	t.lines = append(t.lines, arc)

	return arc
}

func (t *Topology) extractRing(ring [][]float64) *arc {
	n := len(ring)
	for i := 0; i < n; i++ {
		t.coordinates = append(t.coordinates, ring[i])
	}

	index := len(t.coordinates) - 1
	arc := &arc{Start: index - n + 1, End: index}
	t.rings = append(t.rings, arc)

	return arc
}
