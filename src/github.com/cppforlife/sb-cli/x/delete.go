package x

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

func (d Dependency) Delete(rawArgs []string) error {
	for _, dep := range d.Dependencies {
		err := dep.Delete(rawArgs)
		if err != nil {
			return err
		}
	}

	err := d.delete(d.DeploymentName, rawArgs)
	if err != nil {
		return err
	}

	return nil
}

func (d Dependency) delete(deplName string, rawArgs []string) error {
	args := []string{"delete-deployment", "--tty", "-n", "-d", deplName}

	args = append(args, rawArgs...)

	err := d.executeWithInfo(args, nil)
	if err != nil {
		return bosherr.WrapErrorf(err, "Deleting deployment '%s'", deplName)
	}

	return nil
}
