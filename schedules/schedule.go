package schedules

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/skgt_api/common"
	"github.com/jbowtie/gokogiri/xml"
)

type Time struct {
	Hours   int
	Minutes int
}

func NewTime(hours int, minutes int) *Time {
	return &Time{
		Hours:   hours,
		Minutes: minutes,
	}
}

type Course []*Time

type ScheduleType int

const (
	None                 ScheduleType = 0
	Workday                           = 1
	Holiday                           = 2
	PreHoliday                        = 4
	HolidayAndPreHoliday              = 6
	All                               = 7
)

type Route struct {
	Direction string
	Stops     []int
	Schedules map[ScheduleType][]Course
}

type Timetable struct {
	Line   *common.Line
	Routes []*Route
}

type routeData struct {
	Direction string
	Stops     []int
	Courses   []Course
}

func GetTimetable(settings *htmlparsing.Settings, line *common.Line) (*Timetable, error) {
	page, err := htmlparsing.NewClient(settings).ParsePage(
		fmt.Sprintf(
			"https://schedules.sofiatraffic.bg/%s/%s",
			vehicle(line.Vehicle),
			line.Number,
		),
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to parse schedule page: %s", err)
	}

	typeDivs, err := page.Search(
		`//div[contains(@class, 'schedule_active_list_content')]`,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to find schedule type divs: %s", err)
	}

	routes := make(map[string]*Route)

	for i := range typeDivs {
		scheduleTypeHeader, err := htmlparsing.First(typeDivs[i], `.//h3`)
		if err != nil {
			return nil, fmt.Errorf("unable to find schedule type: %s", err)
		}
		scheduleType, err := parseType(scheduleTypeHeader.Content())
		if err != nil {
			return nil, fmt.Errorf("unable to parse schedule type: %s", err)
		}

		directionDivs, err := typeDivs[i].Search(
			`.//div[contains(@class, 'schedule_view_direction_content')]`,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to find directions: %s", err)
		}

		for j := range directionDivs {
			data, err := parseSchedule(directionDivs[j])
			if err != nil {
				return nil, fmt.Errorf("unable to parse schedule: %s", err)
			}

			route, ok := routes[data.Direction]
			if !ok {
				route = &Route{
					Direction: data.Direction,
					Stops:     data.Stops,
					Schedules: make(map[ScheduleType][]Course),
				}
				routes[data.Direction] = route
			}

			if !sameStops(route.Stops, data.Stops) {
				return nil, fmt.Errorf("stops for same route are different on different days")
			}

			route.Schedules[scheduleType] = data.Courses
		}
	}

	routeSlice := make([]*Route, len(routes))
	i := 0
	for _, route := range routes {
		routeSlice[i] = route
		i++
	}

	return &Timetable{
		Line:   line,
		Routes: routeSlice,
	}, nil
}

func sameStops(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func parseType(typeName string) (ScheduleType, error) {
	switch strings.TrimSpace(strings.Split(typeName, "-")[0]) {
	case "делник":
		return Workday, nil
	case "предпразник, празник":
		return HolidayAndPreHoliday, nil
	case "предпразник":
		return PreHoliday, nil
	case "празник":
		return Holiday, nil
	case "делник, предпразник, празник":
		return All, nil
	default:
		return None, fmt.Errorf("unknown schedule type: %s", typeName)
	}
}

func parseSchedule(scheduleDiv xml.Node) (*routeData, error) {
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

	return &routeData{
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
	timeCells, err := scheduleDiv.Search(
		`.//div[contains(@class, 'hours_cell')]/a`,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to find timetable cells", err)
	}

	courses := make([]Course, len(timeCells))
	for i := range timeCells {
		jsLink, ok := timeCells[i].Attributes()["onclick"]
		if !ok {
			return nil, fmt.Errorf("time link has no onclick attribute")
		}

		courses[i], err = parseCourse(jsLink.Value())
		if err != nil {
			return nil, fmt.Errorf("unable to parse course: %s", err)
		}
	}

	return courses, nil
}

func parseCourse(jsCall string) (Course, error) {
	// we shall now parse javascript code.
	// here be dragons.

	// a sample jsCall:
	// Raz.exec ('show_course', ['9caca5ad9', '4,761,763,765,766,767,768,770,772,775,777,780,783,785,788,792,794,796,798,800,802,804,806,807']); return false;

	// we need the '4,761,...' part.
	groups := regexp.MustCompile(
		`.+\[\'\w+\', \'([0-9, ]+)\'\].*`,
	).FindStringSubmatch(jsCall)
	if len(groups) < 2 {
		return nil, fmt.Errorf("unable to find time list in this javascript code: %s", jsCall)
	}

	times := strings.Split(groups[1], ",")[1:]

	course := make(Course, len(times))
	var err error
	for i := range times {
		course[i], err = parseTime(times[i])
		if err != nil {
			return nil, fmt.Errorf("unable to parse time: %s", err)
		}
	}

	return course, nil
}

func parseTime(skgtTime string) (*Time, error) {
	if len(skgtTime) == 0 {
		return nil, nil // empty string is valid - transport doesn't reach stop
	}

	// skgtTime is actually an integer: number of minutes since 00:00
	minutes, err := strconv.Atoi(skgtTime)
	if err != nil {
		return nil, err
	}

	return &Time{
		Hours:   minutes / 60,
		Minutes: minutes % 60,
	}, nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	if t == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(fmt.Sprintf("%02d:%02d", t.Hours, t.Minutes))
}

func vehicle(transport common.VehicleType) string {
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
