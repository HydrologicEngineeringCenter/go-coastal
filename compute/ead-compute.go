package compute

import (
	"fmt"
	"log"
	"strings"

	"github.com/HydrologicEngineeringCenter/go-coastal/hazardprovider"
	"github.com/HydrologicEngineeringCenter/go-coastal/resultswriters"
	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/structureprovider"
	"github.com/USACE/go-consequences/structures"
)

func ExpectedAnnualDamages(hazardfp string, inventoryfp string) {
	outputPathParts := strings.Split(hazardfp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}
	outfp += "_ead_consequences.json"
	sw, err := consequences.InitGeoJsonResultsWriterFromFile(outfp)
	if err != nil {
		panic("error creating ead output")
	}
	defer sw.Close()
	hp := hazardprovider.Init(hazardfp)
	defer hp.Close()
	nsp, err := structureprovider.InitGPK(inventoryfp, "nsi")
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
			lends := len(ds)
			for i, d := range ds {
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
					}
					sdams[i] = r.Result[6].(float64)
					cdams[i] = r.Result[7].(float64)
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
				}
			}

		}
	})
}
func ExpectedAnnualDamagesGPK(hazardfp string, inventoryfp string) {
	outputPathParts := strings.Split(hazardfp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}
	outfp += "_ead_consequences.gpkg"
	sw, err := consequences.InitGpkResultsWriter(outfp, "EAD_RESULTS") //swap to geopackage.
	if err != nil {
		panic("error creating ead output")
	}
	defer sw.Close()
	hp := hazardprovider.Init(hazardfp)
	defer hp.Close()
	nsp, err := structureprovider.InitGPK(inventoryfp, "nsi")
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
			lends := len(ds)
			for i, d := range ds {
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
					}
					sdams[i] = r.Result[6].(float64)
					cdams[i] = r.Result[7].(float64)
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
				}
			}

		}
	})
}

/*
WIP
func ExpectedAnnualDamagesGPK_FIPS(hazardfp string, inventoryfp string, fips string) {
	outputPathParts := strings.Split(hazardfp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}
	hp := hazardprovider.Init(hazardfp)
	defer hp.Close()
	nsp := structureprovider.InitGPK(inventoryfp, "nsi")
	nsp.SetDeterministic(true)

	m := census.StateToCountyFipsMap()
	counties := m[fips]
	frequencies := []float64{.5, .2, .1, .05, .02, .01, .005, .002, .001, .0002, .0001}
	fmt.Println("Computing State By County")
	maxcons := 50
	loops := len(counties) / maxcons
	if loops > 0 {
		for i := 0; i < loops; i++ {
			//
			endpoint := (i + 1) * maxcons
			if endpoint > len(counties) {
				endpoint = len(counties)
			}
			tmpcounties := counties[i*maxcons : endpoint]
			waitcount := endpoint - (i * maxcons)
			wg := sync.WaitGroup{}
			wg.Add(waitcount)
			for _, ccccc := range tmpcounties {
				go func(cfips string) {
					fmt.Println("Computing " + cfips)
					defer wg.Done()
					outfpccccc := outfp + "_" + cfips + "_ead_consequences.gpkg"
					sw := consequences.InitGpkResultsWriter(outfpccccc, "EAD_RESULTS") //swap to geopackage.
					defer sw.Close()
					//get FilterStructures
					nsp.ByFips(cfips, func(f consequences.Receptor) {
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
							lends := len(ds)
							for i, d := range ds {
								if d.Has(hazards.Depth) {

									if d.Depth() > 0.0 && d.Depth() < 9999.0 {
										r := f.Compute(d)
										if i == lends-1 {
											ret.Result[0] = r.Result[0]
											ret.Result[1] = r.Result[1]
											ret.Result[2] = r.Result[2]
											ret.Result[4] = r.Result[4]
											ret.Result[5] = r.Result[5]
											ret.Result[8] = r.Result[8]
											ret.Result[9] = r.Result[9]
											ret.Result[10] = r.Result[10]
											ret.Result[11] = r.Result[11]
										}
										sdams[i] = r.Result[6].(float64)
										cdams[i] = r.Result[7].(float64)
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
								}
							}
						}
					})
				}(ccccc)
			}
			wg.Wait() //does this work with the last bit?
		}
	} else {
		//no loops?
		fmt.Println("loops was zero")
	}

	fmt.Println("Compute Complete")
}
*/
func ExpectedAnnualDamages_OSEOutput(hazardfp string, inventoryfp string) {
	frequencies := []float64{.5, .2, .1, .05, .02, .01, .005, .002, .001, .0002, .0001}
	outputPathParts := strings.Split(hazardfp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}
	out2fp := outfp
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
	hp := hazardprovider.Init(hazardfp)
	defer hp.Close()
	nsp, err := structureprovider.InitGPK(inventoryfp, "nsi")
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
				}
			}

		}
	})
}
