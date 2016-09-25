package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DexterLB/skgt_api/backend"
	"github.com/DexterLB/skgt_api/config"
	"github.com/cep21/xdgbasedir"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "skgt api server"
	app.Usage = "totally legal."

	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "initialise an empty database (deletes any old data)",
			Action: runInit,
			Flags:  []cli.Flag{},
		},
		{
			Name:   "update",
			Usage:  "update the database with data parsed from the site",
			Action: runUpdate,
			Flags:  []cli.Flag{},
		},
		{
			Name:   "serve",
			Usage:  "start a http server with the API",
			Action: runServer,
			Flags:  []cli.Flag{},
		},
		{
			Name:   "apikey",
			Usage:  "operate on API keys stored in the database",
			Action: runAPIKey,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "check, c",
					Usage: "check if an API key is valid",
				},
				cli.BoolFlag{
					Name:  "new, n",
					Usage: "create a new API key and print it to stdout",
				},
				cli.StringFlag{
					Name:  "delete, d",
					Usage: "delete an API key",
				},
			},
		},
	}

	defaultConfigFile, _ := xdgbasedir.GetConfigFileLocation("skgt.toml")

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config-file, c",
			Usage: "location of the config file",
			Value: defaultConfigFile,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Printf("error: %s", err)
	}
}

func initBackend(config *config.Config) (*backend.Backend, error) {
	log.Printf("initialising backend and connecting to database")
	backend, err := backend.New(config.Database.URN())
	log.Printf("finished backend initialisation")

	if err != nil {
		return nil, fmt.Errorf("unable to initialise backend: %s", err)
	}

	return backend, nil
}

func parseConfig(c *cli.Context) (*config.Config, error) {
	var err error

	filename := c.GlobalString("config-file")

	config, err := config.Load(filename)
	if err != nil {
		return nil, fmt.Errorf("can't load config file: %s", err)
	}

	return config, nil
}
