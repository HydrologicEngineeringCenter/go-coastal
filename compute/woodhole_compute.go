package compute

import (
	"fmt"
	"log"
	"math"

	"github.com/HydrologicEngineeringCenter/go-coastal/hazardprovider"
	"github.com/HydrologicEngineeringCenter/go-coastal/resultswriters"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/structureprovider"
	"github.com/USACE/go-consequences/structures"
)

type WoodHoleSimulationSettings struct {
	DataSets        []WoodHoleFrequencyDataset
	BaseYear        int
	DiscountRate    float64
	InventoryPath   string
	OutputDirectory string //do a ead output per dataset and another one with discounted values and computed EEAD
	//terrain path?
	//occtype definitions?
	//life loss parameters?
	//seed?
}
type WoodHoleFrequencyDataset struct {
	Year                  int
	WaterSurfaceGridPaths []string
	//uncertainty data??
	WavePaths   []string
	Frequencies []float64
}

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

// CreateDiscountFactor discounts a future dollar based on a discount rate and a number of years in the future. The resulting
// value represents the value of a dollar that many years in the future at the specified discount rate. The value can be used
// as a factor to reduce future values to their present worth.
func CreateDiscountFactor(rate float64, numYearsInFuture int) float64 {
	//https://en.wikipedia.org/wiki/Discounted_cash_flow
	if numYearsInFuture <= 0 {
		return 1
	}
	return 1 / (math.Pow(1+rate, float64(numYearsInFuture))) //calcuation of a discount factor basising on 1 dollar to create a multiplier.
}
func WoodHoleMultiYearDeterministicEAD(settings WoodHoleSimulationSettings) error {
	inventory, err := structureprovider.InitGPK(settings.InventoryPath, "nsi")
	//create an aggregator results writer and inject it into the single compute
	if err != nil {
		return err
	}
	for _, d := range settings.DataSets {
		//make hazardproviders
		hps := make([]hazardproviders.HazardProvider, len(d.Frequencies))
		for idx, fp := range d.WaterSurfaceGridPaths {
			hps[idx] = hazardprovider.InitWoodHoleGroupTif(fp, d.WavePaths[idx])
		}
		//make results writer.
		rwfp := fmt.Sprintf("%vEAD_%v.gpkg", settings.OutputDirectory, d.Year)
		numYearsInFuture := d.Year - settings.BaseYear
		rw, err := resultswriters.InitwoodHoleResultsWriterFromFile(rwfp, d.Frequencies, CreateDiscountFactor(settings.DiscountRate, numYearsInFuture))
		if err != nil {
			return err
		}
		WoodHoleDeterministicEAD(hps, d.Frequencies, inventory, rw)
	}
	return nil
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
				depth := c.Depth() - s.GroundElevation
				if depth > 0 {
					c.SetDepth(depth) //set ground elevation on structures in go consequences, and pull from it to convert to depth... annoying.
					if err2 == nil {
						//hasLoss = true
						r, err := s.Compute(d)
						if err == nil {
							rw.Write(r)
						}
					}
				}

			}
		}

	})
}
