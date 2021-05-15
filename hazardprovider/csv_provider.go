package hazardprovider

import (
	"errors"

	"github.com/HydrologicEngineeringCenter/go-coastal/geometry"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

type csvHazardProvider struct {
	//csv *csv.Reader
	ds *geometry.Tin
}

//Init creates and produces an unexported csvHazardProvider
func Init(fp string) csvHazardProvider {
	// Open the file
	t, err := process_TIN(fp, int(OneHundred))
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
	return geography.BBox{}, errors.New("stop asking these questions...")
}
func (csv csvHazardProvider) Close() {
	//do nothing?

}
