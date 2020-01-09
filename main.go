package main

import (
	"log"
	"os"

	"github.com/alec-rabold/copypasta/commands"
	"github.com/mitchellh/cli"
)

func main() {
	ui := &cli.BasicUi{
		Writer:      os.Stdout,
		Reader:      os.Stdin,
		ErrorWriter: os.Stdout,
	}

	uiColored := &cli.ColoredUi{
		OutputColor: cli.UiColorNone,
		InfoColor:   cli.UiColorNone,
		ErrorColor:  cli.UiColorRed,
		Ui:          ui,
	}

	c := cli.NewCLI("copypasta", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"": func() (cli.Command, error) {
			return &commands.CopyPasteCommand{UI: uiColored}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
