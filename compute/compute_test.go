package compute

import (
	"fmt"
	"strings"
	"testing"

	"github.com/HydrologicEngineeringCenter/go-coastal/hazardprovider"
	gcrw "github.com/USACE/go-consequences/resultswriters"
)

func Test_Event(t *testing.T) {
	f := hazardprovider.Fifty
	hp := "/workspaces/go-coastal/data/SACS/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315.csv"
	sp := "/workspaces/go-coastal/data/nsi.gpkg"
	Event(hp, sp, int(f), f.String()) //pass in frequency
}
func Test_Event_Grid_CSV(t *testing.T) {
	//f := hazardprovider.TenThousand
	cellsize := .0001
	hp := "/workspaces/go-coastal/data/SACS/PR/CHS-SACS_PR_PCHA_Nodal_Inundation_Depth_SLC2_BE_vOct2023.csv"
	gp := "/workspaces/go-coastal/data/SACS/sacs_prusvi_base_g001.grd"
	outputPathParts := strings.Split(hp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}

	//hazp := hazardprovider.InitWithGrd(hp, gp)
	//offset to zero based position.
	//defer hazp.Close()
	for _, f := range hazardprovider.CoastalFrequencies {
		//hazp.SelectFrequency(int(f) - int(hazardprovider.Two))
		Event_Grid(hp, gp, int(f), f.String(), cellsize) //pass in frequency
	}

}
func Test_Event_Grid_HDF(t *testing.T) {
	cellsize := .0001
	hmop := "/workspaces/go-coastal/data/NACS/SLC1/CHS-NA_CC_SimBslc1_Post1RT_Nodes_Hm0_AEF.h5"
	swlp := "/workspaces/go-coastal/data/NACS/SLC1/CHS-NA_CC_SimBslc1_Post1RT_Nodes_SWL_AEF.h5"
	gp := "/workspaces/go-coastal/data/NACS/CHS-NA_Spat_Sim0_Post0_Nodes_ADCIRC_Locations.h5"
	outputPathParts := strings.Split(swlp, ".")
	outfp := outputPathParts[0]

	hazp, err := hazardprovider.NewHdfAdcercHazardProvider(gp, swlp, hmop, "Best Estimate AEF")
	//offset to zero based position.
	if err != nil {
		panic(err)
	}
	defer hazp.Close()
	for i, f := range hazp.Frequencies() {
		if f != 0.00002 {
			if f < 0.11 {
				outputfilepath := fmt.Sprintf("%v_%2.6f.tif", outfp, f)
				hazp.SelectFrequency(i)
				Event_Grid_new(outputfilepath, hazp, cellsize) //pass in frequency
			}

		}

	}

}
func Test_EAD(t *testing.T) {
	hp := "/workspaces/go-coastal/data/SACS/NC/CHS-SACS_NC_PCHA_Nodal_Inundation_Depth_SLC2_BE_vOct2023.csv"
	gp := "/workspaces/go-coastal/data/SACS/sacs_sa_base_g001.grd"
	sp := "/workspaces/go-coastal/data/nsi_2022.gpkg"
	complianceRate := 0.75
	seed := 1234
	ExpectedAnnualDamages(hp, gp, sp, complianceRate, int64(seed))
}
func Test_EAD_resultsWriter(t *testing.T) {
	hp := "/workspaces/go-coastal/data/SACS/AL/CHS-SACS_AL_PCHA_Nodal_Inundation_Depth_SLC0_BE_v2023.csv"
	gp := "/workspaces/go-coastal/data/SACS/sacs_gm_base_g001.grd"
	sp := "/workspaces/go-coastal/data/nsi_2022.gpkg"
	outputPathParts := strings.Split(hp, ".")
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
	complianceRate := 0.75
	seed := 1234
	ExpectedAnnualDamages_ResultsWriter(hp, gp, sp, sw, complianceRate, int64(seed))
}

// @TODO:export this as a c function
func Test_EADGpk_WithWaves(t *testing.T) {
	fp := "/Working/hec/go-coastal/data/NAC2014_R01_ClosedRivers.grd"
	swlp := "/Working/hec/go-coastal/data/NACS_Nantucket_PCHA_SLC0_SWL_BE_v20210722.csv"
	hmop := "/Working/hec/go-coastal/data/NACS_Nantucket_PCHA_SLC0_Hm0_BE_v20210722.csv"
	sp := "/Working/hec/go-coastal/data/nsi.gpkg"
	complianceRate := 0.75
	seed := 1234
	ExpectedAnnualDamagesGPK_WithWAVE(fp, swlp, hmop, sp, complianceRate, int64(seed))
}

