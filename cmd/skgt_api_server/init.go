package main

import (
	"fmt"
	"log"

	"github.com/urfave/cli"
)

func runInit(c *cli.Context) error {
	config, err := parseConfig(c)
	if err != nil {
		return err
	}
	backend, err := initBackend(config)
	if err != nil {
		return err
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
