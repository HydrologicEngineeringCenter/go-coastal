package hazardprovider

import (
	"testing"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"

	"github.com/USACE/go-consequences/structureprovider"
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
func Test_Compute(t *testing.T) {
	nsp := structureprovider.InitNSISP()
	f := OneHundred
	hp := Init("/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315_a.csv", int(f)) //pass in frequency?
	sw := consequences.InitGeoJsonResultsWriterFromFile("/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315_a_consequences.json")
	//fmt.Println("FIPS Code is " + "12") //for florida
	compute.StreamAbstract(hp, nsp, sw)
}
