package realtime

import (
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/skgt_api/common"
	"github.com/jbowtie/gokogiri/xml"
)

const pageURL = "https://skgt-bg.com/VirtualBoard/Web/SelectByStop.aspx"
const captchaURL = "https://skgt-bg.com/VirtualBoard/Services/Captcha.ashx"

var location, _ = time.LoadLocation("Europe/Sofia")

type StopData struct {
	Parameters    map[string]string
	Lines         map[common.Line]int
	Captcha       io.Reader
	CaptchaResult string
	client        *htmlparsing.Client
	Name          string
	Description   string
}

type Arrival struct {
	Time            time.Time
	Calculated      time.Time
	AirConditioning bool
	Accessibility   bool
}

func (s *StopData) Arrivals(lineID int) ([]*Arrival, error) {
	s.Parameters["ctl00$ContentPlaceHolder1$ddlLine"] = fmt.Sprintf("%d", lineID)
	s.Parameters["ctl00$ContentPlaceHolder1$CaptchaInput"] = s.CaptchaResult

	page, err := s.client.ParsePage(pageURL, urlValues(s.Parameters))
	if err != nil {
		return nil, fmt.Errorf("cannot get line info page: %s", err)
	}
	defer page.Free()

	s.Parameters, err = getFormValues(page)
	if err != nil {
		return nil, fmt.Errorf("cannot get POST parameters: %s", err)
	}

	rows, err := page.Search(
		`//table[contains(@id,"ctl00_ContentPlaceHolder1_gvTimes")]/tr[not(contains(@class, "Header"))]`,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot get arrivals table rows: %s", err)
	}

	arrivals := make([]*Arrival, len(rows))
	for i := range rows {
		arrivals[i], err = parseArrival(rows[i])
		if err != nil {
			return nil, fmt.Errorf("unable to parse arrival: %s", err)
		}
	}

	return arrivals, nil
}

func (s *StopData) LoadCaptcha() error {
	var err error
	s.Captcha, err = getCaptcha(s.client)
	return err
}

func (s *StopData) BreakCaptcha() error {
	err := s.LoadCaptcha()
	if err != nil {
		return fmt.Errorf("unable to load captcha: %s", err)
	}

	result, err := htmlparsing.BreakSimpleCaptcha(s.Captcha)
	if err != nil {
		return fmt.Errorf("unable to break captcha: %s", err)
	}

	s.CaptchaResult = result

	return nil
}

func LookupStop(settings *htmlparsing.Settings, id int) (*StopData, error) {
	client, err := htmlparsing.NewCookiedClient(settings)
	if err != nil {
		return nil, fmt.Errorf("unable to initialise http client: %s", err)
	}

	page, err := client.ParsePage(pageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot parse search page: %s", err)
	}
	defer page.Free()

	parameters, err := getFormValues(page)
	if err != nil {
		return nil, fmt.Errorf("unable to get hidden values: %s", err)
	}

	parameters["ctl00$ContentPlaceHolder1$tbStopCode"] = fmt.Sprintf("%04d", id)
	parameters["ctl00$ContentPlaceHolder1$btnSearchLine.x"] = fmt.Sprintf("%d", rand.Intn(53))
	parameters["ctl00$ContentPlaceHolder1$btnSearchLine.y"] = fmt.Sprintf("%d", rand.Intn(16))

	page, err = client.ParsePage(pageURL, urlValues(parameters))
	if err != nil {
		return nil, fmt.Errorf("cannot parse selection page: %s", err)
	}
	defer page.Free()

	htmlparsing.DumpHTML(page, "/tmp/bleh.html")

	stopName, err := htmlparsing.First(
		page,
		`//span[contains(@id, 'lblStopName')]`,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to get stop name: %s", err)
	}

	stopDescription, err := htmlparsing.First(
		page,
		`//span[contains(@id, 'lblDescription')]`,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to get stop description: %s", err)
	}

	parameters, err = getFormValues(page)
	if err != nil {
		return nil, fmt.Errorf("unable to get hidden values: %s", err)
	}

	lines, err := getLines(page)
	if err != nil {
		return nil, fmt.Errorf("unable to get lines: %s", err)
	}
	data := &StopData{
		client:      client,
		Parameters:  parameters,
		Lines:       lines,
		Name:        stopName.Content(),
		Description: strings.Trim(stopDescription.Content(), "\r\n\t \";"),
	}

	return data, nil
}

