package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/skgthack/realtime"
)

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
		http.Error(w, fmt.Sprintf("unable to get data: %s", err), 500)
		return
	}

	data, err := json.MarshalIndent(arrivals, "", "    ")
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to marshal data: %s", err), 500)
		return
	}

	fmt.Fprintf(w, string(data))
}

func main() {
	http.HandleFunc("/", stopSearch)
	http.ListenAndServe(":8080", nil)
}
