package appcloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/pkg/errors"
)

// EnableSSL installs an existing SSL certificate.
func (p *Plugin) EnableSSL(c plugin.CliConnection, domain string, hostname string) error {
	un, err := c.Username()
	if err != nil {
		return errors.Wrap(err, "Couldn't get your username")
	}

	fullDomain := domain
	if hostname != "" {
		fullDomain = strings.Join([]string{hostname, domain}, ".")
	}

	p.ui.Say("Enabling SSL for %s as %s...", terminal.EntityNameColor(fullDomain), terminal.EntityNameColor(un))

	s, err := c.GetCurrentSpace()
	if err != nil {
		return errors.Wrap(err, "Couldn't retrieve current space")
	}

	req := SSLCertificateRequest{
		SpaceID:        s.Guid,
		FullDomainName: fullDomain,
	}
	reqData, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "Couldn't parse JSON data")
	}

	url := "/custom/certifications/install"
	resLines, err := c.CliCommandWithoutTerminalOutput("curl", "-X", "PUT", url, "-d", string(reqData))
	if err != nil {
		return errors.Wrap(err, "Couldn't enable SSL for route")
	}

	resString := strings.Join(resLines, "")
	var res SSLCertificateResponse
	err = json.Unmarshal([]byte(resString), &res)
	if err != nil {
		return errors.Wrap(err, "Couldn't read JSON response from server")
	}

	if res.ErrorCode != "" {
		return fmt.Errorf("Error response from server: %s", res.Description)
	}

	p.ui.Say(terminal.SuccessColor("OK\n"))

	p.ui.Say("SSL enabled")

	return nil
}
