package main

import (
	"fmt"
	"log"

	"github.com/DexterLB/skgt_api/backend"
	"github.com/urfave/cli"
)

func runAPIKey(c *cli.Context) error {
	config := parseConfig(c)

	log.Printf("initialising backend and connecting to database")
	backend, err := backend.New(config.Database.URN())
	log.Printf("finished backend initialisation")

	if err != nil {
		return fmt.Errorf("unable to initialise backend: %s", err)
	}

	apiKey := c.String("check")
	if apiKey != "" {
		err := backend.CheckAPIKey(apiKey)
		if err != nil {
			return err
		}
		fmt.Printf("API key is valid.\n")
	}

	apiKey = c.String("delete")
	if apiKey != "" {
		err := backend.CheckAPIKey(apiKey)
		if err != nil {
			return err
		}
		err = backend.DeleteAPIKey(apiKey)
		if err != nil {
			return err
		}
		fmt.Printf("deleted API key.\n")
	}

	if c.Bool("new") {
		apiKey, err := backend.NewAPIKey()
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", apiKey)
	}

	return nil
}
