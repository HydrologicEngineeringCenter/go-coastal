package compute

import (
	"testing"

	"github.com/HydrologicEngineeringCenter/go-coastal/hazardprovider"
)

func Test_Event(t *testing.T) {
	f := hazardprovider.OneHundred
	hp := "/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315.csv" //pass in frequency
	sp := "/workspaces/go-coastal/data/nsiv2_12.gpkg"
	Event(hp, sp, int(f))
}
