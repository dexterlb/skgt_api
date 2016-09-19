package schedules

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/skgt_api/common"
)

func AllTimetables(settings *htmlparsing.Settings) ([]*Timetable, error) {
	lines, err := AllLines(settings)
	if err != nil {
		return nil, fmt.Errorf("unable to get list of lines")
	}

	infos := make([]*Timetable, len(lines))
	for i := range lines {
		infos[i], err = GetTimetable(settings, lines[i])
		if err != nil {
			return nil, fmt.Errorf("unable to get schedule info: %s", err)
		}
	}

	return infos, nil
}

func AllLines(settings *htmlparsing.Settings) ([]*common.Line, error) {
	page, err := htmlparsing.NewClient(settings).ParsePage(
		`https://schedules.sofiatraffic.bg/`, nil,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to parse line list page: %s", err)
	}

	links, err := page.Search(
		`//div[contains(@class, 'lines_section')]/ul/li/a`,
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

func GetStops(infos []*Timetable) []int {
	stopSet := make(map[int]struct{})
	for _, info := range infos {
		for _, route := range info.Routes {
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
	default:
		return nil, fmt.Errorf("unknown vehicle type: %s", groups[0])
	}

	return &common.Line{
		Vehicle: vehicle,
		Number:  groups[1],
	}, nil
}
