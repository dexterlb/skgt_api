package schedules

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/skgt_api/common"
)

func prettyPrint(t *testing.T, data interface{}, w io.Writer) {
	s, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	_, err = fmt.Fprintf(w, "%s\n", string(s))
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTimetable(t *testing.T) {
	stops := make(map[int]string)
	stopNames := make(chan *StopName)

	var timetable *Timetable
	go func(timetable **Timetable) {
		var err error

		*timetable, err = GetTimetable(
			htmlparsing.SensibleSettings(),
			&common.Line{
				Vehicle: common.Tram,
				Number:  "10",
			},
			stopNames,
		)

		if err != nil {
			t.Fatal(err)
		}

		close(stopNames)
	}(&timetable)

	for stop := range stopNames {
		stops[stop.ID] = stop.Name
	}

	prettyPrint(
		t,
		struct {
			Timetable *Timetable
			Stops     map[int]string
		}{
			timetable,
			stops,
		},
		os.Stdout,
	)
}

func TestGetTimetable_Subway(t *testing.T) {
	stops := make(map[int]string)
	stopNames := make(chan *StopName)

	var timetable *Timetable
	go func(timetable **Timetable) {
		var err error

		*timetable, err = GetTimetable(
			htmlparsing.SensibleSettings(),
			&common.Line{
				Vehicle: common.Subway,
				Number:  "1",
			},
			stopNames,
		)

		if err != nil {
			t.Fatal(err)
		}

		close(stopNames)
	}(&timetable)

	for stop := range stopNames {
		stops[stop.ID] = stop.Name
	}

	prettyPrint(
		t,
		struct {
			Timetable *Timetable
			Stops     map[int]string
		}{
			timetable,
			stops,
		},
		os.Stdout,
	)
}

func TestAllLines(t *testing.T) {
	lines, err := AllLines(htmlparsing.SensibleSettings())

	if err != nil {
		t.Fatal(err)
	}

	prettyPrint(t, lines, os.Stdout)
}

func TestAllTimetables(t *testing.T) {
	timetables, stops, err := AllTimetables(htmlparsing.SensibleSettings(), 8)

	if err != nil {
		t.Fatal(err)
	}

	prettyPrint(
		t,
		struct {
			Timetables []*Timetable
			Stops      []*common.Stop
		}{
			timetables,
			stops,
		},
		os.Stdout,
	)
}
