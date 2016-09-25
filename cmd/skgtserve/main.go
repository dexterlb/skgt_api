package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/skgt_api/realtime"
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

func main() {
	http.HandleFunc("/arrivals", stopSearch)
	log.Printf("exit: %s", http.ListenAndServe(":8080", nil))
}
