package domain

import (
	"crypto/tls"
	"fmt"
	"time"
	"crypto/x509"
	"strings"
	"github.com/mitchellh/cli"
	"flag"
	"github.com/likexian/whois-go"
	"github.com/likexian/whois-parser-go"
	"github.com/olekukonko/tablewriter"
	"os"
)

const whoisTimeFormat string = "2006-01-02T15:04:05.00Z"

var sunset = map[x509.SignatureAlgorithm]string{
	x509.MD2WithRSA: "MD2 with RSA",
	x509.MD5WithRSA: "MD5 with RSA",
	x509.SHA1WithRSA: "SHA1 with RSA",
	x509.DSAWithSHA1: "DSA with SHA1",
	x509.ECDSAWithSHA1: "ECDSA with SHA1",
}
type DomainVerifier struct {
	Ui cli.Ui
}

func (d *DomainVerifier) Run(args []string) int {
	// TODO refactor
	cmdFlags := flag.NewFlagSet("host", flag.ExitOnError)
	cmdFlags.Usage = func() { d.Ui.Output(d.Help()) }
	host := ""
	cmdFlags.StringVar(&host, "host", "", "The host to check")

	if len(args) < 1 {
		cmdFlags.Usage()
		return 1
	}

	if err := cmdFlags.Parse(args); err != nil {
		cmdFlags.Usage()
		return 1
	}

	if host == "" {
		cmdFlags.Usage()
		return 1
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Status", "Type", "Message"})
	table.SetRowLine(true)
	d.Ui.Info("---------- Validating Certificate and Domain: ----------")
	d.ValidateCert(host, time.Now(), table)
	d.ValidateDomain(host, time.Now(), table)
	table.Render()
	return 0
}

func (d *DomainVerifier) Help() string {
	helpText := `
		Usage: janitor domain --host
		  Checks the host SSL domain expiry and if it's using an algo that is unsafe
		Options:
		  -host  the host to check (mandatory)
		`

	return strings.TrimSpace(helpText)
}

func (d *DomainVerifier) Synopsis() string {
	return "Checks the host SSL domain expiry and if it's using an algo that is unsafe"
}

func (d DomainVerifier) ValidateCert(host string, whenToWarn time.Time, table *tablewriter.Table) {
	d.Ui.Info(fmt.Sprintf("Checking expiry of domain for host=%s", host))
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", host, "443"), nil)
	if err != nil {
		d.Ui.Error(err.Error())
		return
	}
	defer conn.Close()

	for _, chain := range conn.ConnectionState().VerifiedChains {
		for i, cert := range chain {
			if i != len(chain)-1 {
				algorithm := sunset[cert.SignatureAlgorithm]
				if algorithm != "" {
					table.Append([]string{"CRITICAL", "Certificate", fmt.Sprintf("Cert is using a unsafe algo: %s, dns names: %+v", algorithm, cert.DNSNames)})
				}
			}
			if contains(cert.DNSNames, host) {
				if whenToWarn.After(cert.NotAfter) {
					table.Append([]string{"CRITICAL", "Certificate", fmt.Sprintf("Cert already expired --- Host: %+v, Expiry: %+v", cert.DNSNames, cert.NotAfter)})
				}
				if int(cert.NotAfter.Sub(whenToWarn) / (24 * time.Hour)) < 30 {
					table.Append([]string{"WARNING", "Certificate", fmt.Sprintf("Cert expiring within 30 days --- Host: %+v, Expiry: %+v", cert.DNSNames, cert.NotAfter)})
				} else {
					table.Append([]string{"OK", "Certificate", fmt.Sprintf("Host: %+v, Expiry: %+v", cert.DNSNames, cert.NotAfter)})
				}
			}
		}
	}
}

func (d DomainVerifier) ValidateDomain(host string, whenToWarn time.Time, table *tablewriter.Table) {
	whoisResult, err := whois.Whois(host)
	if err != nil {
		d.Ui.Error(err.Error())
	}
	parsed, err := whois_parser.Parser(whoisResult)
	if err != nil {
		d.Ui.Error(err.Error())
	}
	// Print the domain status
	if strings.HasPrefix(parsed.Registrar.DomainStatus, "ok") {
		table.Append([]string{"OK", "Domain", fmt.Sprintf("Domain status: %s", parsed.Registrar.DomainStatus)})
	} else {
		table.Append([]string{"WARNING", "Domain", fmt.Sprintf("Domain status: %s", parsed.Registrar.DomainStatus)})
	}

	// Print the domain expiration date
	expirationDate, err := time.Parse(whoisTimeFormat, parsed.Registrar.ExpirationDate)
	if err != nil {
		d.Ui.Error(err.Error())
	}

	if whenToWarn.After(expirationDate) {
		table.Append([]string{"CRITICAL", "Domain", fmt.Sprintf("Domain (%s) has already expired: %s", host, expirationDate)})
	} else if int(expirationDate.Sub(whenToWarn) / (24 * time.Hour)) < 30 {
		table.Append([]string{"WARNING", "Domain", fmt.Sprintf("Domain (%s) is expiring within 30 days: %s", host, expirationDate)})
	} else {
		table.Append([]string{"OK", "Domain", fmt.Sprintf("Domain (%s) is expiring at: %s", host, expirationDate)})
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}