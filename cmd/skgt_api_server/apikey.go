package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func runAPIKey(c *cli.Context) error {
	config, err := parseConfig(c)
	if err != nil {
		return err
	}
	backend, err := initBackend(config)
	if err != nil {
		return err
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
