package hazardprovider

import (
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/geography"
)

var HDF_TC_FILE string = "/Working/hec/go-coastal/CHS_LACS_Grid_Information.h5"
var HDF_SWL_FILE string = "/Working/hec/go-coastal/CHS_LACS_AEF_SWL_SLC0.h5"
var HDF_HM0_FILE string = "/Working/hec/go-coastal/CHS_LACS_AEF_Hm0_SLC0.h5"

func TestNewHdfAdcercHazardProvider(t *testing.T) {
	hzp, err := NewHdfAdcercHazardProvider(HDF_TC_FILE, HDF_SWL_FILE, HDF_HM0_FILE, "BE (standard)")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hzp)

	haz, err := hzp.ProvideHazards(geography.Location{X: -90.099577, Y: 29.949915})
	fmt.Println(haz)
}

func TestReadNodesHdf(t *testing.T) {
	nodes, err := ReadNodes(HDF_TC_FILE)
	if err != nil {
		t.Fatal(err)
	}
	var l int32 = int32(len(nodes))
	node := nodes[l-1] //last element in test dataset

	if l != 1239389 || node.Point.Z != 3.2825396 {
		t.Fatalf("Expected len:1239389 and ZElev:3.2825396 got %d, %.8f", l, node.Point.Z)
	}
}

func TestReadTrianglesHdf(t *testing.T) {
	triangles, err := ReadTriangles(HDF_TC_FILE)
	if err != nil {
		t.Fatal(err)
	}
	var l int32 = int32(len(triangles))
	triangle := triangles[3067211] //last element in test dataset

	if l != 2406785 || triangle.NodeA != 1563151 {
		t.Fatalf("Expected len:1239389 and NodeA:1563151 got %d, %d", l, triangle.NodeA)
	}
}

func TestReadAefHdf(t *testing.T) {
	aef, err := ReadAEF(HDF_TC_FILE)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(aef)

	if aef[0] != 0.1 || aef[15] != 10000 {
		t.Fatalf("Expected 0.1 and 10000 got %f, %f", aef[0], aef[15])
	}
}

func TestReadProb(t *testing.T) {
	probs, err := ReadProbabilities(HDF_SWL_FILE, "84% (standard)")
	if err != nil {
		t.Fatal(err)
	}
	defer probs.Close()
	rowdata := []float64{}
	probs.ReadRow(0, &rowdata)
	fmt.Println(rowdata)
}
