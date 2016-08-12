package topojson

import (
	geojson "github.com/paulmach/go.geojson"
)

func (t *Topology) ToGeoJSON() *geojson.FeatureCollection {
	fc := geojson.NewFeatureCollection()

	for _, obj := range t.Objects {
		feat := geojson.NewFeature(t.toGeometry(obj))
		feat.ID = obj.ID
		feat.Properties = obj.Properties
		fc.AddFeature(feat)
	}

	return fc
}

func (t *Topology) toGeometry(g *Geometry) *geojson.Geometry {
	switch g.Type {
	case geojson.GeometryPoint:
		return geojson.NewPointGeometry(g.Point)
	case geojson.GeometryMultiPoint:
		return geojson.NewMultiPointGeometry(g.MultiPoint...)
	case geojson.GeometryLineString:
		return geojson.NewLineStringGeometry(t.packLinestring(g.LineString))
	case geojson.GeometryMultiLineString:
		return geojson.NewMultiLineStringGeometry(t.packMultiLinestring(g.MultiLineString)...)
	case geojson.GeometryPolygon:
		return geojson.NewPolygonGeometry(t.packMultiLinestring(g.Polygon))
	case geojson.GeometryMultiPolygon:
		polygons := make([][][][]float64, len(g.MultiPolygon))
		for i, poly := range g.MultiPolygon {
			polygons[i] = t.packMultiLinestring(poly)
		}
		return geojson.NewMultiPolygonGeometry(polygons...)
	case geojson.GeometryCollection:
		geometries := make([]*geojson.Geometry, len(g.Geometries))
		for i, geometry := range g.Geometries {
			geometries[i] = t.toGeometry(geometry)
		}
		return geojson.NewCollectionGeometry(geometries...)
	}
	return nil
}

func (t *Topology) packLinestring(ls []int) [][]float64 {
	result := make([][]float64, 0)
	for _, a := range ls {
		reverse := false
		if a < 0 {
			a = ^a
			reverse = true
		}
		arc := t.Arcs[a]

		if reverse {
			for j := len(arc) - 1; j >= 0; j-- {
				result = append(result, arc[j])
			}
		} else {
			result = append(result, arc...)
		}
	}
	return result
}

func (t *Topology) packMultiLinestring(ls [][]int) [][][]float64 {
	result := make([][][]float64, len(ls))
	for i, l := range ls {
		result[i] = t.packLinestring(l)
	}
	return result
}
