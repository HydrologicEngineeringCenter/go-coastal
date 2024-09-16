package hazardprovider

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazards"
)

type IanCSVFileProvider struct {
	data        map[geography.Location]hazards.CoastalEvent
	boundingBox geography.BBox
}

func InitIanCSVFileProvider(fp string) (*IanCSVFileProvider, error) {
	f, err := os.Open(fp)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	var MinX, MinY, MaxX, MaxY float64
	MinX = -180
	MinY = -180
	MaxX = 180
	MaxY = 180
	hp := IanCSVFileProvider{}
	scanner := bufio.NewScanner(f)
	scanner.Scan() //burn the header.
	data := make(map[geography.Location]hazards.CoastalEvent)
	for scanner.Scan() {
		lines := strings.Split(scanner.Text(), ",")
		x, err := strconv.ParseFloat(lines[0], 64)
		if err != nil {
			return &hp, errors.New("could not parse x")
		}
		y, err := strconv.ParseFloat(lines[1], 64)
		if err != nil {
			return &hp, errors.New("could not parse y")
		}
		if MinX > x {
			MinX = x
		}
		if MinY > y {
			MinY = y
		}
		if MaxX < x {
			MaxX = x
		}
		if MaxY < y {
			MaxY = y
		}
		location := geography.Location{
			X: x,
			Y: y,
			//SRID: "4326",
		}
		ground_elev_ft, err := strconv.ParseFloat(lines[2], 64)
		if err != nil {
			return &hp, errors.New("could not parse ground elevation")
		}
		surge_elev_m, err := strconv.ParseFloat(lines[3], 64)
		if err != nil {
			return &hp, errors.New("could not parse surge elevation")
		}
		wave_elev_m, err := strconv.ParseFloat(lines[4], 64)
		if err != nil {
			return &hp, errors.New("could not parse wave elevation")
		}
		ch := hazards.CoastalEvent{}
		surge_elev_ft := surge_elev_m * 3.28084     //convert from meters to feet
		wave_elev_ft := wave_elev_m * 3.28084       //convert from meters to feet
		wave_height := wave_elev_ft - surge_elev_ft //subtract surge component from wave component to derive wave height.
		depth := surge_elev_ft - ground_elev_ft     //subtract ground component from surge component to derive surge depth.
		ch.SetDepth(depth)
		ch.SetSalinity(true) //assume saline water in case it is ever used in loss functions, and to differentiate from non coastal hazards.
		ch.SetWaveHeight(wave_height)
		data[location] = ch
	}

	bbox := make([]float64, 4)
	bbox[0] = MinX //upper left x
	bbox[1] = MaxY //upper left y
	bbox[2] = MaxX //lower right x
	bbox[3] = MinY //lower right y
	hp.boundingBox = geography.BBox{Bbox: bbox}
	hp.data = data
	return &hp, nil
}
func (csv *IanCSVFileProvider) HazardBoundary() (geography.BBox, error) {
	return csv.boundingBox, nil
}
func (csv *IanCSVFileProvider) Hazard(l geography.Location) (hazards.HazardEvent, error) {
	he := csv.data[l]
	if he.Has(hazards.WaveHeight) {
		return he, nil
	}
	return he, errors.New("the hazard event has no wave")
}
func (csv *IanCSVFileProvider) Close() {
	fmt.Println("closing hazard provider")
}
