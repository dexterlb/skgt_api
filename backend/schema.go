package backend

/*
Stop(_id, name<string>, description<string>, location<gps>)

Transport(_id, type<bus, tram, trolley>, number<string>)

Route(_id, transport_id, direction<string>)

RouteStop(route_id, number<int>, stop_id)

Arrival(route_id, stop_id, course<int>, time<int, hour * 60 + minute>, type<workday, holiday etc>)
*/

const schema = `
	create table stop(
		id int primary key,
		name varchar(1024),
		description varchar(2048),
		latitude real,
		longtitude real
	);

	create table line(
		id bigserial primary key,
		vehicle int,
		number varchar(10)
	);

	create table route(
		id bigserial primary key,
		line bigint references line(id),
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
		course int,
		time int,
		day_type int,

		foreign key(route, stop) references route_stop(route, stop)
	);

	create index arrival_route_stop on arrival(route, stop);

	create table api_key(
		value char(64) primary key
	);
`

const dropSchema = `
	drop table api_key;
	drop index arrival_route_stop;
	drop table arrival;
	drop table route_stop;
	drop table route;
	drop table stop;
	drop table line;
`

const clearTransportSchema = `
	delete from arrival;
	delete from route_stop;
	delete from route;
	delete from stop;
	delete from line;
`