func Test_EADGpk_WithWavesHdf5(t *testing.T) {

	fp := "/workspaces/go-coastal/data/CHS_LACS_Grid_Information.h5"
	swlp := "/workspaces/go-coastal/data/CHS_LACS_AEF_SWL_SLC0.h5"
	hmop := "/workspaces/go-coastal/data/CHS_LACS_AEF_Hm0_SLC0.h5"
	sp := "/workspaces/go-coastal/data/nsi.gpkg"
	complianceRate := 0.75
	seed := 1234
	ExpectedAnnualDamagesGPK_WithWAVE_HDF(fp, swlp, hmop, "BE (standard)", sp, complianceRate, int64(seed))
}
func Test_EADGpk_WithWavesHdf5_LACS(t *testing.T) {

	fp := "/workspaces/go-coastal/data/LACS/CHS-LA_Spat_Sim0_Post0_Nodes_ADCIRC_Locations.h5"
	swlp := "/workspaces/go-coastal/data/LACS/detq/CHS-LA_TS_SimBrfc2_Post1RT_Nodes_SWL_AEF.h5"
	hmop := "/workspaces/go-coastal/data/LACS/detq/CHS-LA_TS_SimBrfc2_Post0_Nodes_Hm0_AEF.h5"
	/*


				hmop := "/workspaces/go-coastal/data/TXCS/SLC0/CHS-TX_TS_SimB_Post0_Nodes_Hm0_AEF.h5"
				swlp := "/workspaces/go-coastal/data/TXCS/SLC0/CHS-TX_TS_SimB_Post1RT_Nodes_SWL_AEF.h5"
				fp := "/workspaces/go-coastal/data/TXCS/CHS-TX_Spat_Sim0_Post0_Nodes_ADCIRC_Locations.h5"

		hmop := "/workspaces/go-coastal/data/NACS/SLC1/CHS-NA_CC_SimBslc1_Post1RT_Nodes_Hm0_AEF.h5"
		swlp := "/workspaces/go-coastal/data/NACS/SLC1/CHS-NA_CC_SimBslc1_Post1RT_Nodes_SWL_AEF.h5"
		fp := "/workspaces/go-coastal/data/NACS/CHS-NA_Spat_Sim0_Post0_Nodes_ADCIRC_Locations.h5"
	*/
	sp := "/workspaces/go-coastal/data/nsi_2022.gpkg"
	complianceRate := 0.75
	seed := 1234
	ExpectedAnnualDamagesGPK_WithWAVE_HDF(fp, swlp, hmop, "Best Estimate AEF", sp, complianceRate, int64(seed))
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

/*
func Test_WoodHole_Event(t *testing.T) {

	wsefp := "/workspaces/go-coastal/data/woodhole/20 Year 2030_wgs84.tif"
	wavefp := "/workspaces/go-coastal/data/woodhole/20 Year 2030 waves_wgs84.tif"
	spfp := "/workspaces/go-coastal/data/nsi.gpkg"
	rwfp := "/workspaces/go-coastal/data/woodhole/wh_2030_20Y.gpkg"
	hp := hazardprovider.InitWoodHoleGroupTif(wsefp, wavefp)
	defer hp.Close()
	sp, err := structureprovider.InitStructureProvider(spfp, "nsi","GPKG")
	defer hp.Close()
	if err != nil {
		panic("error creating inventory provider")
	}
	rw, err := gcrw.InitSpatialResultsWriter(rwfp, "results","GPKG")
	if err != nil {
		panic("error creating results writer")
	}
	defer rw.Close()
	if err != nil {
		panic("error creating results writer")
	}
	WoodHoleEvent(hp, sp, rw)
}
func Test_writeSettingsFile(t *testing.T) {
	ds2030 := WoodHoleFrequencyDataset{
		Year: 2030,
		WaterSurfaceGridPaths: []string{
			"/workspaces/go-coastal/data/woodhole/20 Year 2030_wgs84.tif",
			"/workspaces/go-coastal/data/woodhole/50 Year 2030_wgs84.tif",
			"/workspaces/go-coastal/data/woodhole/100 Year 2030_wgs84.tif",
			"/workspaces/go-coastal/data/woodhole/200 Year 2030_wgs84.tif",
			"/workspaces/go-coastal/data/woodhole/500 Year 2030_wgs84.tif",
			"/workspaces/go-coastal/data/woodhole/1000 Year 2030_wgs84.tif",
		},
		WavePaths: []string{
			"/workspaces/go-coastal/data/woodhole/20 Year 2030 waves_wgs84.tif",
			"/workspaces/go-coastal/data/woodhole/50 Year 2030 waves_wgs84.tif",
			"/workspaces/go-coastal/data/woodhole/100 Year 2030 waves_wgs84.tif",
			"/workspaces/go-coastal/data/woodhole/200 Year 2030 waves_wgs84.tif",
			"/workspaces/go-coastal/data/woodhole/500 Year 2030 waves_wgs84.tif",
			"/workspaces/go-coastal/data/woodhole/1000 Year 2030 waves_wgs84.tif",
		},
		Frequencies: []float64{
			1.0 / 20.0,
			1.0 / 50.0,
			1.0 / 100.0,
			1.0 / 200.0,
			1.0 / 500.0,
			1.0 / 1000.0,
		},
	}
	simSettings := WoodHoleSimulationSettings{
		DataSets:        []WoodHoleFrequencyDataset{ds2030},
		BaseYear:        2025,
		DiscountRate:    0.025,
		InventoryPath:   "/workspaces/go-coastal/data/nsi.gpkg",
		OutputDirectory: "/workspaces/go-coastal/data/woodhole/results/",
	}
	bytes, err := json.Marshal(simSettings)
	if err != nil {
		t.Fail()
	}
	ioutil.WriteFile("/workspaces/go-coastal/data/woodhole/settings.json", bytes, 0600)
}
func Test_WoodHole_EEAD(t *testing.T) {
	bytes, err := ioutil.ReadFile("/workspaces/go-coastal/data/woodhole/settings.json")
	if err != nil {
		t.Fail()
	}
	whss := WoodHoleSimulationSettings{}
	err = json.Unmarshal(bytes, &whss)
	if err != nil {
		t.Fail()
	}
	err = WoodHoleMultiYearDeterministicEEAD(whss)
	if err != nil {
		t.Fail()
	}
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
	rate := .025
	numYearsInFuture := 0
	rw, err := resultswriters.InitwoodHoleResultsWriterFromFile(rwfp, frequencies, CreateDiscountFactor(rate, numYearsInFuture), 2030, nil)

	if err != nil {
		panic("error creating results writer")
	}
	defer rw.Close()
	if err != nil {
		panic("error creating results writer")
	}
	WoodHoleDeterministicEAD(hps, frequencies, sp, rw)
}
*/
