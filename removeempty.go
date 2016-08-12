package topojson

import geojson "github.com/paulmach/go.geojson"

func (t *Topology) removeEmpty() {
	objs := make([]*Geometry, 0, len(t.Objects))
	for _, o := range t.Objects {
		obj := t.removeEmptyObjects(o)
		if obj != nil {
			objs = append(objs, obj)
		}
	}
	t.Objects = objs
}

func (t *Topology) removeEmptyObjects(obj *Geometry) *Geometry {
	switch obj.Type {
	case geojson.GeometryCollection:
		geoms := make([]*Geometry, 0, len(obj.Geometries))
		for _, g := range obj.Geometries {
			geom := t.removeEmptyObjects(g)
			if geom != nil {
				geoms = append(geoms, geom)
			}
		}
		if len(geoms) == 0 {
			return nil
		}
		obj.Geometries = geoms
	case geojson.GeometryLineString:
		if len(obj.LineString) == 0 {
			return nil
		}
	case geojson.GeometryMultiLineString:
		linestrings := make([][]int, 0, len(obj.MultiLineString))
		for _, ls := range obj.MultiLineString {
			if len(ls) > 0 {
				linestrings = append(linestrings, ls)
			}
		}
		if len(linestrings) == 0 {
			return nil
		}
		obj.MultiLineString = linestrings
	case geojson.GeometryPolygon:
		rings := t.removeEmptyPolygon(obj.Polygon)
		if rings == nil {
			return nil
		}
		obj.Polygon = rings
	case geojson.GeometryMultiPolygon:
		polygons := make([][][]int, 0, len(obj.MultiPolygon))
		for _, polygon := range obj.MultiPolygon {
			rings := t.removeEmptyPolygon(polygon)
			if rings != nil {
				polygons = append(polygons, rings)
			}
		}
		if len(polygons) == 0 {
			return nil
		}
		obj.MultiPolygon = polygons
	}

	return obj
}

func (t *Topology) removeEmptyPolygon(polygon [][]int) [][]int {
	rings := make([][]int, 0, len(polygon))
	for i, ring := range polygon {
		if i == 0 && len(ring) == 0 {
			return nil // Outer ring empty
		}
		if len(ring) > 0 {
			rings = append(rings, ring)
		}
	}
	return rings
}
