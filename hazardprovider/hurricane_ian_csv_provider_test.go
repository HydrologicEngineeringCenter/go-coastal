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

func Test_Hurricane_IAN(t *testing.T) {
	scenario := "low_surge_wave"
	//load hazard data
	hfp := fmt.Sprintf("/workspaces/go-coastal/data/ian/%v_30.csv", scenario)
	hp, err := InitIanCSVFileProvider(hfp)
	defer hp.Close()
	if err != nil {
		panic(err)
	}
	//load structure data
	sfp := "/workspaces/go-coastal/data/ian/30BuildingsPointData.shp"
	sp, err := structureprovider.InitStructureProvider(sfp, "30BuildingsPointData", "ESRI Shapefile")

	if err != nil {
		panic(err)
	}

	//choose a results writer.
	rfp := fmt.Sprintf("/workspaces/go-coastal/data/ian/30BuildingsPointData_%v_results.json", scenario)
	rw, err := resultswriters.InitSpatialResultsWriter(rfp, "results", "GeoJSON")
	if err != nil {
		panic(err)
	}
	defer rw.Close()
	//compute.StreamAbstract(hp, sp, rw) //doesnt seem to write out wave by default - fixing that manually.

	fmt.Println("Getting bbox")
	bbox, err := hp.HazardBoundary()
	if err != nil {
		log.Panicf("Unable to get the raster bounding box: %s", err)
	}
	fmt.Println(bbox.ToString())
	sp.ByBbox(bbox, func(f consequences.Receptor) {
		//ProvideHazard works off of a geography.Location
		he, err2 := hp.Hazard(geography.Location{X: f.Location().X, Y: f.Location().Y})
		//compute damages based on hazard being able to provide depth
		if err2 == nil {
			r, err3 := f.Compute(he)
			r.Headers = append(r.Headers, "wave_h_ft")
			//jsonstring := string(b)
			r.Result = append(r.Result, he.WaveHeight())
			if err3 == nil {
				rw.Write(r)
			}
		}
	})

}
