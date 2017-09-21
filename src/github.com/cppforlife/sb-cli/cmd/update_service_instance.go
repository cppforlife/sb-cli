package cmd

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type UpdateServiceInstanceCmd struct {
	siFactory ServiceInstanceFactory
}

func NewUpdateServiceInstanceCmd(siFactory ServiceInstanceFactory) UpdateServiceInstanceCmd {
	return UpdateServiceInstanceCmd{siFactory}
}

func (c UpdateServiceInstanceCmd) Run(opts UpdateServiceInstanceOpts) error {
	params, err := opts.ParamFlags.AsParams()
	if err != nil {
		return err
	}

	finder := ServiceInstanceFinder{
		ID:            opts.Args.ID,
		ServiceID:     opts.ServiceID,
		ServicePlanID: opts.ServicePlanID,
	}

	si, err := c.siFactory.New(finder)
	if err != nil {
		return err
	}

	err = si.Update(params)
	if err != nil {
		return bosherr.WrapError(err, "Updating service instance")
	}

	return nil
}
