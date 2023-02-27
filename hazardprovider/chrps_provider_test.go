package hazardprovider

import (
	"fmt"
	"log"
	"testing"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structureprovider"
)

func Test_CHRPS_GRD_EVENT(t *testing.T) {
	fp := "/workspaces/go-coastal/data/cpra_2023updates_v14a_chk.grd"
	hp, err := InitCHRPS(fp, "/workspaces/go-coastal/data/2008_GUSTAV_Adv_20_PredNodes.txt")
	if err != nil {
		panic(err)
	}
	hp.ds.Hull.ToGeoJson("/workspaces/go-coastal/data/gustav_culled.json")
	fmt.Println(hp)
}
func Test_CHRPS_Compute(t *testing.T) {
	root := "/workspaces/go-coastal/data/2008_GUSTAV"
	fp := "/workspaces/go-coastal/data/cpra_2023updates_v14a_chk.grd"
	hp, err := InitCHRPS(fp, "/workspaces/go-coastal/data/2008_GUSTAV_Adv_20_PredNodes.txt")
	if err != nil {
		panic(err)
	}
	nsp, err := structureprovider.InitGPK("/workspaces/go-coastal/data/nsi.gpkg", "nsi")
	if err != nil {
		panic(err)
	}
	w, _ := resultswriters.InitGeoJsonResultsWriterFromFile(root + "_consequences.json")
	defer w.Close()
	fmt.Println("Getting bbox")
	bbox, err := hp.ProvideHazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	//get FilterStructures
	nsp.ByBbox(bbox, func(f consequences.Receptor) {
		//ProvideHazard works off of a geography.Location
		d, err2 := hp.ProvideHazard(geography.Location{X: f.Location().X, Y: f.Location().Y})
		//compute damages based on hazard being able to provide depth
		if err2 == nil {
			r, err := f.Compute(d)
			if err == nil {
				w.Write(r)
			}
		}
	})
	fmt.Println(hp)
}