package topojson

import (
	geojson "github.com/paulmach/go.geojson"
)

func (t *Topology) ToGeoJSON() *geojson.FeatureCollection {
	fc := geojson.NewFeatureCollection()

	for _, obj := range t.Objects {
		switch obj.Type {
		case geojson.GeometryCollection:
			for _, geometry := range obj.Geometries {
				feat := geojson.NewFeature(t.toGeometry(geometry))
				feat.ID = geometry.ID
				feat.Properties = geometry.Properties
				feat.BoundingBox = geometry.BoundingBox
				fc.AddFeature(feat)
			}
		default:
			feat := geojson.NewFeature(t.toGeometry(obj))
			feat.ID = obj.ID
			feat.Properties = obj.Properties
			feat.BoundingBox = obj.BoundingBox
			fc.AddFeature(feat)
		}
	}

	return fc
}

func (t *Topology) toGeometry(g *Geometry) *geojson.Geometry {
	switch g.Type {
	case geojson.GeometryPoint:
		return geojson.NewPointGeometry(t.packPoint(g.Point))
	case geojson.GeometryMultiPoint:
		return geojson.NewMultiPointGeometry(t.packPoints(g.MultiPoint)...)
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

func (t *Topology) packPoint(in []float64) []float64 {
	if t.Transform == nil {
		return in
	}

	out := make([]float64, len(in))
	for i, v := range in {
		out[i] = v
		if i < 2 {
			out[i] = v*t.Transform.Scale[i] + t.Transform.Translate[i]
		}
	}

	return out
}

func (t *Topology) packPoints(in [][]float64) [][]float64 {
	out := make([][]float64, len(in))
	for i, p := range in {
		out[i] = t.packPoint(p)
	}
	return out
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

		// Copy arc
		newArc := make([][]float64, len(arc))
		for i, point := range arc {
			newArc[i] = append([]float64{}, point...)
		}

		if t.Transform != nil {
			x := float64(0)
			y := float64(0)

			for k, p := range newArc {
				x += p[0]
				y += p[1]

				newArc[k][0] = x*t.Transform.Scale[0] + t.Transform.Translate[0]
				newArc[k][1] = y*t.Transform.Scale[1] + t.Transform.Translate[1]
			}
		}

		if reverse {
			for j := len(newArc) - 1; j >= 0; j-- {
				if len(result) > 0 && pointEquals(result[len(result)-1], newArc[j]) {
					continue
				}
				result = append(result, newArc[j])
			}
		} else {
			for j := 0; j < len(newArc); j++ {
				if len(result) > 0 && pointEquals(result[len(result)-1], newArc[j]) {
					continue
				}
				result = append(result, newArc[j])
			}
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
