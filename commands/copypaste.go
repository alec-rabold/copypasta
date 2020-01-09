package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/alec-rabold/copypasta/runcommands"
	"github.com/alec-rabold/copypasta/store"
	"github.com/alec-rabold/copypasta/utils"
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
	// recipient := copyPasteCommand.String("recipient", "", "defines the gpg recipient")

	var recipient string
	copyPasteCommand.StringVar(&recipient, "r", "", "defines the gpg recipient")

	if err := copyPasteCommand.Parse(args); err != nil {
		return 1
	}

	store, err := store.NewStore(target)
	if err != nil {
		log.Errorf("Error creating minio s3 client %s", err.Error())
		return 1
	}
	if *optionPaste {
		if err := doPaste(target, store); err != nil {
			log.Errorf("Error retrieving content: %s", err.Error())
			return 1
		}
	} else {
		if recipient == "" {
			println(c.Help())
			os.Exit(1)
		}
		println("Encrypting data for recipient: " + recipient)
		if err := doCopy(target, store, recipient); err != nil {
			log.Errorf("Error copying data: %s", err.Error())
		} else {
			println("Data copied to s3, ready to paste!")
		}
	}
	return 0
}

func doCopy(target *runcommands.S3Target, s store.Store, recipient string) error {
	cipher, err := utils.EncryptFile(os.Stdin, recipient)
	if err != nil {
		log.Errorf("Error encrypting data: %s", err.Error())
		return err
	}
	if err = s.Write(cipher); err != nil {
		log.Errorf("Error writing to the bucket: %s", err.Error())
		return err
	}
	return nil
}

func doPaste(target *runcommands.S3Target, s store.Store) error {
	cipher, err := s.Read()
	// cipherBytes, err := ioutil.ReadAll(reader)
	// cipher := string(cipherBytes)
	if err != nil {
		log.Errorf("Error reading from the bucket.. Have you copied yet? %s", err.Error())
		return err
	}
	content, err := utils.DecryptFile(cipher)
	if err != nil {
		log.Errorf("Error decrypting message: %s", err.Error())
		return err
	}
	fmt.Print(content)

	return nil

}

// Help string
func (c *CopyPasteCommand) Help() string {
	return `Usage to paste: copypasta [--paste]
Usage to copy: <some command with output> | copypasta [--recipient <name>]
    Copy or paste using copypasta. Use --paste to force copypasta to
		ignore its stdin and output from the current target.
`
}

// Synopsis is the short help string
func (c *CopyPasteCommand) Synopsis() string {
	return "Copy or paste using copypasta"
}
