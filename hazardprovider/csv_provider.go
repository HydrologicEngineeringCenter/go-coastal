package hazardprovider

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
	"github.com/dewberry/gdal"
)

type csvHazardProvider struct {
	//csv *csv.Reader
	ds *gdal.DataSource
}

//Init creates and produces an unexported csvHazardProvider
func Init(fp string) csvHazardProvider {
	// Open the file
	ds := gdal.OpenDataSource(fp, int(gdal.ReadOnly))

	fmt.Println(ds.Driver().Name())
	fmt.Println(ds.Name())
	fmt.Println(ds.LayerCount())
	fmt.Println(ds.LayerByIndex(0).Definition().FieldCount())
	//need to figure out how to set the featuredefinition

	//fmt.Println(ds.LayerByIndex(0).Extent(true)) //produces "Illegal Error"
	//fmt.Println(ds.LayerByIndex(0).FeatureCount(true))
	return csvHazardProvider{ds: &ds}
}
func processCSV2Tif(file string, zidx int) {
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
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	var minx, miny, maxx, maxy float64
	nodata := -9999.0
	xidx := 0
	yidx := 1
	//assuming meters?
	yRes := 10.0
	xRes := 10.0
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
	nX := uint(math.Round(math.Abs(maxx-minx) / xRes))
	nY := uint(math.Round(math.Abs(maxy-miny) / yRes))
	//create regular grid with inverse distance
	grid, err := gdal.GridCreate(
		gdal.GA_InverseDistancetoAPower,
		gdal.GridInverseDistanceToAPowerOptions{NoDataValue: nodata},
		xvals, yvals, wse,
		minx, maxx, miny, maxy,
		nX, nY,
		gdal.DummyProgress,
		nil,
	)
	fmt.Println(err)
	//convert to tif file
	//crsWkt := `PROJCS["USA_Contiguous_Albers_Equal_Area_Conic_USGS_version",GEOGCS["GCS_North_American_1983",DATUM["D_North_American_1983",SPHEROID["GRS_1980",6378137.0,298.257222101]],PRIMEM["Greenwich",0.0],UNIT["Degree",0.01745329251    99433]],PROJECTION["Albers"],PARAMETER["False_Easting",0.0],PARAMETER["False_Northing",0.0],PARAMETER["Central_Meridian",-96.0],PARAMETER["Standard_Parallel_1",29.5],PARAMETER["Standard_Parallel_2",45.5],PARAMETER["Latitude_Of_Origin",2    3.0],UNIT["Foot",0.3048]]`
	errwse := writeTif2("abcd", "crswkt", int(nX), int(nY), minx, miny, xRes, yRes, nodata, grid)
	if errwse != nil {
		panic(errwse)
	}
	//clip to hull
}
func (csv csvHazardProvider) Close() {
	//do nothing?
	csv.ds.Destroy()
}
func (csv csvHazardProvider) ProvideHazard(l geography.Location) (hazards.HazardEvent, error) {

	h := hazards.CoastalEvent{}
	h.SetDepth(123.45) //update from the actual file
	h.SetSalinity(true)
	return h, nil
}
func (csv csvHazardProvider) ProvideHazardBoundary() (geography.BBox, error) {

	return geography.BBox{}, errors.New("stop asking these questions...")
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
