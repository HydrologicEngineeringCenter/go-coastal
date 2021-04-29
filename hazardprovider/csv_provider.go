package hazardprovider

import (
	"errors"
	"fmt"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
	"github.com/dewberry/gdal"
)

type csvHazardProvider struct {
	//csv *csv.Reader
	ds *gdal.DataSource
}

//Init creates and produces an unexported csvHazardProvider
func Init(fp string) csvHazardProvider {
	// Open the file
	ds := gdal.OpenDataSource(fp, int(gdal.ReadOnly))

	fmt.Println(ds.Driver().Name())
	fmt.Println(ds.Name())
	fmt.Println(ds.LayerCount())
	fmt.Println(ds.LayerByIndex(0).Extent(true)) //produces "Illegal Error"
	fmt.Println(ds.LayerByIndex(0).FeatureCount(true))
	return csvHazardProvider{ds: &ds}
}
func (csv csvHazardProvider) Close() {
	//do nothing?
	csv.ds.Destroy()
}
func (csv csvHazardProvider) ProvideHazard(l geography.Location) (hazards.HazardEvent, error) {
	/*
		This is where the hard part lives...
		//raw data is in meters - unknown projection atm.
		//There is no header, first row is real data.
		//delimiter is comma
		//lon,lat is not dependable, negative values should indicate lon...
		//lon,lat,elevation(meters),.5,.2,.1,.05,.02,.01,.005,.002,.001,.0002,.0001 (all depths in meter, still water level)
		//some xy locations have all zeros for depths.
		//xy location represents nodes on a triangular irregular mesh.
		//BE represents Best Estimate, CL90 90% exceedacnce value, SLC0 1996, SLC1 mid future condition, SLC2 high future condition
		//data organized by state.

	*/
	h := hazards.CoastalEvent{}
	h.SetDepth(123.45) //update from the actual file
	h.SetSalinity(true)
	return h, nil
}
func (csv csvHazardProvider) ProvideHazardBoundary() (geography.BBox, error) {

	return geography.BBox{}, errors.New("stop asking these questions...")
}
