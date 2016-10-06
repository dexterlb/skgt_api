package server

import (
	"fmt"

	"github.com/DexterLB/skgt_api/common"
	"github.com/julienschmidt/httprouter"
)

func (s *Server) transports(params httprouter.Params) (interface{}, error) {
	transports, err := s.backend.Transports()

	return transports, err
}

func (s *Server) routes(params httprouter.Params) (interface{}, error) {
	number := params.ByName("number")
	vehicle, err := common.ParseVehicle(params.ByName("vehicle"))

	if err != nil {
		return nil, fmt.Errorf("could not parse vehicle type: %s", err)
	}

	routes, err := s.backend.Routes(number, vehicle)
	if err != nil {
		return nil, fmt.Errorf("could not get routes: %s", err)
	}

	return routes, nil
}
