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
	fmt.Println(ds.LayerByIndex(0).Definition().FieldCount())
	//need to figure out how to set the featuredefinition

	//fmt.Println(ds.LayerByIndex(0).Extent(true)) //produces "Illegal Error"
	//fmt.Println(ds.LayerByIndex(0).FeatureCount(true))
	return csvHazardProvider{ds: &ds}
}
func (csv csvHazardProvider) ProvideHazard(l geography.Location) (hazards.HazardEvent, error) {

	h := hazards.CoastalEvent{}
	h.SetDepth(123.45) //update from the actual file
	h.SetSalinity(true)
	return h, nil
}
func (csv csvHazardProvider) ProvideHazardBoundary() (geography.BBox, error) {

	return geography.BBox{}, errors.New("stop asking these questions...")
}
func (csv csvHazardProvider) Close() {
	//do nothing?
	csv.ds.Destroy()
}
