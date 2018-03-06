package tfa

import (
	"github.com/mitchellh/cli"
	"strings"
	"github.com/freddd/janitor/tfa/github"
	"github.com/freddd/janitor/tfa/google"
)

type TfaCommand struct {
	Ui cli.Ui
}

func (t *TfaCommand) Run(args []string) int {
	tfa := cli.NewCLI("tfa", "")
	tfa.Args = args

	tfa.Commands = map[string]cli.CommandFactory{
		"github": func() (cli.Command, error) {
			return &github.GitHub{Ui: t.Ui}, nil
		},
		"gsuite": func() (cli.Command, error) {
			return &google.Gsuite{Ui: t.Ui}, nil
		},
	}

	if exitStatus, err := tfa.Run(); err != nil {
		t.Ui.Error(err.Error())
		return exitStatus
	} else {
		return exitStatus
	}
	return 0
}

func (t *TfaCommand) Help() string {
	helpText := `
		Usage: janitor tfa <github/gsuite>
		  Gets the current status of TFA usage on github/gsuite
		`

	return strings.TrimSpace(helpText)
}

func (t *TfaCommand) Synopsis() string {
	return "Check TFA status in either github or gsuite"
}