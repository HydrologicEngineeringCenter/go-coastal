package hazardprovider

import (
	"fmt"
	"strings"
	"testing"
)

func TestOpenCSV_WithCSVProvider(t *testing.T) {
	hp := Init("/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315.csv")
	hp.Close()
}
func Test_triangulation(t *testing.T) {
	process_TIN("/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315_a.csv")
}

func Test_ConcaveHull(t *testing.T) {
	f := OneHundred
	fp := "/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315.csv"
	hp := Init(fp)
	hp.SetFrequency(int(f) - int(Two)) //offset to zero based index
	s := strings.TrimRight(fp, ".csv")
	hp.ds.Hull.ToGeoJson(s + "_concavehull.json")
}
func Test_ConcaveHull_GRD(t *testing.T) {
	fp := "/workspaces/go-coastal/data/NAC2014_R01_ClosedRivers.grd"
	hp := InitWithGrd(fp, "")
	hp.ds.Hull.ToGeoJson("/workspaces/go-coastal/data/NAC2014_R01_ClosedRivers.json")
	fmt.Println(hp)
}
