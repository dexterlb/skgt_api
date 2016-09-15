package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/skgt_api/realtime"
	"github.com/DexterLB/skgt_api/schedules"
)

func info(w http.ResponseWriter, req *http.Request) {
	parameters := req.URL.Query()

	key := parameters.Get("key")
	if key != "42" {
		http.Error(w, "fuck you", 403)
		return
	}

	scheduleInfos, err := schedules.AllSchedules(htmlparsing.SensibleSettings())
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to get schedules: %s", err), 500)
		return
	}

	stopInfos, err := realtime.GetStopsInfo(
		htmlparsing.SensibleSettings(),
		schedules.GetStops(scheduleInfos),
		8,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to get stops: %s", err), 500)
		return
	}

	object := &struct {
		Stops     []*realtime.StopInfo
		Schedules []*schedules.ScheduleInfo
	}{
		Stops:     stopInfos,
		Schedules: scheduleInfos,
	}

	data, err := json.MarshalIndent(object, "", "    ")
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to marshal data: %s", err), 500)
		return
	}

	fmt.Fprintf(w, string(data))
}

func stopSearch(w http.ResponseWriter, req *http.Request) {
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

func main() {
	http.HandleFunc("/arrivals", stopSearch)
	http.HandleFunc("/info", info)
	http.ListenAndServe(":8080", nil)
}
