package github

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"os"
	"strings"
	"github.com/mitchellh/cli"
	"flag"
)

const (
	BaseUrl = "https://api.github.com/"
	tfaPath = "orgs/%s/members?filter=2fa_disabled"
	githubKey = "GITHUB_KEY"
	githubOrg = "GITHUB_ORG"
)

type TFAResponse []struct {
	Login            string `json:"login"`
	ID               int    `json:"id"`
	AvatarURL        string `json:"avatar_url"`
	GravatarID       string `json:"gravatar_id"`
	URL              string `json:"url"`
	HTMLURL          string `json:"html_url"`
	OrganizationsURL string `json:"organizations_url"`
	Type             string `json:"type"`
	SiteAdmin        bool   `json:"site_admin"`
}

type GitHub struct {
	Ui  cli.Ui
}

func (github *GitHub) FindAllUsersWithoutTFA(organization string, apiKey string) ([]string, []error) {
	targetUrl := BaseUrl + fmt.Sprintf(tfaPath, organization)
	request := gorequest.New()
	res, body, errs := request.Get(targetUrl).Set("Authorization", fmt.Sprintf("token %s", apiKey)).End()

	if res.StatusCode != 200 {
		github.Ui.Error(fmt.Sprintf("CRITICAL: Got status code %d from Github", res.StatusCode))
	}
	if errs != nil && len(errs) > 0 {
		return nil, errs
	}

	response := TFAResponse{}
	json.Unmarshal([]byte(body), &response)

	users := make([]string, len(response))
	for _, user := range response {
		users = append(users, user.Login)
	}

	return users, nil
}

func (github *GitHub) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("github", flag.ExitOnError)
	cmdFlags.Usage = func() { github.Ui.Output(github.Help()) }
	apiKey := ""
	organization := ""
	cmdFlags.StringVar(&apiKey, "apiKey", "", "The api key")
	cmdFlags.StringVar(&organization, "organization", "", "The organization")

	if err := cmdFlags.Parse(args); err != nil {
		cmdFlags.Usage()
		return 1
	}

	if apiKey == "" {
		apiKey = os.Getenv(githubKey)
		if apiKey == "" {
			cmdFlags.Usage()
			return 1
		}
	}

	if organization == "" {
		organization = os.Getenv(githubOrg)
		if organization == "" {
			cmdFlags.Usage()
			return 1
		}
	}


	return 0
}

func (github *GitHub) Help() string {
	helpText := `
		Usage: janitor tfa github --apiKey <key> --organization <org>
		  Gets the current status of TFA usage in an organization on github using an api key
		Options:
		  --apiKey  the key with permissions to get the info from github (can also be set using the GITHUB_KEY env variable)
		  --organization the organization we aim to get the info from (can also be set using the GITHUB_ORG env variable)
		`

	return strings.TrimSpace(helpText)
}

func (github *GitHub) Synopsis() string {
	return "Gets the current status of TFA usage in an organization on github using an api key"
}
