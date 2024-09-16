package hazardprovider

import (
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structureprovider"
)

func Test_Hurricane_IAN(t *testing.T) {
	scenario := "low_surge_wave"
	//load hazard data
	hfp := fmt.Sprintf("/workspaces/go-coastal/data/ian/%v.csv", scenario)
	hp, err := InitIanCSVFileProvider(hfp)
	defer hp.Close()
	if err != nil {
		panic(err)
	}
	//load structure data
	sfp := "/workspaces/go-coastal/data/ian/PearlStBuildingsPointData.shp"
	sp, err := structureprovider.InitStructureProvider(sfp, "PearlStBuildingsPointData", "ESRI Shapefile")

	if err != nil {
		panic(err)
	}

	//choose a results writer.
	rfp := fmt.Sprintf("/workspaces/go-coastal/data/ian/PearlStBuildingsPointData_%v_results.json", scenario)
	rw, err := resultswriters.InitSpatialResultsWriter(rfp, "results", "GeoJSON")
	if err != nil {
		panic(err)
	}
	defer rw.Close()
	compute.StreamAbstract(hp, sp, rw)

}
