package topojson

import (
	"fmt"

	"github.com/paulmach/go.geojson"
)

func NewTestFeature(id string, geom *geojson.Geometry) *geojson.Feature {
	feature := geojson.NewFeature(geom)
	feature.SetProperty("id", id)
	return feature
}

func GetFeature(topo *Topology, id string) *topologyObject {
	for _, o := range topo.objects {
		if o.ID == id {
			return o
		}
	}
	panic(fmt.Sprintf("No such object: %s", id))
}
