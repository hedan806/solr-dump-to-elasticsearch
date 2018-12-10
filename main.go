package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "greet"
	app.Usage = "fight the loneliness!"
	app.Action = func(c *cli.Context) error {
		name := "Nefertiti"
		if c.String("lang") == "spanish" {
			fmt.Println("Hola", name)
		} else {
			fmt.Println("Hello", name)
		}
		return nil
	}

	flags := []cli.Flag{
		cli.StringFlag{Name: "load"},
	}

	app.Commands = []cli.Command{
		{
			Name:    "dump",
			Aliases: []string{"d"},
			Usage:   "dump solr data source",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "field, f"},
				cli.StringFlag{Name: "index, i"},
				cli.StringFlag{Name: "src, s"},
				cli.StringFlag{Name: "target, t"},
				cli.IntFlag{Name: "consumer, c"},
			},
			Action: func(c *cli.Context) error {
				src := c.String("src")
				target := c.String("target")
				index := c.String("index")
				field := c.String("field")
				consumerCount := c.Int("consumer")
				//Export(server, true, field)
				Export2ES(consumerCount, src,
					target, index, "doc", field)
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a task to the list",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	app.Flags = flags

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
