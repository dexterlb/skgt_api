// generated by jsonenums -type=ScheduleType; DO NOT EDIT

package schedules

import (
	"encoding/json"
	"fmt"
)

var (
	_ScheduleTypeNameToValue = map[string]ScheduleType{
		"None": None,
	}

	_ScheduleTypeValueToName = map[ScheduleType]string{
		None: "None",
	}
)

func init() {
	var v ScheduleType
	if _, ok := interface{}(v).(fmt.Stringer); ok {
		_ScheduleTypeNameToValue = map[string]ScheduleType{
			interface{}(None).(fmt.Stringer).String(): None,
		}
	}
}

// MarshalJSON is generated so ScheduleType satisfies json.Marshaler.
func (r ScheduleType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(r).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _ScheduleTypeValueToName[r]
	if !ok {
		return nil, fmt.Errorf("invalid ScheduleType: %d", r)
	}
	return json.Marshal(s)
}

// UnmarshalJSON is generated so ScheduleType satisfies json.Unmarshaler.
func (r *ScheduleType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("ScheduleType should be a string, got %s", data)
	}
	v, ok := _ScheduleTypeNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid ScheduleType %q", s)
	}
	*r = v
	return nil
}
