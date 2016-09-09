package realtime

import (
	"fmt"

	"github.com/DexterLB/htmlparsing"
)

func Arrivals(settings *htmlparsing.Settings, stopID int, line *Line) ([]*Arrival, error) {
	data, err := LookupStop(settings, stopID)
	if err != nil {
		return nil, fmt.Errorf("unable to get stop data: %s", err)
	}

	err = data.BreakCaptcha()
	if err != nil {
		return nil, err
	}

	lineID, ok := data.Lines[*line]
	if !ok {
		return nil, fmt.Errorf("no such line: %v", line)
	}

	return data.Arrivals(lineID)
}

type LineArrivals struct {
	Line     *Line
	Arrivals []*Arrival
}

func AllArrivals(settings *htmlparsing.Settings, stopID int) ([]*LineArrivals, error) {
	data, err := LookupStop(settings, stopID)
	if err != nil {
		return nil, fmt.Errorf("unable to get stop data: %s", err)
	}

	lineArrivals := make([]*LineArrivals, len(data.Lines))

	i := 0
	for line, lineID := range data.Lines {
		err = data.BreakCaptcha()
		if err != nil {
			return nil, err
		}

		arrivals, err := data.Arrivals(lineID)
		if err != nil {
			return nil, fmt.Errorf("unable to get arrivals: %s", err)
		}

		lineArrivals[i] = &LineArrivals{
			Line:     &Line{},
			Arrivals: arrivals,
		}
		*(lineArrivals[i].Line) = line

		i++
	}

	return lineArrivals, nil
}

type StopInfo struct {
	ID          int
	Name        string
	Description string
}

func GetStopInfo(settings *htmlparsing.Settings, stopID int) (*StopInfo, error) {
	data, err := LookupStop(settings, stopID)
	if err != nil {
		return nil, fmt.Errorf("unable to get stop data: %s", err)
	}

	return &StopInfo{
		ID:          stopID,
		Name:        data.Name,
		Description: data.Description,
	}, nil
}
