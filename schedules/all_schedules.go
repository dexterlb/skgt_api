package schedules

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/skgt_api/common"
)

// AllTimetables returns the timetables for all lines and all stops,
// making at most parrallelRequests http requests at the same time.
// Stops might not have the full infurmation, and need to be processed
// further.
func AllTimetables(
	settings *htmlparsing.Settings,
	parallelRequests int,
) (
	[]*Timetable,
	[]*common.Stop,
	error,
) {
	lines, err := AllLines(settings)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get list of lines")
	}

	in := make(chan *common.Line, parallelRequests*2)
	out := make(chan *Timetable, parallelRequests*2)
	stops := make(chan *StopName, parallelRequests*2)
	errors := make(chan error)

	go func() {
		for i := range lines {
			in <- lines[i]
		}
		close(in)
	}()

	wg := &sync.WaitGroup{}
	wg.Add(parallelRequests)
	for i := 0; i < parallelRequests; i++ {
		go func() {
			defer wg.Done()

			for line := range in {
				info, err := GetTimetable(settings, line, stops)
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

	timetables := make([]*Timetable, len(lines))
	go func() {
		i := 0
		for timetable := range out {
			timetables[i] = timetable
			i++
		}
		close(stops)
	}()

	stopNameSet := make(map[int]string)
	go func() {
		for stop := range stops {
			stopNameSet[stop.ID] = stop.Name
		}
		close(errors)
	}()

	for err := range errors {
		return nil, nil, fmt.Errorf("unable to get stop info: %s", err)
	}

	return timetables, stopNamesToStops(stopNameSet), nil
}

// AllLines returns all lines
func AllLines(settings *htmlparsing.Settings) ([]*common.Line, error) {
	page, err := htmlparsing.NewClient(settings).ParsePage(
		`https://schedules.sofiatraffic.bg/`, nil,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to parse line list page: %s", err)
	}

	links, err := page.Search(
		`//div[contains(@class, 'lines_section')]/ul/li/a | //a[contains(@class, 'quicksearch') and not(contains(@href, './'))]`,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to get line links: %s", err)
	}

	lines := make([]*common.Line, len(links))
	for i := range links {
		href, ok := links[i].Attributes()["href"]
		if !ok {
			return nil, fmt.Errorf("link has no href")
		}
		lines[i], err = parseLine(href.Value())
		if err != nil {
			return nil, fmt.Errorf("unable to parse line: %s", err)
		}
	}

	return lines, nil
}

// GetStops returns the IDs of all stops mentioned in a list of timetables
func GetStops(timetables []*Timetable) []int {
	stopSet := make(map[int]struct{})
	for _, timetable := range timetables {
		for _, route := range timetable.Routes {
			for _, stop := range route.Stops {
				stopSet[stop] = struct{}{}
			}
		}
	}

	stops := make([]int, len(stopSet))
	i := 0
	for stop := range stopSet {
		stops[i] = stop
		i++
	}

	return stops
}

// parseLine parses a line from a link such as
// "https://schedules.sofiatraffic.bg/autobus/18"
func parseLine(link string) (*common.Line, error) {
	originalLink, err := url.QueryUnescape(link)
	if err != nil {
		return nil, fmt.Errorf("invalid link: %s", err)
	}

	groups := strings.Split(originalLink, "/")
	if len(groups) != 2 {
		return nil, fmt.Errorf("link has wrong number of items")
	}

	var vehicle common.VehicleType

	switch groups[0] {
	case "autobus":
		vehicle = common.Bus
	case "tramway":
		vehicle = common.Tram
	case "trolleybus":
		vehicle = common.Trolley
	case "metro":
		vehicle = common.Subway
	default:
		return nil, fmt.Errorf("unknown vehicle type: %s", groups[0])
	}

	return &common.Line{
		Vehicle: vehicle,
		Number:  groups[1],
	}, nil
}

func stopNamesToStops(stopNames map[int]string) []*common.Stop {
	stops := make([]*common.Stop, len(stopNames))
	i := 0
	for id, name := range stopNames {
		stops[i] = &common.Stop{
			ID:   id,
			Name: name,
		}
		i++
	}
	return stops
}
