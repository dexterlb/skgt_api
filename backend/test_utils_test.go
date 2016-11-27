package backend

import (
	"encoding/json"
	"reflect"
	"testing"
)

func assertEqualJSON(expected interface{}, actual interface{}, t *testing.T) {
	areEqual := reflect.DeepEqual(expected, actual)
	data := &struct {
		Expected interface{}
		Actual   interface{}
	}{
		Expected: expected,
		Actual:   actual,
	}

	if !areEqual {
		jsonData, err := json.MarshalIndent(data, "", "    ")
		if err == nil {
			t.Errorf("things are not equal: %s", string(jsonData))
		} else {
			t.Errorf("Cant jsonise data %s", err)
		}
	}
}
