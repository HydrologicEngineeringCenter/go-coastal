package compute

import (
	"fmt"
	"log"

	"github.com/HydrologicEngineeringCenter/go-coastal/resultswriters"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/structures"
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

func WoodHoleDeterministicEAD(hps []hazardproviders.HazardProvider, frequencies []float64, sp consequences.StreamProvider, rw *resultswriters.WoodHoleResultsWriter) {
	fmt.Println("Getting bbox")
	bbox, err := hps[len(hps)-1].ProvideHazardBoundary() //get the biggest depth grid.
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	//get FilterStructures
	sp.ByBbox(bbox, func(f consequences.Receptor) {
		//sdamages := make([]float64, len(frequencies))
		//cdamages := make([]float64, len(frequencies))
		s, sdok := f.(structures.StructureDeterministic)
		//s.SampleStructure()
		if sdok {
			for idx, _ := range frequencies {
				rw.UpdateFrequencyIndex(idx)
				//ProvideHazard works off of a geography.Location
				d, err2 := hps[idx].ProvideHazard(geography.Location{X: f.Location().X, Y: f.Location().Y})
				//compute damages based on hazard being able to provide a hazard.
				c := d.(hazards.CoastalEvent)
				c.SetDepth(c.Depth() - s.GroundElevation) //set ground elevation on structures in go consequences, and pull from it to convert to depth... annoying.
				if err2 == nil {
					//hasLoss = true
					r, err := s.Compute(d)
					if err == nil {
						rw.Write(r)
					}
				}
			}
		}

	})
}
