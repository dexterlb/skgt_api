package backend

import (
	"fmt"

	"github.com/DexterLB/skgt_api/common"
	"github.com/DexterLB/skgt_api/schedules"
	"github.com/jmoiron/sqlx"
)

// Fill populstes the database with the given stops and timetables
// (replacing all previous content)
func (b *Backend) Fill(stops []*common.Stop, timetables []*schedules.Timetable) (err error) {
	tx, err := b.db.Beginx()
	if err != nil {
		return fmt.Errorf("cannot open transaction: %s", err)
	}

	defer func() {
		// if something has failed, and tx hasn't been commited, roll it back
		if tx != nil {
			txErr := tx.Rollback()
			if txErr != nil {
				err = fmt.Errorf("trying to rollback transaction because of error [%s] failed: %s", err, txErr)
			}
		}
	}()

	// first clear all of the old data
	_, err = tx.Exec(clearTransportSchema)
	if err != nil {
		return fmt.Errorf("cannot clear data before inserting new data: %s", err)
	}

	// insert stops
	for i := range stops {
		err = insertStop(tx, stops[i])
		if err != nil {
			return fmt.Errorf("unable to insert stop: %s", err)
		}
	}

	// insert timetables
	for i := range timetables {
		err = insertTimetable(tx, timetables[i])
		if err != nil {
			return fmt.Errorf("unable to insert timetable: %s", err)
		}
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cannot commit transaction: %s", err)
	}

	// set tx to nil so that the deferred error check knows not to rollback
	tx = nil

	return nil
}

func insertStop(tx *sqlx.Tx, stop *common.Stop) error {
	_, err := tx.NamedExec(
		`insert into stop(id, name, description, latitude, longtitude)
	     values (:id, :name, :description, :latitude, :longtitude)`,
		stop,
	)
	return err
}

func insertTimetable(tx *sqlx.Tx, timetable *schedules.Timetable) error {
	var lineID uint64
	err := tx.Get(
		&lineID,
		`insert into line(id, vehicle, number)
		 values(default, $1, $2) returning id`,
		timetable.Line.Vehicle, timetable.Line.Number,
	)
	if err != nil {
		return err
	}

	for i := range timetable.Routes {
		err = insertRoute(tx, timetable.Routes[i], lineID)
		if err != nil {
			return fmt.Errorf("unable to insert route: %s", err)
		}
	}

	return nil
}

func insertRoute(tx *sqlx.Tx, route *schedules.Route, lineID uint64) error {
	var routeID uint64
	err := tx.Get(
		&routeID,
		`insert into route(id, line, direction)
		 values(default, $1, $2) returning id`,
		lineID, route.Direction,
	)
	if err != nil {
		return err
	}

	for i := range route.Stops {
		_, err = tx.Exec(
			`insert into route_stop(route, index, stop)
			 values($1, $2, $3)`,
			routeID, i+1, route.Stops[i],
		)
		if err != nil {
			return fmt.Errorf("unable to insert route-stop connection: %s", err)
		}
	}

	for scheduleType := range route.Schedules {
		for courseIndex, course := range route.Schedules[scheduleType] {
			for stopIndex := range course {
				_, err = tx.Exec(
					`insert into arrival(route, stop, course, time, day_type)
				 	 values($1, $2, $3, $4, $5)`,
					routeID,
					route.Stops[stopIndex],
					courseIndex+1,
					course[stopIndex],
					scheduleType,
				)
			}
			if err != nil {
				return fmt.Errorf("unable to insert arrival: %s", err)
			}
		}
	}

	return nil
}
