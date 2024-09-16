package compute

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/HydrologicEngineeringCenter/go-coastal/hazardprovider"
	"github.com/HydrologicEngineeringCenter/go-coastal/resultswriters"
	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/criticalinfrastructure"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/lifeloss"
	gcrw "github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structureprovider"
	"github.com/USACE/go-consequences/structures"
	"github.com/USACE/go-consequences/warning"
)

func initLifeLossEngine(complianceRate float64, rng *rand.Rand) lifeloss.LifeLossEngine {
	ws := warning.InitComplianceBasedWarningSystem(rng.Int63(), complianceRate)
	return lifeloss.Init(rng.Int63(), ws)
}
func ExpectedAnnualDamages(hazardfp string, grdfp string, inventoryfp string, complianceRate float64, seed int64) {
	outputPathParts := strings.Split(hazardfp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}
	outfp += "_ead_consequences.gpkg"
	sw, err := gcrw.InitSpatialResultsWriter(outfp, "results", "GPKG")
	if err != nil {
		panic("error creating ead output")
	}
	defer sw.Close()
	hp := hazardprovider.InitWithGrd(hazardfp, grdfp)
	defer hp.Close()
	nsp, err := structureprovider.InitStructureProvider(inventoryfp, "nsi(shape)", "GPKG")
	nsp.SetDeterministic(true)

	if err != nil {
		panic("error creating ead output")
	}
	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	frequencies := []float64{.5, .2, .1, .05, .02, .01, .005, .002, .001, .0002, .0001}
	rng := rand.New(rand.NewSource(seed))
	//get FilterStructures
	lle := initLifeLossEngine(complianceRate, rng)
	nsp.ByBbox(bbox, func(f consequences.Receptor) {
		ScopingToolProcess(f, hp, frequencies, sw, lle, rng)
	})
}
func ExpectedAnnualDamages_ResultsWriter(hazardfp string, gridfp string, inventoryfp string, sw consequences.ResultsWriter, complianceRate float64, seed int64) {

	hp := hazardprovider.InitWithGrd(hazardfp, gridfp)
	defer hp.Close()
	//@TODO handle structure file not found better
	nsp, err := structureprovider.InitStructureProvider(inventoryfp, "nsi", "GPKG")
	if err != nil {
		panic("error creating ead output")
	}
	nsp.SetDeterministic(true)
	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	frequencies := []float64{.5, .2, .1, .05, .02, .01, .005, .002, .001, .0002, .0001}
	//get FilterStructures
	rng := rand.New(rand.NewSource(seed))
	lle := initLifeLossEngine(complianceRate, rng)
	nsp.ByBbox(bbox, func(f consequences.Receptor) {
		ScopingToolProcess(f, hp, frequencies, sw, lle, rng)
	})
}
func CriticalInfrastructure_ResultsWriter(hazardfp string, gridfp string, inventoryfp string, sw consequences.ResultsWriter) {

	hp := hazardprovider.InitWithGrd(hazardfp, gridfp)
	defer hp.Close()
	//@TODO handle structure file not found better
	cisp, err := criticalinfrastructure.InitCriticalInfrastructureProvider(inventoryfp, "criticalInfrastructure", "GPKG")
	if err != nil {
		panic("error creating ci output")
	}
	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	frequencies := []float64{.5, .2, .1, .05, .02, .01, .005, .002, .001, .0002, .0001}
	//get FilterStructures
	cisp.ByBbox(bbox, func(f consequences.Receptor) {
		ScopingToolCriticalInfrastructureProcess(f, hp, frequencies, sw)
	})
}
func ExpectedAnnualDamagesGPK_WithWAVE_HDF(grdfp string, swlfp string, hmofp string, dataset string, inventoryfp string, complianceRate float64, seed int64) {
	outputPathParts := strings.Split(swlfp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}
	outfp += fmt.Sprintf("_ead_consequences_%v.gpkg", complianceRate)
	sw, err := gcrw.InitSpatialResultsWriter(outfp, "EAD_RESULTS", "GPKG") //swap to geopackage.
	if err != nil {
		panic("error creating ead output")
	}
	defer sw.Close()
	hp, err := hazardprovider.NewHdfAdcercHazardProvider(grdfp, swlfp, hmofp, dataset)
	if err != nil {
		panic("error reading hazard input")
	}
	defer hp.Close()
	nsp, err := structureprovider.InitStructureProvider(inventoryfp, "nsi(shape)", "GPKG")
	if err != nil {
		panic("error loading structure inventory")
	}
	nsp.SetDeterministic(true)
	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	//frequencies := []float64{10.0, 5.0, 2.0, 1.0, .5, .2, .1, .05, .02, .01, .005, .002, .001, .0005, .0002, .0001, .00005, .00002, .00001, .000005, .000002, .000001}
	frequencies := hp.Frequencies()
	//get FilterStructures
	rng := rand.New(rand.NewSource(seed))
	lle := initLifeLossEngine(complianceRate, rng)
	nsp.ByBbox(bbox, func(f consequences.Receptor) {
		ScopingToolProcess(f, hp, frequencies, sw, lle, rng)
	})
}
func CriticalInfrastructureGPK_WithWAVE_HDF(grdfp string, swlfp string, hmofp string, dataset string, inventoryfp string) {
	outputPathParts := strings.Split(swlfp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}
	outfp += fmt.Sprintf("_criticalinfrastructure_consequences.gpkg")
	sw, err := gcrw.InitSpatialResultsWriter(outfp, "CI_RESULTS", "GPKG") //swap to geopackage.
	if err != nil {
		panic("error creating CI output")
	}
	defer sw.Close()
	hp, err := hazardprovider.NewHdfAdcercHazardProvider(grdfp, swlfp, hmofp, dataset)
	if err != nil {
		panic("error reading hazard input")
	}
	defer hp.Close()
	cisp, err := criticalinfrastructure.InitCriticalInfrastructureProvider(inventoryfp, "criticalInfrastructure", "GPKG")
	if err != nil {
		panic("error loading structure inventory")
	}

	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	//frequencies := []float64{10.0, 5.0, 2.0, 1.0, .5, .2, .1, .05, .02, .01, .005, .002, .001, .0005, .0002, .0001, .00005, .00002, .00001, .000005, .000002, .000001}
	frequencies := hp.Frequencies()
	//get FilterStructures

	cisp.ByBbox(bbox, func(f consequences.Receptor) {
		ScopingToolCriticalInfrastructureProcess(f, hp, frequencies, sw)
	})
}
func ExpectedAnnualDamagesGPK_WithWAVE(grdfp string, swlfp string, hmo string, inventoryfp string, complianceRate float64, seed int64) {
	outputPathParts := strings.Split(swlfp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}
	outfp += "_ead_consequences.gpkg"
	sw, err := gcrw.InitSpatialResultsWriter(outfp, "EAD_RESULTS", "GPKG") //swap to geopackage.
	if err != nil {
		panic("error creating ead output")
	}
	defer sw.Close()
	hp := hazardprovider.InitWithGrdAndWave(grdfp, swlfp, hmo)
	defer hp.Close()
	nsp, err := structureprovider.InitStructureProvider(inventoryfp, "nsi", "GPKG")
	if err != nil {
		panic("error creating ead output")
	}
	nsp.SetDeterministic(true)
	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	frequencies := []float64{10.0, 5.0, 2.0, 1.0, .5, .2, .1, .05, .02, .01, .005, .002, .001, .0005, .0002, .0001, .00005, .00002, .00001, .000005, .000002, .000001}
	//get FilterStructures
	rng := rand.New(rand.NewSource(seed))
	lle := initLifeLossEngine(complianceRate, rng)
	nsp.ByBbox(bbox, func(f consequences.Receptor) {
		ScopingToolProcess(f, hp, frequencies, sw, lle, rng)
	})
}

