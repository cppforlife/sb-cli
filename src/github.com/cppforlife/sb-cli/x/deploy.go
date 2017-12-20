package x

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

func (d Dependency) Deploy(vars []Variable, opsFiles []OpsFile, rawArgs []string) error {
	vars1, err := d.applicableVars(vars)
	if err != nil {
		return err
	}

	opsFiles1, err := d.applicableOps(opsFiles)
	if err != nil {
		return err
	}

	for _, dep := range d.Dependencies {
		err := dep.Deploy(d.applicableVars2(vars), d.applicableOps2(opsFiles), rawArgs)
		if err != nil {
			return err
		}
	}

	err = d.deploy(vars1, opsFiles1, rawArgs)
	if err != nil {
		return err
	}

	excessDeplNames, err := d.listExcessDeployments()
	if err != nil {
		return bosherr.WrapErrorf(err, "Determining excess deployments")
	}

	for _, deplName := range excessDeplNames {
		err = d.delete(deplName, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d Dependency) deploy(vars []Variable, opsFiles []OpsFile, rawArgs []string) error {
	args := []string{"deploy", "--tty", "-n", "-d", d.DeploymentName, "-"} // todo interactive?

	for _, v1 := range vars {
		args = append(args, []string{"-v", fmt.Sprintf("%s=%s", v1.Name, v1.Value)}...)
	}

	for _, op := range opsFiles {
		args = append(args, []string{"-o", op.Path}...)
	}

	args = append(args, rawArgs...)

	manifest, err := d.manifestWithAdjustedName()
	if err != nil {
		return bosherr.WrapErrorf(err, "Deploying deployment '%s'", d.DeploymentName)
	}

	return d.executeWithInfo(args, bytes.NewReader(manifest))
}

type cliDeploymentsJSON struct{ Tables []cliDeploymentsTableJSON }
type cliDeploymentsTableJSON struct{ Rows []cliDeploymentsRowJSON }
type cliDeploymentsRowJSON struct{ Name string }

func (d Dependency) listExcessDeployments() ([]string, error) {
	args := []string{"deployments", "--tty", "-n", "--column", "name", "--json"}

	stdout, err := d.cli.Execute(args, nil)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Getting deployments output")
	}

	var allDepls cliDeploymentsJSON

	err = json.Unmarshal(stdout, &allDepls)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Unmarshaling deployments output")
	}

	if len(allDepls.Tables) != 1 {
		return nil, bosherr.Errorf("Expected deployments output to contain one table")
	}

	var excessDeplNames []string

	for _, depl := range allDepls.Tables[0].Rows {
		if strings.HasPrefix(depl.Name, d.DeploymentName+"-") {
			var found bool

			for _, dep := range d.Dependencies {
				if depl.Name == dep.DeploymentName {
					found = true
				}
			}

			if !found {
				excessDeplNames = append(excessDeplNames, depl.Name)
			}
		}
	}

	return excessDeplNames, nil
}

func (d Dependency) executeWithInfo(args []string, stdin io.Reader) error {
	// todo use ui
	fmt.Printf("-----> bosh %s\n", strings.Join(args, " "))

	// todo streaming output...
	stdout, err := d.cli.Execute(args, stdin)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", stdout)

	return nil
}
