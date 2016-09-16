package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/skgt_api/realtime"
	"github.com/DexterLB/skgt_api/schedules"
)

func stopSearch(w http.ResponseWriter, req *http.Request) {
	log.Printf("start stop search")
	defer log.Printf("end stop search")

	parameters := req.URL.Query()

	key := parameters.Get("key")
	if key != "42" {
		http.Error(w, "fuck you", 403)
		return
	}

	stop, err := strconv.Atoi(parameters.Get("stop"))
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse stop: %s", err), 500)
		return
	}

	arrivals, err := realtime.AllArrivals(htmlparsing.SensibleSettings(), stop)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to get stop arrivals: %s", err), 500)
		return
	}

	info, err := realtime.GetStopInfo(htmlparsing.SensibleSettings(), stop)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to get stop info: %s", err), 500)
		return
	}

	object := &struct {
		Stop     *realtime.StopInfo
		Arrivals []*realtime.LineArrivals
	}{
		Stop:     info,
		Arrivals: arrivals,
	}

	data, err := json.MarshalIndent(object, "", "    ")
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to marshal data: %s", err), 500)
		return
	}

	fmt.Fprintf(w, string(data))
}

type infoCache struct {
	sync.Mutex

	info *Info
}

type Info struct {
	Stops      []*realtime.StopInfo
	Timetables []*schedules.Timetable
}

func (i *infoCache) GetInfo() (*Info, error) {
	i.Lock()
	defer i.Unlock()

	log.Printf("start get info")
	defer log.Printf("end get info")

	if i.info == nil {
		timetables, err := schedules.AllTimetables(htmlparsing.SensibleSettings())
		if err != nil {
			return nil, fmt.Errorf("unable to get schedules: %s", err)
		}

		stopInfos, err := realtime.GetStopsInfo(
			htmlparsing.SensibleSettings(),
			schedules.GetStops(timetables),
			8,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to get stops: %s", err)
		}

		i.info = &Info{
			Stops:      stopInfos,
			Timetables: timetables,
		}
	}

	return i.info, nil
}

func (i *infoCache) InfoRequest(w http.ResponseWriter, req *http.Request) {
	parameters := req.URL.Query()

	key := parameters.Get("key")
	if key != "42" {
		http.Error(w, "fuck you", 403)
		return
	}

	info, err := i.GetInfo()
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to get info: %s", err), 500)
		return
	}

	data, err := json.MarshalIndent(info, "", "    ")
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to marshal data: %s", err), 500)
		return
	}

	fmt.Fprintf(w, string(data))
}

func main() {
	cache := &infoCache{}

	go func() {
		_, err := cache.GetInfo()
		if err != nil {
			log.Printf("initial get info error: %s", err)
		}
	}()

	http.HandleFunc("/arrivals", stopSearch)
	http.HandleFunc("/info", cache.InfoRequest)
	http.ListenAndServe(":8080", nil)
}
