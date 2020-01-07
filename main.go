package main

import (
	"log"
	"os"

	"github.com/alec-rabold/copypasta/commands"
	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("copypasta", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"": func() (cli.Command, error) {
			return &commands.CopyPasteCommand{Ui: uiColored}, nil
		},

		"foo": fooCommandFactory,
		"bar": barCommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
