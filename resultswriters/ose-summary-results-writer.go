package resultswriters

import (
	"fmt"
	"io"
	"os"

	"github.com/USACE/go-consequences/consequences"
)

type oseResultsWriter struct {
	filepath       string
	frequencies    []float64
	w              io.Writer
	totalsbyfreq   map[string][]int32
	frequencyIndex int
}

func InitOseResultsWriterFromFile(filepath string, frequencies []float64) (*oseResultsWriter, error) {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return &oseResultsWriter{}, err
	}
	//make the maps
	t := make(map[string][]int32, 1)
	return &oseResultsWriter{filepath: filepath, frequencies: frequencies, w: w, totalsbyfreq: t}, nil
}
func (srw *oseResultsWriter) SetFrequencyIndex(index int) {
	srw.frequencyIndex = index
}
func (srw *oseResultsWriter) Write(r consequences.Result) {
	//hardcoding for structures to experiment and think it through.
	var cat = "damage category"
	var ot = "occupancy type"
	var jobs = "pop2pmu65"
	var j int32 = 0
	var o int32 = 0
	var u int32 = 0
	var u65 = "pop2amu65"
	var o65 = "pop2amo65"
	var damcat = ""
	var occtype = ""
	h := r.Headers
	for i, v := range h {
		if v == cat {
			damcat = r.Result[i].(string)
		}
		if v == ot {
			occtype = r.Result[i].(string)
		}
		if v == jobs {
			j = r.Result[i].(int32)
		}
		if v == u65 {
			o = r.Result[i].(int32)
		}
		if v == o65 {
			u = r.Result[i].(int32)
		}

	}
	//update damcat totals.
	t, ok := srw.totalsbyfreq[damcat]
	if ok {
		t[srw.frequencyIndex] += 1
		srw.totalsbyfreq[damcat] = t
	} else {
		dc := make([]int32, len(srw.frequencies))
		dc[srw.frequencyIndex] = 1
		srw.totalsbyfreq[damcat] = dc
	}
	//update occtype totals.
	if occtype == "GOV2" {
		t, ok := srw.totalsbyfreq[occtype]
		if ok {
			t[srw.frequencyIndex] += 1
			srw.totalsbyfreq[occtype] = t
		} else {
			dc := make([]int32, len(srw.frequencies))
			dc[srw.frequencyIndex] = 1
			srw.totalsbyfreq[occtype] = dc
		}
	}
	if occtype == "COM6" {
		t, ok := srw.totalsbyfreq[occtype]
		if ok {
			t[srw.frequencyIndex] += 1
			srw.totalsbyfreq[occtype] = t
		} else {
			dc := make([]int32, len(srw.frequencies))
			dc[srw.frequencyIndex] = 1
			srw.totalsbyfreq[occtype] = dc
		}
	}
	//update totals.
	tot, ok := srw.totalsbyfreq["Total Structure Count"]
	if ok {
		tot[srw.frequencyIndex] += 1
		srw.totalsbyfreq["Total Structure Count"] = tot
	} else {
		dc := make([]int32, len(srw.frequencies))
		dc[srw.frequencyIndex] = 1
		srw.totalsbyfreq["Total Structure Count"] = dc
	}
	//update jobs.
	if damcat != "Res" {
		jbs, ok := srw.totalsbyfreq["Jobs Impacted"]
		if ok {
			jbs[srw.frequencyIndex] += j
			srw.totalsbyfreq["Jobs Impacted"] = jbs
		} else {
			dc := make([]int32, len(srw.frequencies))
			dc[srw.frequencyIndex] = j
			srw.totalsbyfreq["Jobs Impacted"] = dc
		}
	}
	//update u65.
	under, ok := srw.totalsbyfreq["Population Impacted Under 65 Day"]
	if ok {
		under[srw.frequencyIndex] += u
		srw.totalsbyfreq["Population Impacted Under 65 Day"] = under
	} else {
		dc := make([]int32, len(srw.frequencies))
		dc[srw.frequencyIndex] = u
		srw.totalsbyfreq["Population Impacted Under 65 Day"] = dc
	}
	//update o65.
	over, ok := srw.totalsbyfreq["Population Impacted Over 65 Day"]
	if ok {
		over[srw.frequencyIndex] += o
		srw.totalsbyfreq["Population Impacted Over 65 Day"] = over
	} else {
		dc := make([]int32, len(srw.frequencies))
		dc[srw.frequencyIndex] = u
		srw.totalsbyfreq["Population Impacted Over 65 Day"] = dc
	}
}
func (srw *oseResultsWriter) Close() {
	headerstring := "Category of Impact\\Frequency of Flooding"
	for _, f := range srw.frequencies {
		headerstring = fmt.Sprintf("%v, %v", headerstring, f)
	}
	fmt.Fprintln(srw.w, headerstring) //needs new line
	h := srw.totalsbyfreq
	for i, v := range h {
		rowstring := i
		for _, f := range v {
			rowstring = fmt.Sprintf("%v, %v", rowstring, f)
		}
		fmt.Fprintln(srw.w, rowstring) //needs new line
	}
	w2, ok := srw.w.(io.WriteCloser)
	if ok {
		w2.Close()
	}
}
