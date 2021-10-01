package hazardprovider

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/HydrologicEngineeringCenter/go-coastal/geometry"
	"github.com/furstenheim/ConcaveHull"
	"github.com/tidwall/rtree"
)

type CoastalFrequency int

const (
	Two          CoastalFrequency = 3
	Five         CoastalFrequency = 4
	Ten          CoastalFrequency = 5
	Twenty       CoastalFrequency = 6
	Fifty        CoastalFrequency = 7
	OneHundred   CoastalFrequency = 8
	TwoHundred   CoastalFrequency = 9
	FiveHundred  CoastalFrequency = 10
	OneThousand  CoastalFrequency = 11
	FiveThousand CoastalFrequency = 12
	TenThousand  CoastalFrequency = 13
)

var coastalFrequencies = []CoastalFrequency{
	Two,
	Five,
	Ten,
	Twenty,
	Fifty,
	OneHundred,
	TwoHundred,
	FiveHundred,
	OneThousand,
	FiveThousand,
	TenThousand,
}

func (c CoastalFrequency) String() string {
	return [...]string{"Two", "Five", "Ten", "Twenty", "Fifty", "OneHundred", "TwoHundred", "FiveHundred", "OneThousand", "FiveThousand", "TenThousand"}[c-3]
}

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

func process_TIN(fp string) (*geometry.Tin, error) {
	f, err := os.Open(fp)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	nodata := 0.0
	xidx := 0
	yidx := 1
	firstrow := true
	//we dont know how big the file will be, so we have to make a guess.
	dimSize := 0
	points := make([]geometry.PointZ, dimSize)
	count := 0
	ps := make([]float64, dimSize)
	for scanner.Scan() {
		lines := strings.Split(scanner.Text(), ",")
		//check if first value is negative to determine lat/lon
		if firstrow {
			testval, err := strconv.ParseFloat(lines[0], 64)
			if err != nil {
				panic(err)
			}
			if testval < 0 {
				//0 is negative, must be lon
			} else {
				yidx = 0
				xidx = 1
			}
			firstrow = false
		}
		xval, err := strconv.ParseFloat(lines[xidx], 64)
		if err != nil {
			panic(err)
		}
		ps = append(ps, xval)
		yval, err := strconv.ParseFloat(lines[yidx], 64)
		if err != nil {
			panic(err)
		}
		ps = append(ps, yval)
		terrain, err := strconv.ParseFloat(lines[2], 64)
		if err != nil {
			panic(err)
		}
		//loop over z values
		zvals := make([]float64, len(coastalFrequencies))
		for i, zconst := range coastalFrequencies {
			zval, err := strconv.ParseFloat(lines[int(zconst)], 64) //need to read all values and load into an array now.
			if err != nil {
				panic(err)
			}
			if terrain < 0 {
				zval = zval + terrain //(minus a negative to get value above sea level...)
			}
			if zval == 0 {
				zval = nodata
			} else {
				zval *= 3.28084 //convert from meters to feet
			}
			zvals[i] = zval
		}
		p := geometry.Point{X: xval, Y: yval}
		points = append(points, geometry.PointZ{Point: &p, Z: zvals})
		count++
	}
	fmt.Printf("read %v lines from %v\n", count, fp)
	fmt.Println("Creating Concave Hull...")
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
	t, err := geometry.CreateTin(points, nodata, p)
	return t, err
	/*
		pts := t.ConvexHull
		s := "{\"type\": \"FeatureCollection\",\"features\": [{\"type\": \"Feature\",\"geometry\": {\"type\": \"LineString\",\"coordinates\": ["
		for _, p := range pts {
			s += "[" + fmt.Sprintf("%g, %g",p.X, p.Y) + "],"
		}

		s = strings.TrimRight(s, ",")
		s += "]},\"properties\": {\"prop1\": 0.0}}]}"

		fmt.Println(s)

		s := strings.TrimRight(fp,".csv")
		t.Json(s + ".json")
	*/
}
func processGrdAndCSV(grdfp string, csvfp string) (*geometry.Tin, error) {
	grdf, err := os.Open(grdfp)
	if err != nil {
		panic(err)
	}
	defer grdf.Close()

	scanner := bufio.NewScanner(grdf)
	//we dont know how big the file will be, so we have to make a guess.
	scanner.Scan() // burn the header
	scanner.Scan() //count of triangles and points
	row2 := scanner.Text()
	vals := strings.Split(row2, "  ") //not sure this will always work correctly
	dimNSize, _ := strconv.ParseInt(vals[1], 10, 64)
	dimTSize, _ := strconv.ParseInt(vals[0], 10, 64)
	nodes := make(map[int64]geometry.PointZ)
	triangles := make(map[int64]geometry.Triangle)
	var triangleCounter int64
	var pointCounter int64
	triangleCounter = 0
	pointCounter = 0
	var minx, miny, maxx, maxy float64
	minx = 180
	miny = 180
	maxx = -180
	maxy = -180
	var tr rtree.RTree
	ps := make([]float64, 0)
	for scanner.Scan() {
		if pointCounter < dimNSize { //points come first.
			if pointCounter == 0 {
				fmt.Println("reading nodes")
			}
			pointCounter += 1
			line := strings.Split(scanner.Text(), " ") //is there a way to group spaces?
			//nodeid|X|Y|Z
			var nodeid int64
			nodeid = 0
			xval := 0.0
			yval := 0.0
			zval := 0.0
			valcount := 0
			for _, v := range line {
				if v != "" {
					valcount += 1
					switch valcount {
					case 1:
						nodeid, _ = strconv.ParseInt(v, 10, 32)
					case 2:
						xval, _ = strconv.ParseFloat(v, 64)
					case 3:
						yval, _ = strconv.ParseFloat(v, 64)
					case 4:
						zval, _ = strconv.ParseFloat(v, 64)
					}
				}
			}
			//replace with depth values?
			zvals := []float64{zval}
			nodes[nodeid] = geometry.PointZ{Point: &geometry.Point{X: xval, Y: yval}, Z: zvals}
			if xval > maxx {
				maxx = xval
			} else {
				if xval < minx {
					minx = xval
				}
			}
			if yval > maxy {
				maxy = yval
			} else {
				if yval < miny {
					miny = yval
				}
			}
			ps = append(ps, xval)
			ps = append(ps, yval)

		} else if triangleCounter < dimTSize { //triangles come second.
			if triangleCounter == 0 {
				fmt.Println("reading triangles")
			}
			triangleCounter += 1
			line := strings.Split(scanner.Text(), " ") //is there a way to group spaces?
			//triangleid|vertcount|a|b|c
			var triangleid, aidx, bidx, cidx int64
			triangleid = 0
			aidx = 0.0
			bidx = 0.0
			cidx = 0.0
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
			a, aok := nodes[int64(aidx)]
			if !aok {
				panic(fmt.Sprintf("Not a ok! TriangleCounter is %v, a index is %v, and total triangles is %v", triangleCounter, aidx, dimTSize))
			}
			b, bok := nodes[int64(bidx)]
			if !bok {
				panic(fmt.Sprintf("Not b ok! TriangleCounter is %v, b index is %v, and total triangles is %v", triangleCounter, bidx, dimTSize))
			}
			c, cok := nodes[int64(cidx)]
			if !cok {
				panic(fmt.Sprintf("Not c ok! TriangleCounter is %v, c index is %v, and total triangles is %v", triangleCounter, cidx, dimTSize))
			}
			t := geometry.CreateTriangle(&a, &b, &c)
			e := t.Extent()
			triangles[triangleid] = t
			//add to tree?
			tr.Insert(e.LowerLeft.ToXY(), e.UpperRight.ToXY(), &t)
		} else {
			//must be in the hull def?
			break
		}

	}
	fmt.Println("Finished reading, computing Hull")
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
}
