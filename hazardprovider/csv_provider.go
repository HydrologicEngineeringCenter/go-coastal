package hazardprovider

import (
	"encoding/csv"
	"errors"
	"log"
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
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	return csvHazardProvider{csv: r}
}
func (csv csvHazardProvider) Close() {
	//do nothing?
}
func (csv csvHazardProvider) ProvideHazard(l geography.Location) (hazards.HazardEvent, error) {
	h := hazards.DepthEvent{}
	h.SetDepth(123.45)
	return h, nil
}
func (csv csvHazardProvider) ProvideHazardBoundary() (geography.BBox, error) {

	return geography.BBox{}, errors.New("stop asking these questions...")
}
