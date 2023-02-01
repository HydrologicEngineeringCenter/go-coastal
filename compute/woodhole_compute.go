package compute

import (
	"fmt"
	"log"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
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

func WoodHoleDeterministicEAD(hps []hazardproviders.HazardProvider, frequencies []float64, sp consequences.StreamProvider, rw consequences.ResultsWriter) {
	fmt.Println("Getting bbox")
	bbox, err := hps[len(hps)-1].ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	header := []string{"fd_id", "x", "y", "damage category", "occupancy type", "EAD structure", "EAD content", "pop2amu65", "pop2amo65", "pop2pmu65", "pop2pmo65", "cbfips", "s_dam_per", "c_dam_per"}
	results := []interface{}{"updateme", 0.0, 0.0, "dc", "ot", 0.0, 0.0, 0, 0, 0, 0, "CENSUSBLOCKFIPS", 0, 0}
	//get FilterStructures
	sp.ByBbox(bbox, func(f consequences.Receptor) {
		ret := consequences.Result{
			Headers: header,
			Result:  results,
		}
		sdamages := make([]float64, len(frequencies))
		cdamages := make([]float64, len(frequencies))
		s, sdok := f.(structures.StructureDeterministic)
		//s.SampleStructure()
		if sdok {
			hasLoss := false
			for idx, _ := range frequencies {
				if idx == 0 {
					ret.Result[0] = s.BaseStructure.Name
					ret.Result[1] = s.BaseStructure.X
					ret.Result[2] = s.BaseStructure.Y
					ret.Result[3] = s.BaseStructure.DamCat
					ret.Result[4] = s.OccType.Name
					ret.Result[7] = s.Pop2amu65
					ret.Result[8] = s.Pop2amo65
					ret.Result[9] = s.Pop2pmu65
					ret.Result[10] = s.Pop2pmo65
					ret.Result[11] = s.CBFips
				}
				//ProvideHazard works off of a geography.Location
				d, err2 := hps[idx].ProvideHazard(geography.Location{X: f.Location().X, Y: f.Location().Y})
				//compute damages based on hazard being able to provide a hazard.
				if err2 == nil {
					hasLoss = true
					r, err := s.Compute(d)
					if err == nil {
						rs, err := r.Fetch("structure damage")
						if err != nil {
							panic(err)
						}
						rc, err := r.Fetch("content damage")
						if err != nil {
							panic(err)
						}
						sdamages[idx] = rs.(float64)
						cdamages[idx] = rc.(float64)
					}
				} else {
					fmt.Println(err2.Error())
					sdamages[idx] = 0.0
					cdamages[idx] = 0.0
				}
			}
			if hasLoss {
				sead := compute.ComputeSpecialEAD(sdamages, frequencies)
				cead := compute.ComputeSpecialEAD(cdamages, frequencies)
				ret.Result[5] = sead
				ret.Result[6] = cead
				rw.Write(ret)
			}

		}

	})
}
