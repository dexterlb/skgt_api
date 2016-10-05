package server

import (
	"fmt"
	"strconv"

	"github.com/DexterLB/skgt_api/realtime"
	"github.com/julienschmidt/httprouter"
)

func (s *Server) realtimeArrivals(params httprouter.Params) (interface{}, error) {
	stopID, err := strconv.Atoi(params.ByName("stop_id"))
	if err != nil {
		return nil, fmt.Errorf("unable to parse stop ID: %s", err)
	}

	arrivals, err := realtime.AllArrivals(s.parserSettings, stopID)

	return arrivals, err
}
