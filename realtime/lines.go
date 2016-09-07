package realtime

import (
	"fmt"
	"regexp"
	"strconv"
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
	Number int
}

func parseLine(input string) (*Line, error) {
	groups := regexp.MustCompile(
		`(.+) ([\d]+)`,
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

	var err error
	line.Number, err = strconv.Atoi(groups[2])
	if err != nil {
		return nil, fmt.Errorf("unable to parse line number: %s", err)
	}

	return line, nil
}
