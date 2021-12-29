package main

import (
	"fmt"
	"os"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 2 {
		fmt.Println("Expected two arguments, the filepath to the csv input and the file path to the geopackage input")
	} else {
		hfp := argsWithoutProg[0]
		sfp := argsWithoutProg[1]
		fmt.Println(fmt.Sprintf("Computing EAD for %v using an iventory at path %v", hfp, sfp))
		//compute.ExpectedAnnualDamagesGPK(hfp, sfp)
	}
}
