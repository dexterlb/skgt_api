package backend

import (
	"testing"

	"github.com/DexterLB/skgt_api/common"
	"github.com/DexterLB/skgt_api/schedules"
)

func fillDatabase(t *testing.T) *Backend {
	stops := []*common.Stop{
		&common.Stop{
			ID:          1,
			Name:        "foo",
			Description: "FOO",
			Latitude:    42,
			Longtitude:  26,
		},
		&common.Stop{
			ID:          2,
			Name:        "bar",
			Description: "BAR",
			Latitude:    42,
			Longtitude:  26,
		},
		&common.Stop{
			ID:          3,
			Name:        "baz",
			Description: "BAZ",
			Latitude:    42,
			Longtitude:  26,
		},
	}
	timetables := []*schedules.Timetable{
		&schedules.Timetable{
			Line: &common.Line{
				Vehicle: common.Tram,
				Number:  "10",
			},
			Routes: []*schedules.Route{
				&schedules.Route{
					Direction: "A - B",
					Stops:     []int{1, 2, 3},
					Schedules: map[schedules.ScheduleType][]schedules.Course{
						schedules.Workday: []schedules.Course{
							schedules.Course{
								schedules.NewTime(12, 0),
								schedules.NewTime(12, 30),
								schedules.NewTime(13, 0),
							},
							schedules.Course{
								schedules.NewTime(13, 0),
								schedules.NewTime(13, 30),
								schedules.NewTime(14, 0),
							},
						},
						schedules.HolidayAndPreHoliday: []schedules.Course{
							schedules.Course{
								schedules.NewTime(14, 0),
								schedules.NewTime(14, 30),
							},
						},
					},
				},
			},
		},
	}

	backend := openBackend(t)

	err := backend.Fill(stops, timetables)
	if err != nil {
		t.Fatalf("unable to fill database: %s", err)
	}

	return backend
}

func TestBackend_Fill(t *testing.T) {
	backend := fillDatabase(t)
	defer closeBackend(t, backend)
}
