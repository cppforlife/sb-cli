package x

import (
	"bytes"
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Manifest struct {
	Dependency *Dependency
	Bytes      []byte
}

func (d Dependency) Interpolate(vars []Variable, opsFiles []OpsFile, rawArgs []string) ([]Manifest, error) {
	vars1, err := d.applicableVars(vars)
	if err != nil {
		return nil, err
	}

	opsFiles1, err := d.applicableOps(opsFiles)
	if err != nil {
		return nil, err
	}

	var manifests []Manifest

	for _, dep := range d.Dependencies {
		manifests2, err := dep.Interpolate(d.applicableVars2(vars), d.applicableOps2(opsFiles), rawArgs)
		if err != nil {
			return nil, err
		}

		manifests = append(manifests, manifests2...)
	}

	manifest, err := d.interpolate(vars1, opsFiles1, rawArgs)
	if err != nil {
		return nil, err
	}

	return append(manifests, manifest), nil
}

func (d *Dependency) interpolate(vars []Variable, opsFiles []OpsFile, rawArgs []string) (Manifest, error) {
	args := []string{"interpolate", "-n", "-"}

	for _, v1 := range vars {
		args = append(args, []string{"-v", fmt.Sprintf("%s=%s", v1.Name, v1.Value)}...)
	}

	for _, op := range opsFiles {
		args = append(args, []string{"-o", op.Path}...)
	}

	args = append(args, rawArgs...)

	manifest, err := d.manifestWithAdjustedName()
	if err != nil {
		return Manifest{}, err
	}

	stdout, err := d.cli.Execute(args, bytes.NewReader(manifest))
	if err != nil {
		return Manifest{}, bosherr.WrapErrorf(err, "Interpolating deployment '%s'", d.DeploymentName)
	}

	return Manifest{Dependency: d, Bytes: stdout}, nil
}
