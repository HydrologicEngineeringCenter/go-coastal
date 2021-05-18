package hazardprovider

import (
	"strings"
	"testing"
)

func TestOpenCSV_WithCSVProvider(t *testing.T) {
	f := OneHundred
	hp := Init("/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315.csv", int(f))
	hp.Close()
}
func Test_triangulation(t *testing.T) {
	f := OneHundred
	process_TIN("/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315_a.csv", int(f))
}

func Test_ConcaveHull(t *testing.T) {
	f := OneHundred
	fp := "/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315.csv"
	hp := Init(fp, int(f))
	s := strings.TrimRight(fp, ".csv")
	hp.ds.Hull.ToGeoJson(s + "_concavehull.json")
}
