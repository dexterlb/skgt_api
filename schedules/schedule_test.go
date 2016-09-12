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

func TestGetSchedule(t *testing.T) {
	schedule, err := GetScheduleInfo(
		htmlparsing.SensibleSettings(),
		&common.Line{
			Type:   common.Tram,
			Number: "10",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	prettyPrint(t, schedule, os.Stdout)
}
