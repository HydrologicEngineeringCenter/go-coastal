package hazardprovider

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/HydrologicEngineeringCenter/go-coastal/geometry"
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
	"github.com/furstenheim/ConcaveHull"
	"github.com/tidwall/rtree"
)

type chrpsHazardProvider struct {
	//csv *csv.Reader
	ds                       *geometry.Tin
	queryCount               int64
	actualComputedStructures int64
	computeStart             time.Time
	gridFile                 string
	eventFile                string
}

func (csv *chrpsHazardProvider) ProvideHazard(l geography.Location) (hazards.HazardEvent, error) {
	//h := hazards.CoastalEvent{}
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
		v, err := csv.ds.ComputeValues(l.X, l.Y)
		if err != nil {
			return nil, err
		}
		csv.actualComputedStructures++
		return v[0], nil
	}
	notIn := hazardproviders.NoHazardFoundError{Input: "Point Not In Polygon"}
	return nil, notIn
}

//notIn := hazardproviders.NoHazardFoundError{Input: "Point Not In Polygon"}
//h.SetDepth(-9999.0)
//return h, notIn
//}

// implement
func (csv *chrpsHazardProvider) ProvideHazards(l geography.Location) ([]hazards.HazardEvent, error) {
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
		v, err := csv.ds.ComputeValues(l.X, l.Y)
		if err != nil {
			return nil, err
		}
		csv.actualComputedStructures++
		return v, nil
	}
	notIn := hazardproviders.NoHazardFoundError{Input: "Point Not In Polygon"}
	return nil, notIn
}

// implement
func (csv chrpsHazardProvider) ProvideHazardBoundary() (geography.BBox, error) {
	bbox := make([]float64, 4)
	bbox[0] = csv.ds.MinX //upper left x
	bbox[1] = csv.ds.MaxY //upper left y
	bbox[2] = csv.ds.MaxX //lower right x
	bbox[3] = csv.ds.MinY //lower right y
	return geography.BBox{Bbox: bbox}, nil
}

// implement
func (csv *chrpsHazardProvider) Close() {
	//do nothing?
	n := time.Since(csv.computeStart)
	fmt.Print("Compute Complete")
	fmt.Print("Compute Time was: ")
	fmt.Println(n)
	fmt.Println(fmt.Sprintf("Processed %v structures, with %v valid depths", csv.queryCount, csv.actualComputedStructures))

}

type Event struct {
	GPM      string      `json:"GPM"`
	Response [][]float64 `json:"resp"`
}

