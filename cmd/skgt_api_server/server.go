package main

import (
	"log"
	"net/http"

	"github.com/DexterLB/skgt_api/server"
	"github.com/urfave/cli"
)

func runServer(c *cli.Context) error {
	config, err := parseConfig(c)
	if err != nil {
		return err
	}
	backend, err := initBackend(config)
	if err != nil {
		return err
	}

	server := server.New(backend)

	log.Printf("starting HTTP server on address %s", config.Server.ListenAddress)
	log.Printf("exit: %s", http.ListenAndServe(config.Server.ListenAddress, server))

	return nil
}
