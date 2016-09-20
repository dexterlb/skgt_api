package backend

import (
	"fmt"

	"github.com/DexterLB/skgt_api/realtime"
	"github.com/DexterLB/skgt_api/schedules"
	"github.com/jmoiron/sqlx"
)

func (b *Backend) Fill(stops []*realtime.StopInfo, timetables []*schedules.Timetable) (err error) {
	tx, err := b.db.Beginx()
	if err != nil {
		return fmt.Errorf("cannot open transaction: %s", err)
	}

	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	for i := range stops {
		err = insertStop(tx, stops[i])
		if err != nil {
			return fmt.Errorf("unable to insert stop: %s", err)
		}
	}

	for i := range timetables {
		err = insertTimetable(tx, timetables[i])
		if err != nil {
			return fmt.Errorf("unable to insert timetable: %s", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cannot commit transaction: %s", err)
	}
	tx = nil

	return nil
}

func insertStop(tx *sqlx.Tx, stop *realtime.StopInfo) error {
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

	for _, route := range timetable.Routes {

	}

	return nil
}
