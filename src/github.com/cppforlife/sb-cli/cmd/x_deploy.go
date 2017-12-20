package cmd

import (
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	boshx "github.com/cppforlife/sb-cli/x"
)

type XDeployCmd struct {
	ui boshui.UI
	fs boshsys.FileSystem
}

func NewXDeployCmd(ui boshui.UI, fs boshsys.FileSystem) XDeployCmd {
	return XDeployCmd{ui, fs}
}

func (c XDeployCmd) Run(opts XDeployOpts) error {
	tree := boshx.NewDependencyTree(boshx.NewCLIImpl(c.ui), c.fs)

	rootDep, err := tree.Build(opts.Deployment, opts.Args.ManifestPath)
	if err != nil {
		return err
	}

	args := opts.ExtraArgs

	vars, args, err := boshx.NewVarsFromArgs(args)
	if err != nil {
		return err
	}

	opsFiles, args, err := boshx.NewOpsFilesFromArgs(args)
	if err != nil {
		return err
	}

	DeploymentsTable{rootDep, c.ui}.Print()

	_, err = rootDep.Interpolate(vars, opsFiles, args)
	if err != nil {
		return err
	}

	err = rootDep.Deploy(vars, opsFiles, args)
	if err != nil {
		return err
	}

	return nil
}
