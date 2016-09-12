package realtime

import (
	"fmt"
	"sync"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/skgt_api/common"
)

func Arrivals(settings *htmlparsing.Settings, stopID int, line *common.Line) ([]*Arrival, error) {
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
	Line     *common.Line
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
			Line:     &common.Line{},
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

func GetStopsInfo(settings *htmlparsing.Settings, stops []int, parallelRequests int) ([]*StopInfo, error) {
	in := make(chan int)
	out := make(chan *StopInfo)
	errors := make(chan error)

	go func() {
		for i := range stops {
			in <- stops[i]
		}
		close(in)
	}()

	wg := &sync.WaitGroup{}
	wg.Add(parallelRequests)
	for i := 0; i < parallelRequests; i++ {
		go func() {
			defer wg.Done()

			for stop := range in {
				info, err := GetStopInfo(settings, stop)
				if err != nil {
					errors <- err
				} else {
					out <- info
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	infos := make([]*StopInfo, len(stops))
	go func() {
		i := 0
		for info := range out {
			infos[i] = info
			i++
		}
		close(errors)
	}()

	for err := range errors {
		return nil, fmt.Errorf("unable to get stop info: %s", err)
	}

	return infos, nil
}
