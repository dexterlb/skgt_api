package realtime

import (
	"fmt"
	"math/rand"
	"net/url"

	"github.com/DexterLB/htmlparsing"
	"github.com/jbowtie/gokogiri/xml"
)

const pageURL = "https://skgt-bg.com/VirtualBoard/Web/SelectByStop.aspx"

type StopData struct {
	Parameters map[string]string
	client     *htmlparsing.Client
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

	parameters, err = getFormValues(page)
	if err != nil {
		return nil, fmt.Errorf("unable to get hidden values: %s", err)
	}

	htmlparsing.DumpHTML(page, "/tmp/bleh.html")

	data := &StopData{
		client:     client,
		Parameters: parameters,
	}

	return data, fmt.Errorf("not implemented")
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

	fmt.Printf("values: %v\n", values)

	return values
}
