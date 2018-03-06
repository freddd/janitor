package main

import (
	"github.com/freddd/janitor/tracker"
	"github.com/mitchellh/cli"
	"os"
	"github.com/freddd/janitor/tfa"
	"github.com/freddd/janitor/domain"
	"github.com/freddd/janitor/mining"
)

func main() {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	c := cli.NewCLI("Janitor", "0.0.1")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"tracker": func() (cli.Command, error) {
			return &tracker.Tracker{
				Ui: getUi(ui),
			}, nil
		},
		"tfa": func() (cli.Command, error) {
			return &tfa.TfaCommand{
				Ui: getUi(ui),
			}, nil
		},
		"domain": func() (cli.Command, error) {
			return &domain.DomainVerifier{
				Ui: getUi(ui),
			}, nil
		},
		"mining": func() (cli.Command, error) {
			return &mining.Mining{
				Ui: getUi(ui),
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		getUi(ui).Error(err.Error())
	}

	os.Exit(exitStatus)
}

func getUi(ui cli.Ui) *cli.ColoredUi {
	return &cli.ColoredUi{
		Ui:          ui,
		OutputColor: cli.UiColorBlue,
		InfoColor:   cli.UiColorGreen,
		ErrorColor:  cli.UiColorRed,
		WarnColor:   cli.UiColorYellow,
	}
}
