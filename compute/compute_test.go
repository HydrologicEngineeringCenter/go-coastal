package compute

import (
	"strings"
	"testing"

	"github.com/HydrologicEngineeringCenter/go-coastal/hazardprovider"
	gcrw "github.com/USACE/go-consequences/resultswriters"
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

//@TODO:export this as a c function
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
