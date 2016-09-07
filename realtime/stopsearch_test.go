package realtime

import (
	"testing"

	"github.com/DexterLB/htmlparsing"
)

func TestLookupStop(t *testing.T) {
	_, err := LookupStop(htmlparsing.SensibleSettings(), 1700)
	if err != nil {
		t.Fatal(err)
	}
}
