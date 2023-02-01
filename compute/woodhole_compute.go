package compute

import (
	"fmt"
	"log"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
)

func WoodHoleEvent(hp hazardproviders.HazardProvider, sp consequences.StreamProvider, rw consequences.ResultsWriter) {
	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	//get FilterStructures
	sp.ByBbox(bbox, func(f consequences.Receptor) {
		//ProvideHazard works off of a geography.Location
		d, err2 := hp.ProvideHazard(geography.Location{X: f.Location().X, Y: f.Location().Y})
		//compute damages based on hazard being able to provide a hazard.
		if err2 == nil {
			r, err := f.Compute(d)
			if err == nil {
				rw.Write(r)
			}
		} else {
			fmt.Println(err2.Error())
		}
	})
}
