package cmd

import (
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type CreateServiceInstanceCmd struct {
	siFactory ServiceInstanceFactory
	ui        boshui.UI
}

func NewCreateServiceInstanceCmd(siFactory ServiceInstanceFactory, ui boshui.UI) CreateServiceInstanceCmd {
	return CreateServiceInstanceCmd{siFactory, ui}
}

func (c CreateServiceInstanceCmd) Run(opts CreateServiceInstanceOpts) error {
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

	err = si.Provision(params)
	if err != nil {
		return bosherr.WrapError(err, "Provisioning service instance")
	}

	info := boshtbl.Table{
		Header: []boshtbl.Header{
			boshtbl.NewHeader("ID"),
			boshtbl.NewHeader("Service ID"),
			boshtbl.NewHeader("Service Plan ID"),
			boshtbl.NewHeader("Dashboard URL"),
		},
		Rows: [][]boshtbl.Value{
			{
				boshtbl.NewValueString(si.ID()),
				boshtbl.NewValueString(si.ServiceID()),
				boshtbl.NewValueString(si.ServicePlanID()),
				boshtbl.NewValueString(si.DashboardURL()),
			},
		},
		Transpose: true,
	}

	c.ui.PrintTable(info)

	return nil
}
