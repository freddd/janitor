package tracker

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/freddd/janitor/config"
	"github.com/mitchellh/cli"
	"math"
	"os"
	"regexp"
	"strings"
	"github.com/freddd/janitor/util"
	"path/filepath"
)

const (
	base64 string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	hex    string = "1234567890abcdefABCDEF"
	minimumEntropy float64 = 4.8
)

var keys = map[string][]*regexp.Regexp{
	"aws":        {regexp.MustCompile("[0-9a-zA-Z/+]{40}")},
	"bitly":      {regexp.MustCompile("R_[0-9a-f]{32}")},
	"facebook":   {regexp.MustCompile("[0-9a-f]{32}")},
	"flickr":     {regexp.MustCompile("[0-9a-f]{16}")},
	"foursquare": {regexp.MustCompile("[0-9A-Z]{48}")},
	// "linkedin":{regexp.MustCompile("[0-9a-zA-Z]{16}")}, This regexp basically catches everything
	"twitter":   {regexp.MustCompile("[0-9a-zA-Z]{35,44}")},
	"google":    {regexp.MustCompile("(AIza.{35})")},
	"mailchimp": {regexp.MustCompile("[0-9a-z]{32}(-us[12])?")},
	"github":    {regexp.MustCompile("[0-9A-F]{40}")},
	"slack":     {regexp.MustCompile("^xoxb-"), regexp.MustCompile("^xoxp-"), regexp.MustCompile("^xoxa-")},
	"ssh":       {regexp.MustCompile("ssh-rsa AAAA[0-9A-Za-z+/]+[=]{0,3}( [^@]+@[^@]+)?")},
}

type Tracker struct {
	Cfg *config.Tracker
	Ui  cli.Ui
}

func (tracker *Tracker) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("cfg", flag.ExitOnError)
	cmdFlags.Usage = func() { tracker.Ui.Output(tracker.Help()) }
	cfgPath := ""
	cmdFlags.StringVar(&cfgPath, "cfg", "", "Path to the config")

	if len(args) < 1 {
		cmdFlags.Usage()
		return 1
	}

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		tracker.Ui.Error(err.Error())
		return 1
	}
	tracker.Cfg = &cfg.Tracker

	pathToRepo, err := util.CurrentDir()
	if err != nil {
		tracker.Ui.Error(err.Error())
		return 1
	}

	tracker.Ui.Info(fmt.Sprintf("Running on path: %s", pathToRepo))
	files, err := util.FindAllFiles(pathToRepo, []string{})
	if err != nil {
		tracker.Ui.Error(err.Error())
	}

	tracker.Ui.Info("---------- Finding secrets: ------------------------------")
	tracker.FindAllPossibleKeys(files)
	tracker.Ui.Info("----------------------------------------------------------")
	return 0
}

func (tracker *Tracker) Help() string {
	helpText := `
		Usage: janitor tracker
		  Recursively searches for secrets in the current folder
		Options:
		  -cfg  the global config file (mandatory)
		`

	return strings.TrimSpace(helpText)
}

func (tracker *Tracker) Synopsis() string {
	return "Recursively finds secrets in the current dir"
}

func (tracker *Tracker) FindAllPossibleKeys(files []string) {
	for _, file := range files {
		err := tracker.process(file)
		if err != nil {
			tracker.Ui.Error(err.Error())
			continue
		}
	}
}

func (tracker *Tracker) process(path string) error {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return errors.New("not yet implemented")
	} else {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber++
			text := strings.TrimSpace(scanner.Text())
			split := strings.Split(text, " ")

			keyword, seed := tracker.seed(text, path)
			for _, word := range split {
				entropy := shannonEntropy(word, seed)
				if entropy > minimumEntropy {
					tracker.Ui.Info("----------------------------------------------------------")
					matches := matchesRegexp(word)
					if len(matches) > 0 {
						for _, match := range matches {
							if strings.Contains(strings.ToLower(text), match) {
								tracker.Ui.Info(fmt.Sprintf("Matched vendor: %s", match))
							}
						}
					}
					if keyword != "" {
						tracker.Ui.Info(fmt.Sprintf("Matched keyword: %s", keyword))
					}
					tracker.Ui.Info(fmt.Sprintf("File: %s", path))
					tracker.Ui.Info(fmt.Sprintf("Line: %d", lineNumber))
					tracker.Ui.Info(fmt.Sprintf("Entropy: %f", entropy))
					tracker.Ui.Info(fmt.Sprintf("Text: %s", text))
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		return nil
	}
}

// Should be optimized
func (tracker *Tracker) seed(s string, path string) (string, float64) {
	lower := strings.ToLower(s)
	seed := 0.0
	key := ""
	for _, keyword := range tracker.Cfg.Keywords {
		if strings.Contains(lower, keyword) {
			seed += 0.2
			key = keyword
			break
		}
	}

	for _, fileName := range tracker.Cfg.FileNames {
		name := filepath.Base(path)
		if strings.Contains(name, fileName) {
			seed += 0.2
			break
		}
	}

	return key, seed
}

// https://rosettacode.org/wiki/Entropy#Go
func shannonEntropy(s string, seed float64) float64 {
	entropy := seed
	if s == "" || len(s) < 10 {
		return entropy
	}

	for i := 0; i < 256; i++ {
		px := float64(strings.Count(s, string(byte(i)))) / float64(len(s))
		if px > 0 {
			entropy += -px * math.Log2(px)
		}
	}
	return entropy
}

func matchesRegexp(s string) []string {
	var matches []string
	for name, listOfRegexp := range keys {
		for _, regex := range listOfRegexp {
			if regex.MatchString(s) {
				matches = append(matches, name)
			}
		}
	}
	return matches
}
