package main

import (
	"fmt"
	"log"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/skgt_api/backend"
	"github.com/DexterLB/skgt_api/realtime"
	"github.com/DexterLB/skgt_api/schedules"
	"github.com/codegangsta/cli"
)

func runUpdate(c *cli.Context) error {
	config := parseConfig(c)

	log.Printf("initialising backend and connecting to database")
	backend, err := backend.New(config.Database.URN())
	log.Printf("finished backend initialisation")

	if err != nil {
		return fmt.Errorf("unable to initialise backend: %s", err)
	}

	log.Printf("parsing timetables")
	timetables, err := schedules.AllTimetables(htmlparsing.SensibleSettings())
	log.Printf("finished parsing timetables")

	if err != nil {
		return fmt.Errorf("unable to get schedules: %s", err)
	}

	log.Printf("parsing stop info")
	stopInfos, err := realtime.GetStopsInfo(
		htmlparsing.SensibleSettings(),
		schedules.GetStops(timetables),
		config.Parser.ParallelRequests,
	)
	log.Printf("finished parsing stop info")

	if err != nil {
		return fmt.Errorf("unable to get stops: %s", err)
	}

	log.Printf("depositing data to database")
	err = backend.Fill(stopInfos, timetables)
	log.Printf("finished depositing data to database")

	if err != nil {
		return fmt.Errorf("unable to write data to database: %s", err)
	}

	return nil
}
