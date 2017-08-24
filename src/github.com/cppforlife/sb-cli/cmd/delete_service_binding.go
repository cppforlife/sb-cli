package cmd

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type DeleteServiceBindingCmd struct {
	siFactory ServiceInstanceFactory
}

func NewDeleteServiceBindingCmd(siFactory ServiceInstanceFactory) DeleteServiceBindingCmd {
	return DeleteServiceBindingCmd{siFactory}
}

func (c DeleteServiceBindingCmd) Run(opts DeleteServiceBindingOpts) error {
	finder := ServiceInstanceFinder{
		ID:            opts.Args.ServiceInstanceID,
		ServiceID:     opts.ServiceID,
		ServicePlanID: opts.ServicePlanID,
	}

	si, err := c.siFactory.New(finder)
	if err != nil {
		return err
	}

	err = si.Unbind(opts.Args.ID)
	if err != nil {
		return bosherr.WrapError(err, "Unbinding from instance")
	}

	return nil
}
