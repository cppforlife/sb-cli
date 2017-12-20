package cmd

import (
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	boshx "github.com/cppforlife/sb-cli/x"
)

type XInterpolateCmd struct {
	ui boshui.UI
	fs boshsys.FileSystem
}

func NewXInterpolateCmd(ui boshui.UI, fs boshsys.FileSystem) XInterpolateCmd {
	return XInterpolateCmd{ui, fs}
}

func (c XInterpolateCmd) Run(opts XInterpolateOpts) error {
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

	manifests, err := rootDep.Interpolate(vars, opsFiles, args)
	if err != nil {
		return err
	}

	ManifestsTable{manifests, c.ui}.Print()

	return nil
}
