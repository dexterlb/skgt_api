package backend

/*
Stop(_id, name<string>, description<string>, location<gps>)

Transport(_id, type<bus, tram, trolley>, number<string>)

Route(_id, transport_id, direction<string>)

RouteStop(route_id, number<int>, stop_id)

Arrival(route_id, stop_id, time<hour, minute>, type<workday, holiday etc>)
*/

const schema = `
	create table person(
		name varchar(50),
		age integer
	);
`

const dropSchema = `
	drop table person;
`
