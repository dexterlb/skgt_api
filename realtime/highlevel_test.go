package realtime

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/DexterLB/htmlparsing"
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

func TestArrivals(t *testing.T) {
	arrivals, err := Arrivals(
		htmlparsing.SensibleSettings(),
		2045,
		&Line{
			Type:   Tram,
			Number: 10,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	prettyPrint(t, arrivals, os.Stdout)
}

func TestAllArrivals(t *testing.T) {
	arrivals, err := AllArrivals(
		htmlparsing.SensibleSettings(),
		1700,
	)

	if err != nil {
		t.Fatal(err)
	}

	prettyPrint(t, arrivals, os.Stdout)
}
