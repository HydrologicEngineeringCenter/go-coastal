package hazardprovider

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dewberry/gdal"
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

func (c CoastalFrequency) String() string {
	return [...]string{"Two", "Five", "Ten", "Twenty", "Fifty", "OneHundred", "TwoHundred", "FiveHundred", "OneThousand", "FiveThousand", "TenThousand"}[c-3]
}
func processCSV2Tif(infile string, outfile string, zidx int) {
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
	//read the csv file line by line and create an array of x values, y values, and z values (specify frequency.)
	f, err := os.Open(infile)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	var minx, miny, maxx, maxy float64
	minx = 180
	miny = 180
	maxx = -180
	maxy = -180
	nodata := -9999.0
	xidx := 0
	yidx := 1
	yRes := .01
	xRes := .01
	firstrow := true
	//we dont know how big the file will be, so we have to make a guess.
	dimSize := 0
	xvals, yvals, wse := make([]float64, dimSize), make([]float64, dimSize), make([]float64, dimSize)
	count := 0
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
		if maxx < xval {
			maxx = xval
		}
		if minx > xval {
			minx = xval
		}
		xvals = append(xvals, xval)
		yval, err := strconv.ParseFloat(lines[yidx], 64)
		if err != nil {
			panic(err)
		}
		if maxy < yval {
			maxy = yval
		}
		if miny > yval {
			miny = yval
		}
		yvals = append(yvals, yval)
		zval, err := strconv.ParseFloat(lines[zidx], 64)
		if err != nil {
			panic(err)
		}
		if zval == 0 {
			zval = nodata
		}
		//convert from meters to feet?
		wse = append(wse, zval)
		count++
	}
	fmt.Printf("read %v lines\n", count)
	nX := uint(math.Round(math.Abs(maxx-minx) / xRes))
	nY := uint(math.Round(math.Abs(maxy-miny) / yRes))
	//create regular grid with inverse distance
	grid, err := gdal.GridCreate(
		gdal.GA_InverseDistancetoAPower,
		gdal.GridInverseDistanceToAPowerOptions{NoDataValue: nodata},
		xvals, yvals, wse,
		minx, maxx, miny, maxy,
		nX, nY,
		ProgressReport,
		nil,
	)
	if err != nil {
		panic(err)
	}
	//convert to tif file
	sr := gdal.CreateSpatialReference("")
	sr.FromEPSG(4326)
	crsWkt, err := sr.ToWKT()
	if err != nil {
		panic(err)
	}
	errwse := writeTif2(outfile, crsWkt, int(nX), int(nY), minx, miny, xRes, yRes, nodata, grid)
	if errwse != nil {
		panic(errwse)
	}
	//clip to hull
}
func ProgressReport(complete float64, message string, data interface{}) int {
	d := time.Now()
	fmt.Printf("Percent Complete: %f at %s\n", complete, d.Format("3:04:05PM"))
	return gdal.DummyProgress(complete, message, data)
}
func writeTif2(outTifName, crsWKT string, xSize, ySize int, xMin, yMin, xRes, yRes, noDataVal float64, data []float64) error {
	fmt.Printf("Loading driver\n")
	driver, err := gdal.GetDriverByName("GTiff")
	if err != nil {
		return err
	}

	dataset := driver.Create("MEM:::", xSize, ySize, 1, gdal.Float64, nil)
	defer dataset.Close()

	dataset.SetProjection(crsWKT)
	dataset.SetGeoTransform([6]float64{xMin, xRes, 0, yMin, 0, yRes})
	raster := dataset.RasterBand(1)
	raster.SetNoDataValue(noDataVal)
	fmt.Println("Writing to raster band")
	raster.IO(gdal.Write, 0, 0, xSize, ySize, data, xSize, ySize, 0, 0)
	opts := []string{"-t_srs", crsWKT, "-of", "GTiff"}

	outds, err := gdal.Warp(outTifName, []gdal.Dataset{dataset}, opts)
	defer outds.Close()

	fmt.Println("Finished writing", outTifName)
	return nil
}
