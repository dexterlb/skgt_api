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

func TestLookupStop(t *testing.T) {
	data, err := LookupStop(htmlparsing.SensibleSettings(), 1700)
	if err != nil {
		t.Fatal(err)
	}

	data.CaptchaResult, err = htmlparsing.BreakSimpleCaptcha(data.Captcha)
	if err != nil {
		t.Fatal(err)
	}

	arrivals, err := data.Arrivals(57)
	if err != nil {
		t.Fatal(err)
	}

	prettyPrint(t, arrivals, os.Stdout)
}
