package hazardprovider

import (
	"testing"
)

func TestSHPByFips(t *testing.T) {
	hp := Init("/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315.csv") //pass in frequency?
	hp.Close()
}
