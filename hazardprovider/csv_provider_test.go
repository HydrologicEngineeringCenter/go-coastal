package hazardprovider

import (
	"fmt"
	"testing"
)

func TestOpenCSV_WithCSVProvider(t *testing.T) {
	hp := Init("/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315.csv")
	hp.Close()
}
func TestConvertCSV_WithCSVConverter(t *testing.T) {
	f := OneHundred
	processCSV2Tif("/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315.csv", "/workspaces/go-coastal/data/FL_SLC0_BE_"+fmt.Sprint(f)+".tif", int(f))
}
