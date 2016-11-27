package backend

const (
	GET_ALL_LINES = `
		select vehicle, number from line 
		order by vehicle, number;
	`

	GET_DIRECTION_AND_ROUTE_FOR_LINE = `
		select direction, route.id as routeId from route
		left outer join line on line.id = route.line
		where line.number = $1 and line.vehicle = $2;
	`

	GET_STOPS_FOR_ROUTE = `
		select s.* from route_stop r
		left outer join stop s on r.stop = s.id
		where r.route = $1
		order by r.index;
	`
)
