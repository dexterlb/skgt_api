package openstreetmap

import (
	"testing"

	"github.com/DexterLB/htmlparsing"
)

func TestGetStops(t *testing.T) {
	stops, err := GetStops(htmlparsing.SensibleSettings())
	if err != nil {
		t.Fatal(err)
	}

	if len(stops) < 10 {
		t.Errorf("Too few stops. Something's fishy.")
	}
}
