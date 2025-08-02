package main

import (
	"fmt"
	"github.com/nhdms/base-go/pkg/toolkit/cmd"
	"github.com/nhdms/base-go/pkg/toolkit/generator"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "atcli",
		Usage: "A tool to generate various components",
		Commands: []*cli.Command{
			{
				Name:  "consul",
				Usage: "Clone services from a remote Consul server and register them in another",
				Subcommands: []*cli.Command{
					{
						Name:  "clone",
						Usage: "Clone Consul data from source to destination",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "source",
								Usage:    "Source Consul address",
								EnvVars:  []string{"CONSUL_SOURCE_ADDR"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "dest",
								Usage:    "Destination Consul address",
								EnvVars:  []string{"CONSUL_DEST_ADDR"},
								Required: true,
							},
							&cli.BoolFlag{
								Name:  "services",
								Usage: "Clone services",
							},
							&cli.BoolFlag{
								Name:  "kv",
								Usage: "Clone KV store",
							},
						},
						Action: cmd.CloneCommand,
					},
					{
						Name:   "deregister",
						Usage:  "Deregister all services from Consul",
						Action: cmd.DeregisterCommand,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Usage:    "service name",
								Required: false,
							},
						},
					},
				},
			},
			{
				Name:  "gen",
				Usage: "Generate various components",
				Subcommands: []*cli.Command{
					{
						Name:  "proto",
						Usage: "Generate proto file",
						Action: func(c *cli.Context) error {
							filePath := c.Args().First()
							if filePath == "" {
								return fmt.Errorf("name is required")
							}
							fmt.Printf("Generating proto file at: %s\n", filePath)
							generator.GenerateProtoFile(filePath)
							// Add logic to generate proto file here
							return nil
						},
					},
					{
						Name:  "cicd",
						Usage: "Generate service to build",
						Action: func(c *cli.Context) error {
							pathToService := c.Args().First()
							if pathToService == "" {
								return fmt.Errorf("path cicd.json or folder is required")
							}

							generator.GenerateCICD(pathToService)
							// Add logic to generate scheduler here
							return nil
						},
					},
					{
						Name:  "goland-config",
						Usage: "Generate Goland configuration",
						Action: func(c *cli.Context) error {
							// Add logic to generate proto file here
							generator.ExportGolandConfig(nil)
							return nil
						},
					},
					{
						Name:  "api",
						Usage: "Generate API service",
						Action: func(c *cli.Context) error {
							name := c.Args().First()
							if name == "" {
								return fmt.Errorf("name is required")
							}
							fmt.Printf("Generating API service: %s\n", name)
							return generator.GenerateAPI(name)
						},
					},
					{
						Name:  "service",
						Usage: "Generate gRPC service",
						Action: func(c *cli.Context) error {
							name := c.Args().First()
							if name == "" {
								return fmt.Errorf("name is required")
							}
							fmt.Printf("Generating gRPC service: %s\n", name)
							return generator.GenerateService(name)
						},
					},
					{
						Name:  "consumer",
						Usage: "Generate consumer",
						Action: func(c *cli.Context) error {
							name := c.Args().First()
							if name == "" {
								return fmt.Errorf("name is required")
							}
							fmt.Printf("Generating consumer: %s\n", name)
							return generator.GenerateConsumer(name)
						},
					},
					{
						Name:  "scheduler",
						Usage: "Generate scheduler",
						Action: func(c *cli.Context) error {
							name := c.Args().First()
							if name == "" {
								return fmt.Errorf("name is required")
							}
							fmt.Printf("Generating scheduler: %s\n", name)
							// Add logic to generate scheduler here
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
