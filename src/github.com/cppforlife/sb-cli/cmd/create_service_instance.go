package cmd

import (
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"gopkg.in/yaml.v2"
)

type CreateServiceInstanceCmd struct {
	siFactory ServiceInstanceFactory
	ui        boshui.UI
}

func NewCreateServiceInstanceCmd(siFactory ServiceInstanceFactory, ui boshui.UI) CreateServiceInstanceCmd {
	return CreateServiceInstanceCmd{siFactory, ui}
}

func (c CreateServiceInstanceCmd) Run(opts CreateServiceInstanceOpts) error {
	var params map[string]interface{}

	if len(opts.Params.Bytes) > 0 {
		err := yaml.Unmarshal(opts.Params.Bytes, params)
		if err != nil {
			return bosherr.WrapError(err, "Unmarshaling instance params")
		}
	}

	si := c.siFactory.New(ServiceInstanceFinder{
		Name:            opts.Args.Name,
		ServiceName:     opts.Args.ServiceName,
		ServicePlanName: opts.Args.ServicePlanName,
	})

	err := si.Provision(params)
	if err != nil {
		return bosherr.WrapError(err, "Provisioning instance")
	}

	info := boshtbl.Table{
		Header: []boshtbl.Header{
			boshtbl.NewHeader("Name"),
			boshtbl.NewHeader("Service"),
			boshtbl.NewHeader("Service Plan"),
			boshtbl.NewHeader("Dashboard URL"),
		},
		Rows: [][]boshtbl.Value{
			{
				boshtbl.NewValueString(si.Name()),
				boshtbl.NewValueString(si.ServiceName()),
				boshtbl.NewValueString(si.ServicePlanName()),
				boshtbl.NewValueString(si.DashboardURL()),
			},
		},
		Transpose: true,
	}

	c.ui.PrintTable(info)

	return nil
}