func ExpectedAnnualDamages_OSEOutput(hazardfp string, inventoryfp string) {
	frequencies := []float64{.5, .2, .1, .05, .02, .01, .005, .002, .001, .0002, .0001}
	outputPathParts := strings.Split(hazardfp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}
	out2fp := outfp
	out3fp := outfp
	outfp += "_ose_consequences.csv"
	out2fp += "_ead_consequences.csv"
	sw, err := resultswriters.InitSummaryResultsWriterFromFile(out2fp)
	if err != nil {
		panic("error creating ead output")
	}
	defer sw.Close()
	ose_sw, err := resultswriters.InitOseResultsWriterFromFile(outfp, frequencies)
	if err != nil {
		panic("error creating ead output")
	}
	defer ose_sw.Close()
	out3fp += "_ead_consequences.gpkg"
	sw3, err := gcrw.InitSpatialResultsWriter(outfp, "EAD_RESULTS", "GPKG") //swap to geopackage.
	if err != nil {
		panic("error creating ead output")
	}
	defer sw3.Close()
	hp := hazardprovider.Init(hazardfp)
	defer hp.Close()
	nsp, err := structureprovider.InitStructureProvider(inventoryfp, "nsi", "GPKG")
	if err != nil {
		panic("error creating ead output")
	}
	nsp.SetDeterministic(true)
	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	nsp.SetDeterministic(true)
	//get FilterStructures
	nsp.ByBbox(bbox, func(f consequences.Receptor) {
		//ProvideHazard works off of a geography.Location
		ds, err2 := hp.ProvideHazards(geography.Location{X: f.Location().X, Y: f.Location().Y})
		//compute damages based on hazard being able to provide depth
		header := []string{"fd_id", "x", "y", "hazards", "damage category", "occupancy type", "structure EAD", "content EAD", "pop2amu65", "pop2amo65", "pop2pmu65", "pop2pmo65"}
		results := []interface{}{"updateme", 0.0, 0.0, ds, "dc", "ot", 0.0, 0.0, 0, 0, 0, 0}
		var ret = consequences.Result{Headers: header, Result: results}
		if err2 == nil {
			//ds is an array of hazard events
			cdams := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
			sdams := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
			//lends := len(ds)

			for i, d := range ds {
				r, err := f.Compute(d)
				if err == nil {
					//something got "wet"
					ret.Result[0] = r.Result[0]
					ret.Result[1] = r.Result[1]
					ret.Result[2] = r.Result[2]
					ret.Result[4] = r.Result[4]
					ret.Result[5] = r.Result[5]
					ret.Result[8] = r.Result[8]
					ret.Result[9] = r.Result[9]
					ret.Result[10] = r.Result[10]
					ret.Result[11] = r.Result[11]
					sdams[i] = r.Result[6].(float64)
					cdams[i] = r.Result[7].(float64)
					s, sok := f.(structures.StructureDeterministic)
					if sok {
						if s.FoundHt < d.Depth() {
							ose_sw.SetFrequencyIndex(i) //so that results get stored in the right column.
							ose_sw.Write(ret)
						}

					}

				}
			}
			//compute EAD
			cead := compute.ComputeSpecialEAD(cdams, frequencies)
			sead := compute.ComputeSpecialEAD(sdams, frequencies)
			ret.Result[6] = sead
			ret.Result[7] = cead
			if ret.Result[1] != 0.0 {
				if sead != 0 || cead != 0 {
					sw.Write(ret)
					sw3.Write(ret)
				}
			}

		}
	})
}

