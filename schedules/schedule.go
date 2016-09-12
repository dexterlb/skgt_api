package schedules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/skgt_api/common"
	"github.com/jbowtie/gokogiri/xml"
)

type Time [2]int

type Course []Time

type Schedule struct {
	Type      string
	Direction string
	Stops     []int
	Courses   []Course
}

type ScheduleInfo struct {
	Line      *common.Line
	Schedules []*Schedule
}

func GetScheduleInfo(settings *htmlparsing.Settings, line *common.Line) (*ScheduleInfo, error) {
	page, err := htmlparsing.NewClient(settings).ParsePage(
		fmt.Sprintf(
			"https://schedules.sofiatraffic.bg/%s/%s",
			transportType(line.Type),
			line.Number,
		),
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to parse schedule page: %s", err)
	}

	htmlparsing.DumpHTML(page, "/tmp/bleh.html")

	typeDivs, err := page.Search(
		`//div[contains(@class, 'schedule_active_list_content')]`,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to find schedule type divs: %s", err)
	}

	var schedules []*Schedule

	for i := range typeDivs {
		scheduleType, err := htmlparsing.First(typeDivs[i], `.//h3`)
		if err != nil {
			return nil, fmt.Errorf("unable to find schedule type: %s", err)
		}

		directionDivs, err := typeDivs[i].Search(
			`.//div[contains(@class, 'schedule_view_direction_content')]`,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to find directions: %s", err)
		}

		for j := range directionDivs {
			schedule, err := parseSchedule(directionDivs[j])
			if err != nil {
				return nil, fmt.Errorf("unable to parse schedule: %s", err)
			}
			schedule.Type = strings.TrimSpace(scheduleType.Content())

			schedules = append(schedules, schedule)
		}
	}

	return &ScheduleInfo{
		Schedules: schedules,
		Line:      line,
	}, nil
}

func parseSchedule(scheduleDiv xml.Node) (*Schedule, error) {
	direction, err := htmlparsing.First(scheduleDiv, `.//h6`)
	if err != nil {
		return nil, fmt.Errorf("unable to find direction: %s")
	}

	stops, err := parseStops(scheduleDiv)
	if err != nil {
		return nil, fmt.Errorf("unable to parse stops: %s", err)
	}

	courses, err := parseCourses(scheduleDiv)
	if err != nil {
		return nil, fmt.Errorf("unable to parse courses: %s", err)
	}

	return &Schedule{
		Direction: strings.TrimSpace(direction.Content()),
		Courses:   courses,
		Stops:     stops,
	}, nil
}

func parseStops(scheduleDiv xml.Node) ([]int, error) {
	stopItems, err := scheduleDiv.Search(
		`.//ul[contains(@class, 'schedule_direction_signs')]/li/a[contains(@class, 'stop_link')]`,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to find stops: %s", err)
	}

	stops := make([]int, len(stopItems))
	for i := range stopItems {
		stops[i], err = strconv.Atoi(strings.TrimSpace(stopItems[i].Content()))
		if err != nil {
			return nil, fmt.Errorf("unable to parse stop id: %s", err)
		}
	}

	return stops, nil
}

func parseCourses(scheduleDiv xml.Node) ([]Course, error) {
	return nil, nil
}

func transportType(transport common.Transport) string {
	switch transport {
	case common.Bus:
		return "autobus"
	case common.Tram:
		return "tramway"
	case common.Trolley:
		return "trolleybus"
	default:
		return ""
	}
}
