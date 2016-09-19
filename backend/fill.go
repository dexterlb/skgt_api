package backend

import (
	"fmt"

	"github.com/DexterLB/skgt_api/realtime"
	"github.com/DexterLB/skgt_api/schedules"
)

func (b *Backend) Fill(stops []*realtime.StopInfo, timetables []*schedules.Timetable) error {
	tx, err := b.db.Beginx()
	if err != nil {
		return fmt.Errorf("cannot open transaction: %s", err)
	}

	for i := range stops {
		_, err = tx.NamedExec(
			`insert into stop(id, name, description, latitude, longtitude)
			 values (:id, :name, :description, :latitude, :longtitude)`,
			stops[i],
		)
		if err != nil {
			return fmt.Errorf("unable to insert stop: %s", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cannot commit transaction: %s", err)
	}

	return nil
}
