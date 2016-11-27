package schedules

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/skgt_api/common"
	"github.com/jbowtie/gokogiri/xml"
)

// Time represents a schedule time (in the form of hours and minutes)
type Time struct {
	Hours   int
	Minutes int
}

// Scan parses Time from a database value (number of minutes since midnight)
func (t *Time) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	switch minutes := src.(type) {
	case int64:
		if minutes > 60*24 {
			return fmt.Errorf("%d minutes is more than a day", minutes)
		}
		t.Hours = int(minutes / 60)
		t.Minutes = int(minutes % 60)
		return nil
	default:
		return fmt.Errorf("unknown type for 'time': %T", src)
	}
}

// Value stores the Time in a database as number of minutes since midnight
func (t *Time) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}

	return int64(t.Hours*60 + t.Minutes), nil
}

// NewTime initialises Time from hours and minutes
func NewTime(hours int, minutes int) *Time {
	return &Time{
		Hours:   hours,
		Minutes: minutes,
	}
}

// Course is a path a single vehicle takes. It contains the times at which
// the vehicle stops at each of the stops in the route.
type Course []*Time

// ScheduleType represents the day type
type ScheduleType int

const (
	// None is an unknown day type
	None ScheduleType = 0
	// Workday is usualy monday-friday
	Workday = 1
	// Holiday is any national holiday + all sundays
	Holiday = 2
	// PreHoliday are all days scheduled as free around national holidays + all saturdays
	PreHoliday = 4
	// HolidayAndPreHoliday is a combination of Holiday and Preholiday
	HolidayAndPreHoliday = 6
	// All is a combination of all day types
	All = 7
)

// Route is a route which can be performed by a vehicle. Most vehicles have
// two routes - forward and backward
type Route struct {
	// Direction is a human-readable string which describes the route
	// (usually the endpoints)
	Direction string
	// Stops are the stops the vehicle stops at while following this route
	Stops []int
	// Schedules contains lists of all courses for each day
	Schedules map[ScheduleType][]Course
}

// StopName pairs a stop with its name
type StopName struct {
	ID   int
	Name string
}

// Timetable pairs a Line with all of its routes
type Timetable struct {
	Line   *common.Line
	Routes []*Route
}

// routeData is the data for a route we get by parsing a route page
type routeData struct {
	Direction string
	Stops     []int
	Courses   []Course
}

// GetTimetable gets the timetable for a given line, sending
// all stop names it discovers down the stopNames channel
// (there may be duplicate stops)
func GetTimetable(
	settings *htmlparsing.Settings,
	line *common.Line,
	stopNames chan<- *StopName,
) (*Timetable, error) {
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
			data, err := parseSchedule(directionDivs[j], stopNames)
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
				return nil, fmt.Errorf(
					"stops for same route are different on different days (%s %s), directions: (%s, %s), stops: %v, %v",
					line.Vehicle, line.Number,
					route.Direction, data.Direction,
					route.Stops, data.Stops,
				)
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

// sameStops checks if two slices of stop IDs are the same
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

// parseType parses a schedule type from a human-readable string
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

// parseSchedule parses a single schedule, sending all stops it finds
// down the stopNames channel
func parseSchedule(scheduleDiv xml.Node, stopNames chan<- *StopName) (*routeData, error) {
	direction, err := htmlparsing.First(scheduleDiv, `.//h6`)
	if err != nil {
		return nil, fmt.Errorf("unable to find direction: %s", err)
	}

	stops, err := parseStops(scheduleDiv, stopNames)
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

// parseStops parses a stop table, returning a slice of stop IDs and sending
// all stops it sees down the stopNames channel
func parseStops(scheduleDiv xml.Node, stopNames chan<- *StopName) ([]int, error) {
	stopItems, err := scheduleDiv.Search(
		`.//ul[contains(@class, 'schedule_direction_signs')]/li`,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to find stops: %s", err)
	}

	stops := make([]int, len(stopItems))
	for i := range stopItems {
		id, name, err := parseStop(stopItems[i])
		if err != nil {
			return nil, fmt.Errorf("unable to parse stop list item: %s", err)
		}

		stops[i] = id
		if stopNames != nil {
			stopNames <- &StopName{
				ID:   id,
				Name: name,
			}
		}
	}

	return stops, nil
}

func parseStop(stopItem xml.Node) (int, string, error) {
	idLink, err := htmlparsing.First(
		stopItem,
		`.//a[contains(@class, 'stop_link')]`,
	)
	if err != nil {
		return 0, "", fmt.Errorf("unable to find stop id link: %s", err)
	}

	id, err := strconv.Atoi(strings.TrimSpace(idLink.Content()))
	if err != nil {
		return 0, "", fmt.Errorf("unable to parse stop id: %s", err)
	}

	nameLink, err := htmlparsing.First(
		stopItem,
		`.//a[contains(@class, 'stop_change')]`,
	)
	if err != nil {
		return 0, "", fmt.Errorf("unable to find stop name link: %s", err)
	}

	name := strings.TrimSpace(nameLink.Content())

	return id, name, nil
}

func parseCourses(scheduleDiv xml.Node) ([]Course, error) {
	timeCells, err := scheduleDiv.Search(
		`.//div[contains(@class, 'hours_cell')]/a`,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to find timetable cells: %s", err)
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

// MarshalJSON implements the json.Marshaler interface
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
	case common.Subway:
		return "metro"
	default:
		return ""
	}
}
