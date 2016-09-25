package common

import (
	"fmt"
	"regexp"
)

// VehicleType is an enum for types of vehicles
type VehicleType int

//go:generate jsonenums -type=VehicleType

const (
	// Bus is a fossil-fuel powered road vehicle
	Bus VehicleType = iota
	// Tram is an electric rail vehicle
	Tram
	// Trolley is an electric road vehicle
	Trolley
)

// Line is a transport line (e.g. "Tram 10", "Bus 94" etc)
type Line struct {
	Vehicle VehicleType
	Number  string // Why string? For example "4 ТМ"
}

// ParseLine parses a human-readable string like "трамвай 10"
func ParseLine(input string) (*Line, error) {
	groups := regexp.MustCompile(
		`([^\s]+) (.+)`,
	).FindStringSubmatch(input)
	if len(groups) < 3 {
		return nil, fmt.Errorf("unable to parse line info [%s]", input)
	}

	line := &Line{}

	switch groups[1] {
	case "трамвай":
		line.Vehicle = Tram
	case "тролей":
		line.Vehicle = Trolley
	case "автобус":
		line.Vehicle = Bus
	default:
		return nil, fmt.Errorf("unknown transport type [%s]", groups[0])
	}

	line.Number = groups[2]

	return line, nil
}
