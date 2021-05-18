package compute

import (
	"fmt"
	"log"
	"strings"

	"github.com/HydrologicEngineeringCenter/go-coastal/hazardprovider"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/structureprovider"
)

func Event(hazardfp string, inventoryfp string, frequency int) {
	outputPathParts := strings.Split(hazardfp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}
	outfp += "_consequences.json"
	sw := consequences.InitGeoJsonResultsWriterFromFile(outfp)
	defer sw.Close()
	hp := hazardprovider.Init(hazardfp, frequency) //pass in frequency?
	defer hp.Close()
	nsp := structureprovider.InitGPK(inventoryfp, "nsi")
	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	//get FilterStructures
	nsp.ByBbox(bbox, func(f consequences.Receptor) {
		//ProvideHazard works off of a geography.Location
		d, err2 := hp.ProvideHazard(geography.Location{X: f.Location().X, Y: f.Location().Y})
		//compute damages based on hazard being able to provide depth
		if err2 == nil {
			if d.Has(hazards.Depth) {
				if d.Depth() > 0.0 && d.Depth() < 9999.0 {
					sw.Write(f.Compute(d))
				}
			}
		}
	})

}
