package backend

/*
Stop(_id, name<string>, description<string>, location<gps>)

Transport(_id, type<bus, tram, trolley>, number<string>)

Route(_id, transport_id, direction<string>)

RouteStop(route_id, number<int>, stop_id)

Arrival(route_id, stop_id, time<hour, minute>, type<workday, holiday etc>)
*/

type Stop struct {
	ID          uint64
	Name        string
	Description string
	Latitude    float32
	Longtitude  float32
}

type Transport struct {
	ID      uint64
	Vehicle common.Vehicle
	Number  string
}

const schema = `
	create table stop(
		id bigserial primary key,
		name varchar(1024),
		description varchar(2048),
		latitude real,
		longtitude real
	);

	create table transport(
		id bigserial primary key,
		vehicle int,
		number varchar(10)
	);

	create table route(
		id bigserial primary key,
		transport bigint references transport(id),
		direction varchar(1024)
	);

	create table route_stop(
		route bigint references route(id),
		index int,
		stop bigint references stop(id),

		primary key(route, stop)
	);

	create table arrival(
		route bigint not null,
		stop bigint not null,
		time int,
		day_type int,

		foreign key(route, stop) references route_stop(route, stop)
	);

	create table api_key(
		value char(256) primary key
	);
`

const dropSchema = `
	drop table api_key;
	drop table arrival;
	drop table route_stop;
	drop table route;
	drop table stop;
	drop table transport;
`
