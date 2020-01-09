package commands

import (
	"github.com/alec-rabold/copypasta/runcommands"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

// TargetCommand is responsible for setting the target
type TargetCommand struct {
	UI cli.Ui
}

// Help string for function
func (t *TargetCommand) Help() string {
	return `Usage: copy-pasta target [<target>

	Changes the current target to the specified target.
	If no argument is provided, it lists the current target.
	`
}

// Run command for function
func (t *TargetCommand) Run(args []string) int {
	target, err := runcommands.Load()
	if err != nil {
		log.Errorf("Error loading run commands: %s", err.Error())
		return 1
	}
	if len(args) > 0 {

	} else {
		t.UI.Output("copypasta current target:")
		t.UI.Output("	" + target.Name)
		return 0
	}
	return 0
}
