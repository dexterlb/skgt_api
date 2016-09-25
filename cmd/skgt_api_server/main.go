package main

import (
	"log"
	"os"

	"github.com/DexterLB/skgt_api/config"
	"github.com/cep21/xdgbasedir"
	"github.com/codegangsta/cli"
)

func parseConfig(c *cli.Context) *config.Config {
	var err error

	filename := c.GlobalString("config-file")

	config, err := config.Load(filename)
	if err != nil {
		log.Fatalf("can't load config file: %s", err)
	}

	return config
}

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
