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
	finder := ServiceInstanceFinder{
		ID:            opts.Args.ID,
		ServiceID:     opts.ServiceID,
		ServicePlanID: opts.ServicePlanID,
	}

	si, err := c.siFactory.New(finder)
	if err != nil {
		return err
	}

	err = si.Deprovision()
	if err != nil {
		return bosherr.WrapError(err, "Deprovisioning instance")
	}

	return nil
}
