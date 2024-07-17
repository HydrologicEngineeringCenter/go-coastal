package hazardprovider

import (
	"errors"

	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/hazards"
)

type WoodHoleGroupTif struct {
	WSEFilePath  string
	WSE          hazardproviders.HazardProvider
	Wavefilepath string
	Wave         hazardproviders.HazardProvider

	//terrain file?
}

func InitWoodHoleGroupTif(wsefilepath string, wavefilepath string) WoodHoleGroupTif {
	//check input projection and reproject to wgs84?
	whgt := WoodHoleGroupTif{
		WSEFilePath:  wsefilepath,
		WSE:          nil,
		Wavefilepath: wavefilepath,
		Wave:         nil,
	}
	wse, err := hazardproviders.Init(wsefilepath)
	if err != nil {
		panic(err)
	}
	whgt.WSE = wse
	wave, err := hazardproviders.Init(wavefilepath)
	if err != nil {
		panic(err)
	}
	whgt.Wave = wave
	return whgt
}

func (whgt WoodHoleGroupTif) Hazard(l geography.Location) (hazards.HazardEvent, error) {
	c := hazards.CoastalEvent{}
	c.SetSalinity(true)
	d, err := whgt.WSE.Hazard(l)
	if err != nil {
		return c, err
	}
	c.SetDepth(normalizeWSEValues(d.Depth())) //need to pull ground elevation off
	w, err := whgt.Wave.Hazard(l)
	if err != nil {
		return c, err
	}
	c.SetWaveHeight(normalizeWaveValues(w.Depth())) //any actions here? should i reduce it by .7?
	return c, nil
}

// normalizeWSEValues
// 9997 Areas influenced by wave overtopping based flooding only. Cells with this value indicate areas where flooding is caused by intermittent pulses of water from wave overtopping of major coastal structures (e.g., revetments, seawalls) only (i.e., no water directly flows to the location) during simulated events.
// 9998 Dynamic Landform Areas. Cells with this value indicate areas where geomorphology is extremely dynamic and as such expected flooding may vary drastically. These values can appear in any ACFEP level. There are minimal locations of this type and are generally in locations that are regularly flooded and do not have, nor would allow, any type of development.
// 9999 Shallow water flooding during extreme storms. Cells with this value indicate areas where flooding can only be expected during the most extreme events (> 1000-year return period) or where there is only minor water depth during 1000-year return period AEP. These values only appear in 0.1% ACFEP level files.
func normalizeWSEValues(input float64) float64 {
	switch input {
	case 9997:
		return 0.0 //depth no way to convert to wse?
	case 9998:
		panic(errors.New("arrrrgggg wse"))
	case 9999:
		return .5 //depth - no way to convert to wse
	}
	return input
}

// normalizeWaveValues
// 9997 Areas influenced by wave overtopping based flooding only. Cells with this value indicate areas where flooding is caused by intermittent pulses of water from wave overtopping of major coastal structures (e.g., revetments, seawalls) only (i.e., no water directly flows to the location) during simulated events.
// 9998 Dynamic Landform Areas. Cells with this value indicate areas where geomorphology is extremely dynamic and as such expected flooding may vary drastically. These values can appear in any ACFEP level. There are minimal locations of this type and are generally in locations that are regularly flooded and do not have, nor would allow, any type of development.
// 9999 Shallow water flooding during extreme storms. Cells with this value indicate areas where flooding can only be expected during the most extreme events (> 1000-year return period) or where there is only minor water depth during 1000-year return period AEP. These values only appear in 0.1% ACFEP level files.
func normalizeWaveValues(input float64) float64 {
	switch input {
	case 9997:
		return 0 //does this input ever happen in wave?
	case 9998:
		panic(errors.New("arrrrgggg wave"))
	case 9999:
		return 0 //does this happen in the wave files?
	}
	if input > 100 {
		return 0 //correspondence with lisa winter dated 3/9/2023
	}
	return input //could multiply by .7 to reflect "damaging" wave height?
}
func (whgt WoodHoleGroupTif) HazardBoundary() (geography.BBox, error) {
	return whgt.Wave.HazardBoundary()
}

// implement
func (whgt WoodHoleGroupTif) Close() {
	//do nothing?
	whgt.WSE.Close()
	whgt.Wave.Close()
	/*n := time.Since(csv.computeStart)
	fmt.Print("Compute Complete")
	fmt.Print("Compute Time was: ")
	fmt.Println(n)
	fmt.Println(fmt.Sprintf("Processed %v structures, with %v valid depths", csv.queryCount, csv.actualComputedStructures))
	*/
}
