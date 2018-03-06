package mining

import (
	"bufio"
	"fmt"
	"github.com/freddd/janitor/util"
	"github.com/mitchellh/cli"
	"os"
	"strings"
	"regexp"
)

const (
	httpRegexp string = `(http|https)://`
	ipRegexp string = `(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])`
)

type Mining struct {
	Ui cli.Ui
}

func (m *Mining) Run(args []string) int {
	m.Ui.Info("---------- Mining information: ---------------------------")
	pathToRepo, err := util.CurrentDir()
	if err != nil {
		m.Ui.Error(err.Error())
		return 1
	}

	m.Ui.Info(fmt.Sprintf("Running on path: %s", pathToRepo))
	files, err := util.FindAllFiles(pathToRepo, []string{"vendor", ".git", "node_modules"})
	if err != nil {
		m.Ui.Error(err.Error())
	}
	m.Ui.Info("---------- Result: ---------------------------------------")
	hosts := m.findAllHosts(files)
	for _, host := range hosts {
		m.Ui.Info(fmt.Sprintf("OK: Host: %s", host))
	}
	m.Ui.Info("----------------------------------------------------------")
	return 0
}

func (m *Mining) findAllHosts(files []string) []string {
	var result []string
	for _, file := range files {
		hosts, err := m.findHostsInFile(file)
		if err != nil {
			m.Ui.Error(err.Error())
			continue
		}
		result = append(result, hosts...)
	}
	return result
}

func (m *Mining) findHostsInFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hosts []string
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		split := strings.Split(text, " ")
		for _, word := range split {
			// g.cn is the shortest
			if len(word) < 4 {
				continue
			}

			matchesHttp, err := regexp.MatchString(httpRegexp, word)
			if err != nil {
				return nil, err
			}

			matchesIp, err := regexp.MatchString(ipRegexp, word)
			if err != nil {
				return nil, err
			}

			if matchesHttp || matchesIp {
				hosts = append(hosts, word)
			}
		}
	}
	return hosts, nil
}

func (m *Mining) Help() string {
	helpText := `
		Usage: janitor mining
		  Mining the current directory for information (usually done in a repo)
		`

	return strings.TrimSpace(helpText)
}

func (m *Mining) Synopsis() string {
	return "Mining the current directory for information (usually done in a repo)"
}
