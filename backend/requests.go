package backend

import (
	"fmt"

	"github.com/DexterLB/skgt_api/common"
	"github.com/jmoiron/sqlx"
)

func (b *Backend) Transports() ([]*common.Line, error) {
	var transports []*common.Line
	err := b.db.Select(&transports, GET_ALL_LINES)
	if err != nil {
		return nil, fmt.Errorf("unable to select transports from db: %s", err)
	}

	return transports, nil
}

func (b *Backend) Routes(lineNumber int, vehicleType common.VehicleType) ([]*common.Route, error) {
	data, err := b.Wrap(func(tx *sqlx.Tx) (interface{}, error) {
		var routes []*common.Route
		var stops []*common.Stop
		var directionRouteConnection []struct {
			Direction string
			RouteId   int
		}

		err := tx.Select(&directionRouteConnection, GET_DIRECTION_AND_ROUTE_FOR_LINE, lineNumber, vehicleType)

		if err != nil {
			return nil, fmt.Errorf("unable to select directions for line %d of type %s from db: %s", lineNumber, vehicleType, err)
		}

		for i := range directionRouteConnection {
			stops = nil

			routeId := directionRouteConnection[i].RouteId
			direction := directionRouteConnection[i].Direction

			err = tx.Select(&stops, GET_STOPS_FOR_ROUTE, routeId)
			if err != nil {
				return nil, fmt.Errorf("unable to select routes for line %d of type %s from db: %s", lineNumber, vehicleType, err)
			}

			routes = append(routes, &common.Route{direction, stops})
		}

		return routes, nil
	})

	return data.([]*common.Route), err
}
