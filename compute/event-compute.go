package compute

import (
	"fmt"
	"log"
	"strings"

	"github.com/HydrologicEngineeringCenter/go-coastal/hazardprovider"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	gcrw "github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structureprovider"
)

func Event(hazardfp string, inventoryfp string, frequency int, frequencystring string) {
	outputPathParts := strings.Split(hazardfp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}

	outfp += "_" + frequencystring + "_consequences.gpkg"
	//sw, err := gcrw.InitGeoJsonResultsWriterFromFile(outfp)
	sw, err := gcrw.InitSpatialResultsWriter(outfp, "results", "GPKG")
	if err != nil {
		panic("error creating ead output")
	}
	defer sw.Close()
	hp := hazardprovider.Init(hazardfp)
	hp.SelectFrequency(frequency - int(hazardprovider.Two)) //offset to zero based position.
	defer hp.Close()
	nsp, err := structureprovider.InitStructureProvider(inventoryfp, "nsi", "GPKG")
	if err != nil {
		panic("error creating ead output")
	}
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
			r, err := f.Compute(d)
			if err == nil {
				sw.Write(r)
			}
		}
	})
}
