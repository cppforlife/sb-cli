package cmd

import (
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type DeleteServiceInstanceCmd struct {
	siFactory ServiceInstanceFactory
	ui        boshui.UI
}

func NewDeleteServiceInstanceCmd(siFactory ServiceInstanceFactory, ui boshui.UI) DeleteServiceInstanceCmd {
	return DeleteServiceInstanceCmd{siFactory, ui}
}

func (c DeleteServiceInstanceCmd) Run(opts DeleteServiceInstanceOpts) error {
	si := c.siFactory.New(ServiceInstanceFinder{
		Name:            opts.Args.Name,
		ServiceName:     opts.Args.ServiceName,
		ServicePlanName: opts.Args.ServicePlanName,
	})

	err := si.Deprovision()
	if err != nil {
		return bosherr.WrapError(err, "Deprovisioning instance")
	}

	return nil
}
