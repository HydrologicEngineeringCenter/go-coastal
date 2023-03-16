package compute

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/HydrologicEngineeringCenter/go-coastal/hazardprovider"
	"github.com/HydrologicEngineeringCenter/go-coastal/resultswriters"
	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
	gcrw "github.com/USACE/go-consequences/resultswriters"
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
	sw, err := gcrw.InitGeoJsonResultsWriterFromFile(outfp)
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
func ExpectedAnnualDamages_ResultsWriter(hazardfp string, inventoryfp string, sw consequences.ResultsWriter) {

	hp := hazardprovider.Init(hazardfp)
	defer hp.Close()
	//@TODO handle structure file not found better
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
 - initialize hazard provider dataset for spatial querying
    - has a single bounding box/concave hull (e.g. boundary)
 - open nsi gpkg
   - for all points in bounding bbox
     - compute hazard



func (ahp *HdfAdcercHazardProvider) Close() {
	defer ahp.nodeReader.Close()
	defer ahp.elementReader.Close()
	defer ahp.probabilityData.Close()
}

type HdfAdcercHazardProvider struct {
	nodeReader      HdfDataset
	elementReader   HdfDataset
	probabilites    []float64
	probabilityData HdfDataset
}



optionally:
 - allow hazard provider to have one or more boundaries
 - open nsi gpkg
   - for each boundary
      - for all points in bounding box
	    - compute hazard


type HazardProvider interface{
	HazardBoundary(asBbox bool) Polygon //boundary for the entire region covered by this hazard provider
	BoundaryElements() int //number of individual boundary elements
	NextElement() Polygon or nil//
	ProvideHazards(geography.Location) []hazards.HazardEvent
	Close()
}

type HdfAdcercHazardProvider struct{
	nodeReader HdfDataset
	elementReader Hdf5Dataset
	probabilites []float64
	probabilityData HdfDataset
}

//implements all Hazard Provider interfaces
//only change to code is to add an outer loop around nps.Bbox...for example

for {
	elemBound,err:=NextElement()
	if err!=nil{
		return err
	}
	if elemBound==nil{
		break
	}
	nsp.ByPoly(polygon,func(f consequences.Receptor){

	})
}

*/

func ExpectedAnnualDamagesGPK_WithWAVE_HDF(grdfp string, swlfp string, hmofp string, dataset string, inventoryfp string) {
	outputPathParts := strings.Split(swlfp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}
	outfp += "_ead_consequences.gpkg"
	sw, err := gcrw.InitGpkResultsWriter(outfp, "EAD_RESULTS") //swap to geopackage.
	if err != nil {
		panic("error creating ead output")
	}
	defer sw.Close()
	hp, err := hazardprovider.NewHdfAdcercHazardProvider(grdfp, swlfp, hmofp, dataset)
	//defer hp.Close()
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
	//frequencies := []float64{10.0, 5.0, 2.0, 1.0, .5, .2, .1, .05, .02, .01, .005, .002, .001, .0005, .0002, .0001, .00005, .00002, .00001, .000005, .000002, .000001}
	frequencies := hp.Frequencies()
	//get FilterStructures

	nsp.ByBbox(bbox, func(f consequences.Receptor) {
		//ProvideHazard works off of a geography.Location
		ds, err2 := hp.ProvideHazards(geography.Location{X: f.Location().X, Y: f.Location().Y})
		//compute damages based on hazard being able to provide depth
		header := []string{"fd_id", "x", "y", "hazards", "damcat", "occtype", "s EAD", "c EAD", "pop2amu65", "pop2amo65", "pop2pmu65", "pop2pmo65"}
		results := []interface{}{"updateme", 0.0, 0.0, ds, "dc", "ot", 0.0, 0.0, 0, 0, 0, 0}
		var ret = consequences.Result{Headers: header, Result: results}
		if err2 == nil {
			//ds is an array of hazard events
			cdams := make([]float64, len(frequencies))
			sdams := make([]float64, len(frequencies))
			hazardEvents := make([]hazards.CoastalEvent, len(frequencies))
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
					hazardEvents[i] = r.Result[3].(hazards.CoastalEvent)
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
			if ret.Result[1] != 0.0 {
				if sead != 0 || cead != 0 {
					sw.Write(ret)
				}
			}

		}
	})
}

func ExpectedAnnualDamagesGPK_WithWAVE(grdfp string, swlfp string, hmo string, inventoryfp string) {
	outputPathParts := strings.Split(swlfp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}
	outfp += "_ead_consequences.gpkg"
	sw, err := gcrw.InitGpkResultsWriter(outfp, "EAD_RESULTS") //swap to geopackage.
	if err != nil {
		panic("error creating ead output")
	}
	defer sw.Close()
	hp := hazardprovider.InitWithGrd(grdfp, swlfp, hmo)
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
	frequencies := []float64{10.0, 5.0, 2.0, 1.0, .5, .2, .1, .05, .02, .01, .005, .002, .001, .0005, .0002, .0001, .00005, .00002, .00001, .000005, .000002, .000001}
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
			cdams := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
			sdams := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
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
	sw3, err := gcrw.InitGpkResultsWriter(outfp, "EAD_RESULTS") //swap to geopackage.
	if err != nil {
		panic("error creating ead output")
	}
	defer sw3.Close()
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
	sw3, err := gcrw.InitGpkResultsWriter(out3fp, "EAD_RESULTS") //swap to geopackage.
	if err != nil {
		panic("error creating ead output")
	}
	defer sw3.Close()
	hp := hazardprovider.Init(hazardfp)
	defer hp.Close()
	nsp, err := structureprovider.InitGPK(inventoryfp, "nsi")
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
