package hazardprovider

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

type csvHazardProvider struct {
	csv *csv.Reader
}

//Init creates and produces an unexported csvHazardProvider
func Init(fp string) csvHazardProvider {
	// Open the file
	csvfile, err := os.Open("input.csv")
	if err != nil {
		fmt.Println("Couldn't open the csv file") //, err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	return csvHazardProvider{csv: r}
}
func (csv csvHazardProvider) Close() {
	//do nothing?
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
