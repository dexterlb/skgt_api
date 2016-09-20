package backend

import (
	"testing"

	"github.com/DexterLB/skgt_api/common"
	"github.com/DexterLB/skgt_api/realtime"
	"github.com/DexterLB/skgt_api/schedules"
)

func TestBackend_Fill(t *testing.T) {
	stops := []*realtime.StopInfo{
		&realtime.StopInfo{
			ID:          1,
			Name:        "foo",
			Description: "FOO",
			Latitude:    42,
			Longtitude:  26,
		},
		&realtime.StopInfo{
			ID:          2,
			Name:        "bar",
			Description: "BAR",
			Latitude:    42,
			Longtitude:  26,
		},
		&realtime.StopInfo{
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
	defer closeBackend(t, backend)

	err := backend.Fill(stops, timetables)
	if err != nil {
		t.Fatalf("unable to fill database: %s", err)
	}

	pause()
}
