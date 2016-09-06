package realtime

import (
	"fmt"
	"testing"

	"github.com/DexterLB/htmlparsing"
)

func TestLookupStop(t *testing.T) {
	data, err := LookupStop(htmlparsing.SensibleSettings(), 2327)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("data: %v\n", data)
}