func InitCHRPS(gridFile string, eventFile string) (chrpsHazardProvider, error) {
	chp := chrpsHazardProvider{
		gridFile:                 gridFile,
		eventFile:                eventFile,
		actualComputedStructures: 0,
		computeStart:             time.Now(),
	}
	ds, err := ReadCHRPS_GRD_Event(gridFile, eventFile)
	if err != nil {
		panic(err)
	}
	chp.ds = ds
	return chp, nil
}
func ReadCHRPS_GRD_Event(grdfp string, eventFile string) (*geometry.Tin, error) {
	grdf, err := os.Open(grdfp)
	if err != nil {
		panic(err)
	}
	defer grdf.Close()

	scanner := bufio.NewScanner(grdf)
	//we dont know how big the file will be, so we have to make a guess.
	scanner.Scan() // burn the header
	scanner.Scan() //count of triangles and points
	row2 := strings.Trim(scanner.Text(), " ")
	vals := strings.Split(row2, "  ") //not sure this will always work correctly
	dimNSize, _ := strconv.ParseInt(vals[1], 10, 64)
	dimTSize, _ := strconv.ParseInt(vals[0], 10, 64) //test.
	nodes := make(map[int32]geometry.PointZZ)
	triangles := make(map[int64]geometry.TriangleZZ)
	var triangleCounter int64
	var pointCounter int64
	triangleCounter = 0
	pointCounter = 0
	var minx, miny, maxx, maxy float64
	minx = 180
	miny = 180
	maxx = -180
	maxy = -180
	loadData := true
	for scanner.Scan() {
		if pointCounter < dimNSize { //points come first.
			if pointCounter == 0 {
				fmt.Println("reading nodes")
			}
			pointCounter += 1
			cleanLine := strings.Trim(scanner.Text(), " ") //leading and trailing
			line := strings.Split(cleanLine, " ")          //is there a way to group spaces?
			//nodeid|X|Y|Z
			var nodeid int32
			nodeid = 0
			xval := 0.0
			yval := 0.0
			zval := 0.0 //terrain
			valcount := 0
			for _, v := range line {
				if v != "" {
					valcount += 1
					switch valcount {
					case 1:
						tmpInt, err := strconv.ParseInt(v, 10, 32)
						if err != nil {
							panic(err)
						}
						nodeid = int32(tmpInt)
					case 2:
						xval, _ = strconv.ParseFloat(v, 64)
					case 3:
						yval, _ = strconv.ParseFloat(v, 64)
					case 4:
						zval, _ = strconv.ParseFloat(v, 64) //terrain
					}
				}
			}
			nodes[nodeid] = geometry.PointZZ{Point: &geometry.Point{X: xval, Y: yval}, ZElev: zval}
		} else if triangleCounter < dimTSize { //triangles come second.
			if loadData {
				//read the event file to load in data into the nodes.
				eventf, err := os.Open(eventFile)
				if err != nil {
					panic(err)
				}
				defer eventf.Close()
				eventbytes, err := ioutil.ReadAll(eventf)
				if err != nil {
					panic(err)
				}
				e := Event{}
				err = json.Unmarshal(eventbytes, &e)
				if err != nil {
					panic(err)
				}

				for _, data := range e.Response {
					//parse each row.
					nodeid := int32(data[0])
					if err != nil {
						panic(err)
					}

					node := nodes[nodeid]
					zhmo := make([]float64, 1)
					zswl := make([]float64, 1)
					zhmo[0] = 0.0
					zswl[0] = data[3]
					node.ZHm0 = zhmo
					node.ZSwl = zswl
					nodes[nodeid] = node
				}
				loadData = false
			}

			if triangleCounter == 0 {
				fmt.Println("reading triangles")
			}
			triangleCounter += 1
			line := strings.Split(scanner.Text(), " ") //is there a way to group spaces?
			//triangleid|vertcount|a|b|c
			var triangleid, aidx, bidx, cidx int64
			triangleid = 0
			aidx = 0
			bidx = 0
			cidx = 0
			valcount := 0
			for _, v := range line {
				if v != "" {
					valcount += 1
					switch valcount {
					case 1:
						triangleid, _ = strconv.ParseInt(v, 10, 32)
					case 2:
						//skip count of vertices
						break
					case 3:
						aidx, _ = strconv.ParseInt(v, 10, 64)
					case 4:
						bidx, _ = strconv.ParseInt(v, 10, 64)
					case 5:
						cidx, _ = strconv.ParseInt(v, 10, 64)
					}
				}
			}
			a, aok := nodes[int32(aidx)]
			if !aok {
				panic(fmt.Sprintf("Not a ok! TriangleCounter is %v, a index is %v, and total triangles is %v", triangleCounter, aidx, dimTSize))
			}
			b, bok := nodes[int32(bidx)]
			if !bok {
				panic(fmt.Sprintf("Not b ok! TriangleCounter is %v, b index is %v, and total triangles is %v", triangleCounter, bidx, dimTSize))
			}
			c, cok := nodes[int32(cidx)]
			if !cok {
				panic(fmt.Sprintf("Not c ok! TriangleCounter is %v, c index is %v, and total triangles is %v", triangleCounter, cidx, dimTSize))
			}
			t := geometry.CreateTriangleZZ(&a, &b, &c)
			triangles[triangleid] = t

		} else {
			//must be in the hull def?
			break
		}
	}

	//should probably reduce the space to remove triangles.
	var tr rtree.RTree
	ps := make([]float64, 0)
	kept := 0
	culled := 0
	for _, triangle := range triangles {
		if triangle.HasData() {
			//add to tree?
			e := triangle.Extent()
			tr.Insert(e.LowerLeft.ToXY(), e.UpperRight.ToXY(), triangle)
			if e.Max()[0] > maxx {
				maxx = e.Max()[0]
			} else {
				if e.Min()[0] < minx {
					minx = e.Min()[0]
				}
			}
			if e.Max()[1] > maxy {
				maxy = e.Max()[1]
			} else {
				if e.Min()[1] < miny {
					miny = e.Min()[1]
				}
			}
			ps = append(ps, triangle.Points()...)
			kept += 1
		} else {
			culled += 1
		}
	}
	fmt.Println(fmt.Sprintf("kept %v, culled %v", kept, culled))
	fmt.Println("Finished reading, computing Hull")
	if len(ps) > 2 {
		flathull := ConcaveHull.Compute(ConcaveHull.FlatPoints(ps))
		ptcount := len(flathull) / 2
		hull := make([]geometry.Point, ptcount+1)
		index := 0
		for i := 0; i < len(flathull); i += 2 {
			hull[index] = geometry.Point{X: flathull[i], Y: flathull[i+1]}
			index++
		}
		hull[index] = geometry.Point{X: flathull[0], Y: flathull[1]}
		p := geometry.CreatePolygon(hull)
		return &geometry.Tin{MaxX: maxx, MinX: minx, MaxY: maxy, MinY: miny, Tree: tr, Hull: p}, err
	} else {
		fmt.Println("No Points, no Hull")
		return &geometry.Tin{MaxX: maxx, MinX: minx, MaxY: maxy, MinY: miny, Tree: tr}, err
	}

}
