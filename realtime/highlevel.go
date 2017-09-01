package realtime

import (
	"fmt"
	"log"
	"sync"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/skgt_api/common"
)

// Arrivals returns all arrivals on the given line, at the given stop in the next
// hour or so
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

// LineArrivals pairs a line with a list of arrivals
type LineArrivals struct {
	Line     *common.Line
	Arrivals []*Arrival
}

// AllArrivals returns all arrivals at a given stop in the next hour or so
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

// GetStopInfo gets information for the given stop ID
func GetStopInfo(settings *htmlparsing.Settings, stopID int) (*common.Stop, error) {
	data, err := LookupStop(settings, stopID)
	if err != nil {
		return nil, fmt.Errorf("unable to get stop data: %s", err)
	}

	return &common.Stop{
		ID:          stopID,
		Name:        data.Name,
		Description: data.Description,
	}, nil
}

// GetStopsInfo gets information for multiple stops, making at most
// parallelRequests requests in a single moment
func GetStopsInfo(settings *htmlparsing.Settings, stops []int, parallelRequests int) ([]*common.Stop, error) {
	infos := make([]*common.Stop, len(stops))
	for i := range stops {
		infos[i] = &common.Stop{
			ID: stops[i],
		}
	}

	err := UpdateStopsInfo(settings, infos, parallelRequests)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

// GetStopsInfo updates the information for multiple stops, making at most
// parallelRequests requests in a single moment
func UpdateStopsInfo(
	settings *htmlparsing.Settings,
	stops []*common.Stop,
	parallelRequests int,
) error {
	in := make(chan *common.Stop, parallelRequests*2)
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
				err := UpdateStopInfo(settings, stop)
				if err != nil {
					errors <- err
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errors)
	}()

	for err := range errors {
		return fmt.Errorf("unable to get stop info: %s", err)
	}

	return nil
}

// UpdateStopInfo updates the stop information with data from the site
func UpdateStopInfo(settings *htmlparsing.Settings, info *common.Stop) error {
	newInfo, err := GetStopInfo(settings, info.ID)
	if err != nil {
		return err
	}

	if info.Name != newInfo.Name {
		if info.Name == "" {
			info.Name = newInfo.Name
			log.Printf("warning: stop %04d: no name, assigning VT name: %s", info.ID, info.Name)
		} else if newInfo.Name != "" {
			log.Printf("warning: stop %04d: preferring '%s' over '%s' (VT)", info.ID, info.Name, newInfo.Name)
		} else {
			log.Printf("warning: stop %04d: no VT name, leaving %s", info.ID, info.Name)
		}
	} else if info.Name == "" && newInfo.Name == "" {
		log.Printf("warning: stop %04d has no name!")
	}

	info.Description = newInfo.Description
	info.Longitude = newInfo.Longitude
	info.Latitude = newInfo.Latitude

	return nil
}
