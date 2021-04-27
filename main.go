package main

import (
	"fmt"
	"strconv"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/structureprovider"
)

func main() {
	nsp := structureprovider.InitNSISP()

	fmt.Println("Hello World.")
	fmt.Println("FIPS Code is " + "12")
	var index int64 = 0
	nsp.ByFips("12", func(f consequences.Receptor) {
		index++
	})
	fmt.Println("Found " + strconv.FormatInt(index, 10) + " structures.")
}
