package compute

import (
	"strings"

	"github.com/HydrologicEngineeringCenter/go-coastal/hazardprovider"
	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/structureprovider"
)

func Event(hazardfp string, inventoryfp string, frequency int) {
	hp := hazardprovider.Init(hazardfp, frequency) //pass in frequency?
	defer hp.Close()
	outputPathParts := strings.Split(hazardfp, ".")
	outfp := outputPathParts[0]
	for i := 1; i < len(outputPathParts)-1; i++ {
		outfp += "." + outputPathParts[i]
	}
	outfp += "_consequences.json"
	sw := consequences.InitGeoJsonResultsWriterFromFile(outfp)
	defer sw.Close()
	nsp := structureprovider.InitGPK(inventoryfp, "nsi")

	compute.StreamAbstract(hp, nsp, sw)
}
