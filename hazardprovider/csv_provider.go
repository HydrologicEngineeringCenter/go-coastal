package hazardprovider

import (
	"github.com/HydrologicEngineeringCenter/go-coastal/geometry"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

type csvHazardProvider struct {
	//csv *csv.Reader
	ds *geometry.Tin
}

//Init creates and produces an unexported csvHazardProvider
func Init(fp string, zidx int) csvHazardProvider {
	// Open the file
	t, err := process_TIN(fp, zidx)
	if err != nil {
		panic(err)
	}
	return csvHazardProvider{ds: t}
}
func (csv csvHazardProvider) ProvideHazard(l geography.Location) (hazards.HazardEvent, error) {
	h := hazards.CoastalEvent{}
	v, err := csv.ds.ComputeValue(l.X, l.Y)
	if err != nil {
		h.SetDepth(-9999.0)
		return h, err
	}
	h.SetDepth(v)
	h.SetSalinity(true)
	return h, nil
}
func (csv csvHazardProvider) ProvideHazardBoundary() (geography.BBox, error) {
	bbox := make([]float64, 4)
	bbox[0] = csv.ds.MinX //upper left x
	bbox[1] = csv.ds.MaxY //upper left y
	bbox[2] = csv.ds.MaxX //lower right x
	bbox[3] = csv.ds.MinY //lower right y

	return geography.BBox{Bbox: bbox}, nil
}
func (csv csvHazardProvider) Close() {
	//do nothing?

}
