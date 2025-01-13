package main

import (
	"fmt"
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
						Name:  "api",
						Usage: "Generate API service",
						Action: func(c *cli.Context) error {
							name := c.Args().First()
							if name == "" {
								return fmt.Errorf("name is required")
							}
							fmt.Printf("Generating API service: %s\n", name)
							// Add logic to generate API service here
							return nil
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
							// Add logic to generate gRPC service here
							return nil
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
							// Add logic to generate consumer here
							return nil
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
