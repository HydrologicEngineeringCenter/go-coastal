package compute

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/HydrologicEngineeringCenter/go-coastal/hazardprovider"
	"github.com/HydrologicEngineeringCenter/go-coastal/resultswriters"
	"github.com/HydrologicEngineeringCenter/go-coastal/structure_provider"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
)

func Event_Grid(hazardfp string, gridfp string, frequency int, frequencystring string, cellsize float64) {
	outputPathParts := strings.Split(hazardfp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}

	hp := hazardprovider.InitWithGrd(hazardfp, gridfp)
	hp.SelectFrequency(frequency - int(hazardprovider.Two)) //offset to zero based position.
	defer hp.Close()
	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	xdelta := cellsize
	ydelta := cellsize
	//get distance in x domain
	xdist := bbox.Bbox[0] - bbox.Bbox[2]

	//get distance in y domain
	ydist := bbox.Bbox[1] - bbox.Bbox[3]
	//get total number of x and y steps
	xSteps := int(math.Floor(math.Abs(xdist) / xdelta))
	ySteps := int(math.Floor(math.Abs(ydist) / ydelta))

	fmt.Println(bbox.ToString())
	outfp += "_" + frequencystring + "_grid.tif"
	sp, err := structure_provider.InitUniformPointDS(xdelta, ydelta)
	if err != nil {
		panic("error creating grid output")
	}
	sw, err := resultswriters.InitGridWriterFromFIle(outfp, xSteps, ySteps, bbox.Bbox[0], xdelta, bbox.Bbox[3], ydelta)
	if err != nil {
		panic("error creating grid output")
	}
	defer sw.Close()
	//get FilterStructures
	sp.ByBbox(bbox, func(f consequences.Receptor) {
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
func Event_Grid_new(outputfp string, hp hazardprovider.HazardProvider, cellsize float64) {
	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	xdelta := cellsize
	ydelta := cellsize
	//get distance in x domain
	xdist := bbox.Bbox[0] - bbox.Bbox[2]

	//get distance in y domain
	ydist := bbox.Bbox[1] - bbox.Bbox[3]
	//get total number of x and y steps
	xSteps := int(math.Floor(math.Abs(xdist) / xdelta))
	ySteps := int(math.Floor(math.Abs(ydist) / ydelta))

	fmt.Println(bbox.ToString())
	sp, err := structure_provider.InitUniformPointDS(xdelta, ydelta)
	if err != nil {
		panic("error creating grid output")
	}
	sw, err := resultswriters.InitGridWriterFromFIle(outputfp, xSteps, ySteps, bbox.Bbox[0], xdelta, bbox.Bbox[3], ydelta)
	if err != nil {
		panic("error creating grid output")
	}
	defer sw.Close()
	//get FilterStructures
	sp.ByBbox(bbox, func(f consequences.Receptor) {
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
