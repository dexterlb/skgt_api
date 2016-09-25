package main

import (
	"fmt"
	"log"

	"github.com/DexterLB/skgt_api/backend"
	"github.com/codegangsta/cli"
)

func runInit(c *cli.Context) error {
	config := parseConfig(c)

	log.Printf("initialising backend and connecting to database")
	backend, err := backend.New(config.Database.URN())
	log.Printf("finished backend initialisation")

	if err != nil {
		return fmt.Errorf("unable to initialise backend: %s", err)
	}

	log.Printf("dropping old database")
	err = backend.DropDB()
	if err != nil {
		log.Printf("can't drop database: %s (running for the first time?)", err)
	} else {
		log.Printf("dropped old database")
	}

	log.Printf("creating tables")
	err = backend.InitDB()
	log.Printf("finished creating tables")

	if err != nil {
		return fmt.Errorf("unable to initialise database: %s", err)
	}

	return nil
}
