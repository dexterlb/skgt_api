package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DexterLB/skgt_api/backend"
	"github.com/DexterLB/skgt_api/server"
	"github.com/urfave/cli"
)

func runServer(c *cli.Context) error {
	config := parseConfig(c)

	log.Printf("initialising backend and connecting to database")
	backend, err := backend.New(config.Database.URN())
	log.Printf("finished backend initialisation")

	if err != nil {
		return fmt.Errorf("unable to initialise backend: %s", err)
	}

	server := server.New(backend)

	log.Printf("starting HTTP server on address %s", config.Server.ListenAddress)
	log.Printf("exit: %s", http.ListenAndServe(config.Server.ListenAddress, server))

	return nil
}
