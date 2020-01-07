package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/alec-rabold/copypasta/runcommands"
	"github.com/alec-rabold/copypasta/store"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

// CopyPasteCommand is the command for running the copying and pasting procedures
type CopyPasteCommand struct {
	UI cli.Ui
}

// Run function for the tool
func (c *CopyPasteCommand) Run(args []string) int {
	target, err := runcommands.Load()
	if err != nil && target != nil {
		c.UI.Error(fmt.Sprintf("Failed to load the runcommands file: %s", err.Error()))
		os.Exit(1)
	}

	copyPasteCommand := flag.NewFlagSet("", flag.ExitOnError)
	optionPaste := copyPasteCommand.Bool("paste", false, "")

	if err := copyPasteCommand.Parse(args); err != nil {
		os.Exit(1)
	}

	store, err := store.NewStore(target)
	if err != nil {
		log.Errorf("Error creating minio s3 client %s", err.Error())
		os.Exit(1)
	}
	if *optionPaste {
		doPaste(target, store)
	} else {
		doCopy(target, store)
	}
	return 0
}

func doCopy(target *runcommands.S3Target, s store.Store) error {
	if err := s.Write(os.Stdin); err != nil {
		log.Errorf("Error writing to the bucket: %s", err.Error())
		return err
	}
	return nil
}

func doPaste(target *runcommands.S3Target, s store.Store) (string, error) {
	content, err := s.Read()
	if err != nil {
		log.Errorf("Error reading to the bucket.. Have you copied yet? %s", err.Error())
		return "", err
	}
	return content, nil

}
