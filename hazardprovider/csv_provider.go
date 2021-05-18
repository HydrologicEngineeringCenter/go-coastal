package hazardprovider

import (
	"fmt"
	"time"

	"github.com/HydrologicEngineeringCenter/go-coastal/geometry"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
)

type csvHazardProvider struct {
	//csv *csv.Reader
	ds                       *geometry.Tin
	queryCount               int64
	actualComputedStructures int64
	computeStart             time.Time
}

//Init creates and produces an unexported csvHazardProvider
func Init(fp string, zidx int) *csvHazardProvider {
	// Open the file
	t, err := process_TIN(fp, zidx)
	if err != nil {
		panic(err)
	}
	c := time.Now()
	return &csvHazardProvider{ds: t, computeStart: c}
}
func (csv *csvHazardProvider) ProvideHazard(l geography.Location) (hazards.HazardEvent, error) {
	h := hazards.CoastalEvent{}
	csv.queryCount++
	//check if point is in the hull polygon.
	p := geometry.Point{X: l.X, Y: l.Y}
	if csv.queryCount%100000 == 0 {
		n := time.Since(csv.computeStart)
		fmt.Print("Compute Time: ")
		fmt.Println(n)
		fmt.Println(fmt.Sprintf("Processed %v structures, with %v valid depths", csv.queryCount, csv.actualComputedStructures))
	}
	if csv.ds.Hull.Contains(p) {
		v, err := csv.ds.ComputeValue(l.X, l.Y)
		if err != nil {
			h.SetDepth(-9999.0)
			return h, err
		}
		h.SetDepth(v)
		h.SetSalinity(true)
		csv.actualComputedStructures++
		return h, nil
	}
	notIn := hazardproviders.NoHazardFoundError{Input: "Point Not In Polygon"}
	h.SetDepth(-9999.0)
	return h, notIn
}
func (csv csvHazardProvider) ProvideHazardBoundary() (geography.BBox, error) {
	bbox := make([]float64, 4)
	bbox[0] = csv.ds.MinX //upper left x
	bbox[1] = csv.ds.MaxY //upper left y
	bbox[2] = csv.ds.MaxX //lower right x
	bbox[3] = csv.ds.MinY //lower right y
	return geography.BBox{Bbox: bbox}, nil
}
func (csv *csvHazardProvider) Close() {
	//do nothing?
	n := time.Since(csv.computeStart)
	fmt.Print("Compute Complete")
	fmt.Print("Compute Time was: ")
	fmt.Println(n)
	fmt.Println(fmt.Sprintf("Processed %v structures, with %v valid depths", csv.queryCount, csv.actualComputedStructures))

}
