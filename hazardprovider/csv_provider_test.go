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

	f := OneHundred
	hp := Init("/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315_a.csv", int(f)) //pass in frequency?
	defer hp.Close()
	sw := consequences.InitGeoJsonResultsWriterFromFile("/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315_a_consequences.json")
	defer sw.Close()
	//fmt.Println("FIPS Code is " + "12") //for florida

	nsp := structureprovider.InitGPK("/workspaces/go-coastal/data/nsiv2_12.gpkg", "nsi")

	compute.StreamAbstract(hp, nsp, sw)
}
