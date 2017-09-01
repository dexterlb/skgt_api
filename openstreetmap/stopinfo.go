package openstreetmap

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/DexterLB/htmlparsing"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

var OverpassURL = "http://overpass-api.de/api/interpreter"
var StopsQuery = `
<osm-script>

<union>
    <query type="relation">
        <has-kv k="name" v="автобуси в София"/>
    </query>
    <recurse type="relation-relation" />
    <recurse type="relation-node" />
    
    <query type="relation">
        <has-kv k="name" v="Тролеи в София"/>
    </query>
    <recurse type="relation-relation" />
    <recurse type="relation-node" />
    
    <query type="relation">
        <has-kv k="name" v="Трамваи в София"/>
    </query>
    <recurse type="relation-relation" />
    <recurse type="relation-node" />

    <query type="relation">
        <has-kv k="name" v="метро в София"/>
    </query>
    <recurse type="relation-relation" />
    <recurse type="relation-node" />
</union>

<print />

<osm-script>
`

// Stop represents a bus stop as in the OpenStreetMap Overpass API
type Stop struct {
	ID                int
	Name              string
	InternationalName string
	Latitude          float64
	Longitude         float64
}

// GetStops gets all stops
func GetStops(settings *htmlparsing.Settings) (map[int]*Stop, error) {
	data, err := rawData(settings)
	if err != nil {
		return nil, err
	}

	return parse(data)
}

func rawData(settings *htmlparsing.Settings) ([]byte, error) {
	client := htmlparsing.NewClient(settings)

	resp, err := client.Post(
		OverpassURL,
		"text/xml",
		strings.NewReader(StopsQuery),
	)

	if err != nil {
		return nil, fmt.Errorf("HTTP error: %s", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read page: %s", err)
	}

	return data, nil
}

func parse(data []byte) (map[int]*Stop, error) {
	doc, err := gokogiri.ParseXml(data)
	if err != nil {
		return nil, fmt.Errorf("unable to parse XML: %s", err)
	}

	nodes, err := doc.Root().Search(
		`//node[tag/@k = 'ref']`,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to find nodes: %s", err)
	}

	stops := make(map[int]*Stop)

	for i := range nodes {
		stop, err := parseStop(nodes[i])
		if err != nil {
			return nil, fmt.Errorf("unable to parse node: %s", err)
		}
		stops[stop.ID] = stop
	}

	return stops, nil
}

func parseStop(node xml.Node) (*Stop, error) {
	stop := &Stop{}

	tags, err := tags(node)
	if err != nil {
		return nil, err
	}

	attributes := node.Attributes()

	if ref, ok := tags["ref"]; ok {
		stop.ID, err = strconv.Atoi(ref)
		if err != nil {
			return nil, fmt.Errorf("ref is not a number: %s", err)
		}
	} else {
		return nil, fmt.Errorf("no ref tag")
	}

	if name, ok := tags["name"]; ok {
		stop.Name = name
	}

	if intName, ok := tags["name:en"]; ok {
		stop.InternationalName = intName
	}

	if intName, ok := tags["int_name"]; ok {
		stop.InternationalName = intName
	}

	lat, ok := attributes["lat"]
	stop.Latitude, err = strconv.ParseFloat(lat.Value(), 64)
	if !ok || err != nil {
		return nil, fmt.Errorf("node has no latitude")
	}

	lon, ok := attributes["lon"]
	stop.Longitude, err = strconv.ParseFloat(lon.Value(), 64)
	if !ok || err != nil {
		return nil, fmt.Errorf("node has no longitude")
	}

	return stop, nil
}

func tags(node xml.Node) (map[string]string, error) {
	tagNodes, err := node.Search(`.//tag`)
	if err != nil {
		return nil, fmt.Errorf("unable to find tags: %s", err)
	}

	tags := make(map[string]string)
	for i := range tagNodes {
		attributes := tagNodes[i].Attributes()

		key, ok := attributes["k"]
		if !ok {
			return nil, fmt.Errorf("Tag with no key: %s", tagNodes[i])
		}

		value, ok := attributes["v"]
		if !ok {
			return nil, fmt.Errorf("Tag with no value: %s", tagNodes[i])
		}

		tags[key.Value()] = value.Value()
	}

	return tags, nil
}
