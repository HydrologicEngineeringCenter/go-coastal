package hazardprovider

import (
	"reflect"
)

/*
	- open grid dataset
	-

*/

var CHL_HDFGRID_TC string = "Triangular Connectivity"
var CHL_HDFGRID_NI string = "Nodal Information"
var CHL_HDFGRID_AEF string = "AEF"

//@QUESTION: is a 32bit int large enough for big models?
//@QUESTION: is LACS a study, hence we use a study node?
type AcNode struct {
	Point      AcPoint
	Index      int32 //nodes oridinal position in hdf dataset.  Should be equyal to "model node id-1"
	AdcircNode int32
}

type AcTriangle struct {
	ElementId int32
	NodeA     int32
	NodeB     int32
	NodeC     int32
	UrbLat    float64
	UrbLon    float64
	LrbLat    float64
	LrbLon    float64
}

type AcPoint struct {
	X float64
	Y float64
	Z float64
}

func ReadTriangles(hdfFilePath string) (map[int32]AcTriangle, error) {
	triangularConnsOptions := HdfReadOptions{
		Dtype:           reflect.Float64,
		IncrementalRead: true,
		IncrementSize:   1000, //read 1000 rows at a time
	}
	tc, err := NewHdfDataset(hdfFilePath, CHL_HDFGRID_TC, triangularConnsOptions)
	if err != nil {
		return nil, err
	}
	defer tc.Close()
	triangles := make(map[int32]AcTriangle)

	row := []float64{}
	for i := 0; i < tc.Rows(); i++ {
		tc.ReadRow(i, &row)

		node := AcTriangle{
			ElementId: int32(row[0]),
			NodeA:     int32(row[1]),
			NodeB:     int32(row[2]),
			NodeC:     int32(row[3]),
			UrbLat:    row[4],
			UrbLon:    row[5],
			LrbLat:    row[6],
			LrbLon:    row[7],
		}
		triangles[node.ElementId] = node
	}
	return triangles, nil
}

func ReadNodes(hdfFilePath string) (map[int32]AcNode, error) {

	nodalOptions := HdfReadOptions{
		Dtype:           reflect.Float64,
		IncrementalRead: true,
		IncrementSize:   1000, //read 1000 rows at a time
	}
	nodeData, err := NewHdfDataset(hdfFilePath, CHL_HDFGRID_NI, nodalOptions)
	if err != nil {
		return nil, err
	}
	defer nodeData.Close()

	nodes := make(map[int32]AcNode)

	row := []float64{}
	for i := 0; i < nodeData.Rows(); i++ {
		nodeData.ReadRow(i, &row)

		node := AcNode{
			Point: AcPoint{
				X: row[3],
				Y: row[2],
				Z: row[4],
			},
			Index:      int32(i),
			AdcircNode: int32(row[1]),
		}
		nodes[node.Index] = node
	}
	return nodes, nil
}

func ReadAEF(hdfFilePath string) ([]float64, error) {
	options := HdfReadOptions{
		Dtype: reflect.Float64,
	}
	dataset, err := NewHdfDataset(hdfFilePath, CHL_HDFGRID_AEF, options)
	if err != nil {
		return nil, err
	}
	defer dataset.Close()

	cols := dataset.Cols()
	data := make([]float64, cols)
	dataset.ReadInto(&data)
	return data, nil
}

func ReadProbabilities(hdfFilePath string, hdfDataPath string) (*HdfDataset, error) {
	options := HdfReadOptions{
		Dtype: reflect.Float64,
	}
	return NewHdfDataset(hdfFilePath, hdfDataPath, options)
}
