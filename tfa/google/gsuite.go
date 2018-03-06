package google

import (
	"strings"
	"github.com/mitchellh/cli"
)

const gsuiteKey string = "GSUITE_KEY"

type Gsuite struct {
	Ui cli.Ui
}

func (g *Gsuite) Run(args []string) int {
	g.Ui.Info("---------- Finding users in GSuite without TFA: ----------")
	g.Ui.Error("Not yet implemented!")
	g.Ui.Info("----------------------------------------------------------")
	return 0
}

func (g *Gsuite) Help() string {
	helpText := `
		Usage: janitor tfa gsuite
		  Gets the current status of TFA usage on github/gsuite
		Options:
		  --apiKey  the key with permissions to get the info from github (can also be set using the GSUITE_KEY env variable)
		`

	return strings.TrimSpace(helpText)
}

func (g *Gsuite) Synopsis() string {
	return "Check TFA status in gsuite"
}