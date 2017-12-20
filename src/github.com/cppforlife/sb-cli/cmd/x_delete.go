package cmd

import (
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	boshx "github.com/cppforlife/sb-cli/x"
)

type XDeleteCmd struct {
	ui boshui.UI
	fs boshsys.FileSystem
}

func NewXDeleteCmd(ui boshui.UI, fs boshsys.FileSystem) XDeleteCmd {
	return XDeleteCmd{ui, fs}
}

func (c XDeleteCmd) Run(opts XDeleteOpts) error {
	tree := boshx.NewDependencyTree(boshx.NewCLIImpl(c.ui), c.fs)

	rootDep, err := tree.Build(opts.Deployment, opts.Args.ManifestPath)
	if err != nil {
		return err
	}

	DeploymentsTable{rootDep, c.ui}.Print()

	err = rootDep.Delete(opts.ExtraArgs)
	if err != nil {
		return err
	}

	return nil
}
