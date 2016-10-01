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
	schedule, err := GetTimetable(
		htmlparsing.SensibleSettings(),
		&common.Line{
			Vehicle: common.Tram,
			Number:  "10",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	prettyPrint(t, schedule, os.Stdout)
}

func TestGetTimetable_Subway(t *testing.T) {
	schedule, err := GetTimetable(
		htmlparsing.SensibleSettings(),
		&common.Line{
			Vehicle: common.Subway,
			Number:  "1",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	prettyPrint(t, schedule, os.Stdout)
}

func TestAllLines(t *testing.T) {
	lines, err := AllLines(htmlparsing.SensibleSettings())

	if err != nil {
		t.Fatal(err)
	}

	prettyPrint(t, lines, os.Stdout)
}

func TestAllTimetables(t *testing.T) {
	schedules, err := AllTimetables(htmlparsing.SensibleSettings())

	if err != nil {
		t.Fatal(err)
	}

	prettyPrint(t, schedules, os.Stdout)
}
