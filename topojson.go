package topojson

import "github.com/paulmach/go.geojson"

type Topology struct {
	coordinates [][]float64
	objects     map[string]*topologyObject
	lines       []Arc
	rings       []Arc
}

type Arc [2]int

type Point [2]float64

type topologyObject struct {
	Type geojson.GeometryType

	Geometries []*topologyObject // For geometry collections
	Arc        Arc               // For lines
	Arcs       []Arc             // For multi lines and polygons
	MultiArcs  [][]Arc           // For multi polygons
}
