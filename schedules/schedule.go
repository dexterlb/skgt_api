package schedules

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
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

type Course []*Time

type ScheduleType int

const (
	None                              = 0
	Workday              ScheduleType = 1
	Holiday                           = 2
	PreHoliday                        = 4
	HolidayAndPreHoliday              = 6
	All                               = 7
)

type Schedule struct {
	Type      ScheduleType
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
			schedule.Type, err = parseType(scheduleType.Content())
			if err != nil {
				return nil, fmt.Errorf("unable to parse schedule type: %s", err)
			}

			schedules = append(schedules, schedule)
		}
	}

	log.Printf("checking")
	for i := range schedules {
		for j := range schedules {
			if schedules[i].Direction == schedules[j].Direction {
				if !reflect.DeepEqual(schedules[i].Stops, schedules[j].Stops) {
					return nil, fmt.Errorf(":(")
				}
			}
		}
	}

	return &ScheduleInfo{
		Schedules: schedules,
		Line:      line,
	}, nil
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