func ExpectedAnnualDamages_OSEOutput_CT(hazardfp string, inventoryfp string, fipscode string) {
	frequencies := []float64{.5, .2, .1, .05, .02, .01, .005, .002, .001, .0002, .0001}
	outputPathParts := strings.Split(hazardfp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}
	out2fp := outfp
	out3fp := outfp
	outfp += "_" + fipscode + "_ose_consequences.csv"
	out2fp += "_" + fipscode + "_ead_consequences.csv"
	sw, err := resultswriters.InitSummaryResultsWriterFromFile(out2fp)
	if err != nil {
		panic("error creating ead output")
	}
	defer sw.Close()
	ose_sw, err := resultswriters.InitOseResultsWriterFromFile(outfp, frequencies)
	if err != nil {
		panic("error creating ead output")
	}
	defer ose_sw.Close()
	out3fp += "_" + fipscode + "_ead_consequences.gpkg"
	sw3, err := gcrw.InitSpatialResultsWriter(out3fp, "EAD_RESULTS", "GPKG") //swap to geopackage.
	if err != nil {
		panic("error creating ead output")
	}
	defer sw3.Close()
	hp := hazardprovider.Init(hazardfp)
	defer hp.Close()
	nsp, err := structureprovider.InitStructureProvider(inventoryfp, "nsi", "GPKG")
	if err != nil {
		panic("error creating ead output")
	}
	nsp.SetDeterministic(true)
	//get FilterStructures
	nsp.ByFips(fipscode, func(f consequences.Receptor) {
		//ProvideHazard works off of a geography.Location
		ds, err2 := hp.ProvideHazards(geography.Location{X: f.Location().X, Y: f.Location().Y})
		//compute damages based on hazard being able to provide depth
		header := []string{"fd_id", "x", "y", "hazards", "damage category", "occupancy type", "structure EAD", "content EAD", "pop2amu65", "pop2amo65", "pop2pmu65", "pop2pmo65"}
		results := []interface{}{"updateme", 0.0, 0.0, ds, "dc", "ot", 0.0, 0.0, 0, 0, 0, 0}
		var ret = consequences.Result{Headers: header, Result: results}
		if err2 == nil {
			//ds is an array of hazard events
			cdams := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
			sdams := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
			//lends := len(ds)

			for i, d := range ds {
				r, err := f.Compute(d)
				if err == nil {
					//something got "wet"
					ret.Result[0] = r.Result[0]
					ret.Result[1] = r.Result[1]
					ret.Result[2] = r.Result[2]
					ret.Result[4] = r.Result[4]
					ret.Result[5] = r.Result[5]
					ret.Result[8] = r.Result[8]
					ret.Result[9] = r.Result[9]
					ret.Result[10] = r.Result[10]
					ret.Result[11] = r.Result[11]
					sdams[i] = r.Result[6].(float64)
					cdams[i] = r.Result[7].(float64)
					s, sok := f.(structures.StructureDeterministic)
					if sok {
						if s.FoundHt < d.Depth() {
							if sdams[i] != 0 || cdams[i] != 0 {
								ose_sw.SetFrequencyIndex(i) //so that results get stored in the right column.
								ose_sw.Write(ret)
							}
						}

					}

				}
			}
			//compute EAD
			cead := compute.ComputeSpecialEAD(cdams, frequencies)
			sead := compute.ComputeSpecialEAD(sdams, frequencies)
			ret.Result[6] = sead
			ret.Result[7] = cead
			if ret.Result[1] != 0.0 {
				if sead != 0 || cead != 0 {
					sw.Write(ret)
					sw3.Write(ret)
				}
			}

		}
	})
}
func ScopingToolProcess(f consequences.Receptor, hp hazardprovider.HazardProvider, frequencies []float64, sw consequences.ResultsWriter, lle lifeloss.LifeLossEngine, rng *rand.Rand) {
	//ProvideHazard works off of a geography.Location
	s, ok := f.(structures.StructureDeterministic)
	if !ok {
		return
	}
	ds, err2 := hp.ProvideHazards(geography.Location{X: f.Location().X, Y: f.Location().Y})
	//set up destination headers and results for each structure.
	header := []string{"fd_id", "x", "y", "hazards", "damcat", "occtype", "s EAD", "c EAD", "pop2amu65", "pop2amo65", "pop2pmu65", "pop2pmo65", "found_ht", "firmzone", "aall_u65", "aall_o65", "aal_tot"}
	results := []interface{}{"updateme", 0.0, 0.0, ds, "dc", "ot", 0.0, 0.0, 0, 0, 0, 0, 0.0, "", 0.0, 0.0, 0.0}
	for _, f := range frequencies {
		header = append(header, fmt.Sprintf("%2.6fS", f))
		header = append(header, fmt.Sprintf("%2.6fC", f))
		header = append(header, fmt.Sprintf("%2.6fU", f))
		header = append(header, fmt.Sprintf("%2.6fO", f))
		header = append(header, fmt.Sprintf("%2.6fH", f))

		results = append(results, 0.0)
		results = append(results, 0.0)
		results = append(results, 0.0)
		results = append(results, 0.0)
		results = append(results, "no hazard")
	}
	var ret = consequences.Result{Headers: header, Result: results}
	if err2 == nil {
		//ds is an array of hazard events
		cdams := make([]float64, len(frequencies))
		sdams := make([]float64, len(frequencies))
		u65nlls := make([]float64, len(frequencies))
		o65nlls := make([]float64, len(frequencies))
		hazardEvents := make([]hazards.CoastalEvent, len(frequencies))
		lends := len(ds)
		for i, d := range ds {
			//compute economic loss
			r, err := f.Compute(d)
			if err == nil {
				if i == lends-1 { //on the last one get the attributes.
					ret.Result[0] = r.Result[0]
					ret.Result[1] = r.Result[1]
					ret.Result[2] = r.Result[2]
					ret.Result[4] = r.Result[4]
					ret.Result[5] = r.Result[5]
					ret.Result[8] = r.Result[8]
					ret.Result[9] = r.Result[9]
					ret.Result[10] = r.Result[10]
					ret.Result[11] = r.Result[11]
					ret.Result[12] = s.FoundHt
					ret.Result[13] = s.FirmZone
				}
				sdams[i] = r.Result[6].(float64)
				cdams[i] = r.Result[7].(float64)
				hazardEvents[i] = r.Result[3].(hazards.CoastalEvent)
				ret.Result[17+(i*5)] = sdams[i]
				ret.Result[17+(i*5)+1] = cdams[i]
				//update hazard to include dv
				llevent := hazards.DepthandDVEvent{}
				llevent.SetDepth(d.Depth())
				llevent.SetDV(0.0)
				if d.Has(hazards.Velocity) {
					llevent.SetDV(d.Depth() * d.Velocity())
				} else if d.Has(hazards.DV) {
					llevent.SetDV(d.DV())
				} else if d.Has(hazards.WaveHeight) {
					//if waveheight>3 => VE zone
					llevent.SetDV(d.Depth() * 6.5)
				} else {
					switch s.FirmZone {
					case "VE", "V1-30":
						llevent.SetDV(d.Depth() * 6.5)
					}
				}
				//compute life loss
				stability, _ := lle.EvaluateStabilityCriteria(llevent, s)
				llr, err := lle.ComputeLifeLoss(llevent, s, stability)
				if err != nil {
					panic(err)
				}
				u65nll, err := llr.Fetch("ll_u65")
				u65nllint := u65nll.(int32)
				u65nlls[i] = float64(u65nllint)
				o65nll, err := llr.Fetch("ll_o65")
				o65nllint := o65nll.(int32)
				o65nlls[i] = float64(o65nllint)
				ret.Result[17+(i*5)+2] = u65nlls[i]
				ret.Result[17+(i*5)+3] = o65nlls[i]
				b, _ := json.Marshal(d)
				hazard := string(b)
				ret.Result[17+(i*5)+4] = hazard
			}
		}
		//compute EAD
		cead := compute.ComputeSpecialEAD(cdams, frequencies)
		sead := compute.ComputeSpecialEAD(sdams, frequencies)
		stringHazards := ""
		for _, he := range hazardEvents {
			b, _ := json.Marshal(he)
			stringHazards += string(b)
		}
		ret.Result[3] = stringHazards
		ret.Result[6] = sead
		ret.Result[7] = cead
		//compute average annualized life loss numbers
		u65aal := compute.ComputeSpecialEAD(u65nlls, frequencies)
		o65aal := compute.ComputeSpecialEAD(o65nlls, frequencies)
		//update life loss aalls
		ret.Result[14] = u65aal          //u65 aal
		ret.Result[15] = o65aal          //o65 aal
		ret.Result[16] = u65aal + o65aal //tot aal
		if ret.Result[1] != 0.0 {
			if sead != 0 || cead != 0 {
				sw.Write(ret)
			}
		}

	}
}
func ScopingToolCriticalInfrastructureProcess(f consequences.Receptor, hp hazardprovider.HazardProvider, frequencies []float64, sw consequences.ResultsWriter) {
	//ProvideHazard works off of a geography.Location
	_, ok := f.(criticalinfrastructure.CriticalInfrastructureFeature)
	if !ok {
		return
	}
	ds, err2 := hp.ProvideHazards(geography.Location{X: f.Location().X, Y: f.Location().Y})
	//set up destination headers and results for each structure.
	header := []string{"name", "x", "y", "hazards", "sector", "lifeline"}
	results := []interface{}{"updateme", 0.0, 0.0, ds, "s", "l"}
	for _, f := range frequencies {
		header = append(header, fmt.Sprintf("%2.6fH", f))
		results = append(results, "no hazard")
	}
	var ret = consequences.Result{Headers: header, Result: results}
	if err2 == nil {
		//ds is an array of hazard events
		hazardEvents := make([]hazards.CoastalEvent, len(frequencies))
		lends := len(ds)
		for i, d := range ds {
			if d.Depth() > 0 {
				//compute impact
				r, err := f.Compute(d)
				if err == nil {
					if i == lends-1 { //on the last one get the attributes.
						ret.Result[0] = r.Result[0]
						ret.Result[1] = r.Result[1]
						ret.Result[2] = r.Result[2]
						ret.Result[4] = r.Result[4]
						ret.Result[5] = r.Result[3]

					}

					hazardEvents[i] = r.Result[5].(hazards.CoastalEvent)
					b, _ := json.Marshal(d)
					hazard := string(b)
					ret.Result[5+(i*1)+1] = hazard
				}
			}

		}

		stringHazards := ""
		for _, he := range hazardEvents {
			b, _ := json.Marshal(he)
			stringHazards += string(b)
		}
		ret.Result[3] = stringHazards

		if ret.Result[1] != 0.0 {
			sw.Write(ret)
		}

	}
}
