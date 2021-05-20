package main

import (
	"github.com/HydrologicEngineeringCenter/go-coastal/hazardprovider"
	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/consequences"

	"github.com/USACE/go-consequences/structureprovider"
)

func main() {
	nsp := structureprovider.InitNSISP()
	f := hazardprovider.OneHundred
	hp := hazardprovider.Init("/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315_a.csv")
	hp.SetFrequency(int(f)) //pass in frequency?
	sw := consequences.InitGeoJsonResultsWriterFromFile("/workspaces/go-coastal/data/CHS_SACS_FL_Blending_PCHA_depth_SLC0_BE_v2020315_a_consequences.json")
	//fmt.Println("FIPS Code is " + "12") //for florida
	compute.StreamAbstract(hp, nsp, sw)

}
