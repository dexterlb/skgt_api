package main

import (
	"fmt"
	"log"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/skgt_api/realtime"
	"github.com/DexterLB/skgt_api/schedules"
	"github.com/urfave/cli"
)

func runUpdate(c *cli.Context) error {
	config, err := parseConfig(c)
	if err != nil {
		return err
	}
	backend, err := initBackend(config)
	if err != nil {
		return err
	}

	log.Printf("parsing timetables")
	timetables, stopInfos, err := schedules.AllTimetables(
		htmlparsing.SensibleSettings(),
		config.Parser.ParallelRequests,
	)
	log.Printf("finished parsing timetables")

	if err != nil {
		return fmt.Errorf("unable to get schedules: %s", err)
	}

	log.Printf("parsing stop info")
	err = realtime.UpdateStopsInfo(
		htmlparsing.SensibleSettings(),
		stopInfos,
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
