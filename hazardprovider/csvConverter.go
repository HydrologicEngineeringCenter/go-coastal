package hazardprovider

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/HydrologicEngineeringCenter/go-coastal/geometry"
	"github.com/furstenheim/ConcaveHull"
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
