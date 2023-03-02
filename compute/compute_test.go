package compute

import (
	"strings"
	"testing"

	"github.com/HydrologicEngineeringCenter/go-coastal/hazardprovider"
	"github.com/HydrologicEngineeringCenter/go-coastal/resultswriters"
	"github.com/USACE/go-consequences/hazardproviders"
	gcrw "github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structureprovider"
)

func Test_Event(t *testing.T) {
	f := hazardprovider.OneHundred
	hp := "/Working/hec/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315_a.csv"
	sp := "/Working/hec/go-coastal/data/nsiv2_12.gpkg"
	Event(hp, sp, int(f)) //pass in frequency
}
func Test_EAD(t *testing.T) {
	hp := "/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315_a.csv"
	sp := "/workspaces/go-coastal/data/nsiv2_12.gpkg"
	ExpectedAnnualDamages(hp, sp)
}
func Test_EAD_resultsWriter(t *testing.T) {
	hp := "/Working/hec/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315.csv"
	sp := "/Working/hec/go-coastal/data/nsiv2_12.gpkg"
	outputPathParts := strings.Split(hp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}
	outfp += "_ead_consequences.shp"
	sw, err := gcrw.InitShpResultsWriter(outfp, "EADResults") //swap to geopackage.
	if err != nil {
		panic("error creating ead output")
	}
	defer sw.Close()
	ExpectedAnnualDamages_ResultsWriter(hp, sp, sw)
}

// @TODO:export this as a c function
func Test_EADGpk_WithWaves(t *testing.T) {
	fp := "/Working/hec/go-coastal/data/NAC2014_R01_ClosedRivers.grd"
	swlp := "/Working/hec/go-coastal/data/NACS_Nantucket_PCHA_SLC0_SWL_BE_v20210722.csv"
	hmop := "/Working/hec/go-coastal/data/NACS_Nantucket_PCHA_SLC0_Hm0_BE_v20210722.csv"
	sp := "/Working/hec/go-coastal/data/nsi.gpkg"
	ExpectedAnnualDamagesGPK_WithWAVE(fp, swlp, hmop, sp)
}

func Test_EADGpk_WithWavesHdf5(t *testing.T) {

	fp := "/Working/hec/go-coastal/CHS_LACS_Grid_Information.h5"
	swlp := "/Working/hec/go-coastal/CHS_LACS_AEF_SWL_SLC0.h5"
	hmop := "/Working/hec/go-coastal/CHS_LACS_AEF_Hm0_SLC0.h5"
	sp := "/Working/hec/go-coastal/data/nsi.gpkg"
	ExpectedAnnualDamagesGPK_WithWAVE_HDF(fp, swlp, hmop, "BE (standard)", sp)
}
func Test_EAD_OSE(t *testing.T) {
	hp := "/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315.csv"
	sp := "/workspaces/go-coastal/data/nsiv2_12.gpkg"
	ExpectedAnnualDamages_OSEOutput(hp, sp)
}

func Test_EAD_OSE_CT(t *testing.T) {
	hp := "/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315.csv"
	sp := "/workspaces/go-coastal/data/nsiv2_12.gpkg"
	ExpectedAnnualDamages_OSEOutput_CT(hp, sp, "12086008900")
}

func Test_WoodHole_Event(t *testing.T) {

	wsefp := "/workspaces/go-coastal/data/woodhole/20 Year 2030_wgs84.tif"
	wavefp := "/workspaces/go-coastal/data/woodhole/20 Year 2030 waves_wgs84.tif"
	spfp := "/workspaces/go-coastal/data/nsi.gpkg"
	rwfp := "/workspaces/go-coastal/data/woodhole/wh_2030_20Y.gpkg"
	hp := hazardprovider.InitWoodHoleGroupTif(wsefp, wavefp)
	defer hp.Close()
	sp, err := structureprovider.InitGPK(spfp, "nsi")
	defer hp.Close()
	if err != nil {
		panic("error creating inventory provider")
	}
	rw, err := gcrw.InitGpkResultsWriter(rwfp, "results")
	if err != nil {
		panic("error creating results writer")
	}
	defer rw.Close()
	if err != nil {
		panic("error creating results writer")
	}
	WoodHoleEvent(hp, sp, rw)
}

func Test_WoodHole_EAD(t *testing.T) {

	wsefp := []string{
		"/workspaces/go-coastal/data/woodhole/20 Year 2030_wgs84.tif",
		"/workspaces/go-coastal/data/woodhole/50 Year 2030_wgs84.tif",
		"/workspaces/go-coastal/data/woodhole/100 Year 2030_wgs84.tif",
		"/workspaces/go-coastal/data/woodhole/200 Year 2030_wgs84.tif",
		"/workspaces/go-coastal/data/woodhole/500 Year 2030_wgs84.tif",
		"/workspaces/go-coastal/data/woodhole/1000 Year 2030_wgs84.tif",
	}
	wavefp := []string{
		"/workspaces/go-coastal/data/woodhole/20 Year 2030 waves_wgs84.tif",
		"/workspaces/go-coastal/data/woodhole/50 Year 2030 waves_wgs84.tif",
		"/workspaces/go-coastal/data/woodhole/100 Year 2030 waves_wgs84.tif",
		"/workspaces/go-coastal/data/woodhole/200 Year 2030 waves_wgs84.tif",
		"/workspaces/go-coastal/data/woodhole/500 Year 2030 waves_wgs84.tif",
		"/workspaces/go-coastal/data/woodhole/1000 Year 2030 waves_wgs84.tif",
	}
	frequencies := []float64{
		1.0 / 20.0,
		1.0 / 50.0,
		1.0 / 100.0,
		1.0 / 200.0,
		1.0 / 500.0,
		1.0 / 1000.0,
	}
	spfp := "/workspaces/go-coastal/data/nsi.gpkg"
	rwfp := "/workspaces/go-coastal/data/woodhole/wh_EAD.gpkg"
	hps := make([]hazardproviders.HazardProvider, len(wsefp))
	for idx, wse := range wsefp {
		hp := hazardprovider.InitWoodHoleGroupTif(wse, wavefp[idx])
		hps[idx] = hp
	}
	sp, err := structureprovider.InitGPK(spfp, "nsi")
	if err != nil {
		panic("error creating inventory provider")
	}
	sp.SetDeterministic(true)
	//rw, err := gcrw.InitGpkResultsWriter(rwfp, "EAD results")
	rw, err := resultswriters.InitwoodHoleResultsWriterFromFile(rwfp, frequencies)

	if err != nil {
		panic("error creating results writer")
	}
	defer rw.Close()
	if err != nil {
		panic("error creating results writer")
	}
	WoodHoleDeterministicEAD(hps, frequencies, sp, rw)
}