func parseArrival(row xml.Node) (*Arrival, error) {
	arrival := &Arrival{}

	accessibilityMarkers, err := row.Search(
		`.//img[contains(@id, "imgPlatform")]`,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to search for accessibility markers: %s", err)
	}

	arrival.Accessibility = (len(accessibilityMarkers) > 0)

	airConditioningMarkers, err := row.Search(
		`.//img[contains(@id, "imgAirCondition")]`,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to search for air conditioning markers: %s", err)
	}

	arrival.AirConditioning = (len(airConditioningMarkers) > 0)

	timeString, err := htmlparsing.First(
		row,
		`.//div[contains(@id, "dvItem")]`,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to find arrival time: %s", err)
	}

	arrival.Time, arrival.Calculated, err = parseArrivalTime(timeString.Content())
	if err != nil {
		return nil, fmt.Errorf("unable to parse arrival time: %s", err)
	}

	return arrival, nil
}

func parseArrivalTime(input string) (time.Time, time.Time, error) {
	groups := regexp.MustCompile(
		`([\d]+)\:([\d]+) изчислено в. ([\d]+\:[\d]+ [\d]+\.[\d]+\.[\d]+)`,
	).FindStringSubmatch(input)
	if len(groups) < 4 {
		return time.Time{}, time.Time{}, fmt.Errorf("unable to find time")
	}

	calculated, err := time.ParseInLocation("15:04 02.01.2006", groups[3], location)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("unable to parse calculated time: %s", err)
	}

	hour, err := strconv.Atoi(groups[1])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("unable to parse arrival time hour: %s", err)
	}

	minute, err := strconv.Atoi(groups[2])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("unable to parse arrival time minute: %s", err)
	}

	year, month, day := calculated.Date()

	arrival := time.Date(year, month, day, hour, minute, 0, 0, location)
	if arrival.Before(calculated) {
		day += 1
		arrival = time.Date(year, month, day, hour, minute, 0, 0, location)
	}

	return arrival, calculated, nil
}

func getCaptcha(client *htmlparsing.Client) (io.Reader, error) {
	response, err := client.Get(captchaURL)
	if err != nil {
		return nil, fmt.Errorf("unable to get captcha: %s", err)
	}

	return response.Body, nil
}

func getLines(page xml.Node) (map[common.Line]int, error) {
	options, err := page.Search(
		`//select/option[@value != ""]`,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to find selector options: %s", err)
	}

	lines := make(map[common.Line]int)

	for i := range options {
		value, ok := options[i].Attributes()["value"]
		if !ok {
			return nil, fmt.Errorf("option element has no value")
		}

		id, err := strconv.Atoi(value.Value())
		if err != nil {
			return nil, fmt.Errorf("option value is not integer: %s", err)
		}

		line, err := common.ParseLine(options[i].Content())
		if err != nil {
			return nil, fmt.Errorf("unable to parse line: %s", err)
		}

		lines[*line] = id
	}

	return lines, nil
}

func getFormValues(page xml.Node) (map[string]string, error) {
	hiddenInputs, err := page.Search(
		`//input`,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to find hidden input elements: %s", err)
	}

	hiddenValues := make(map[string]string)

	for i := range hiddenInputs {
		attributes := hiddenInputs[i].Attributes()
		key, ok := attributes["name"]
		if !ok {
			return nil, fmt.Errorf("input element has no name")
		}

		value, ok := attributes["value"]
		if !ok {
			hiddenValues[key.Value()] = ""
		} else {
			hiddenValues[key.Value()] = value.Value()
		}
	}

	return hiddenValues, nil
}

func urlValues(parameters map[string]string) url.Values {
	values := make(url.Values)

	for key := range parameters {
		values.Set(key, parameters[key])
	}

	return values
}
