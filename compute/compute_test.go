package compute

import (
	"testing"

	"github.com/HydrologicEngineeringCenter/go-coastal/hazardprovider"
)

func Test_Event(t *testing.T) {
	f := hazardprovider.OneHundred
	hp := "/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315_a.csv"
	sp := "/workspaces/go-coastal/data/nsiv2_12.gpkg"
	Event(hp, sp, int(f)) //pass in frequency
}
func Test_EAD(t *testing.T) {
	hp := "/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315_a.csv"
	sp := "/workspaces/go-coastal/data/nsiv2_12.gpkg"
	ExpectedAnnualDamages(hp, sp)
}
func Test_EADGpk(t *testing.T) {
	hp := "/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315.csv"
	sp := "/workspaces/go-coastal/data/nsiv2_12.gpkg"
	ExpectedAnnualDamagesGPK(hp, sp)
}

/*
func Test_EADGpk_parallel(t *testing.T) {
	hp := "/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315.csv"
	sp := "/workspaces/go-coastal/data/nsiv2_12.gpkg"
	ExpectedAnnualDamagesGPK_FIPS(hp, sp, "12")
}
*/
