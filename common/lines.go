package common

import (
	"fmt"
	"regexp"
)

//go:generate jsonenums -type=Transport
type Transport int

const (
	Bus Transport = iota
	Tram
	Trolley
)

type Line struct {
	Type   Transport
	Number string // Why string? For example "4 ТМ"
}

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
		line.Type = Tram
	case "тролей":
		line.Type = Trolley
	case "автобус":
		line.Type = Bus
	default:
		return nil, fmt.Errorf("unknown transport type [%s]", groups[0])
	}

	line.Number = groups[2]

	return line, nil
}
