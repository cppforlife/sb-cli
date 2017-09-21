package cmd

import (
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
	"gopkg.in/yaml.v2"
)

type CreateServiceBindingCmd struct {
	siFactory ServiceInstanceFactory
	uuidGen   boshuuid.Generator
	ui        boshui.UI
}

func NewCreateServiceBindingCmd(siFactory ServiceInstanceFactory, uuidGen boshuuid.Generator, ui boshui.UI) CreateServiceBindingCmd {
	return CreateServiceBindingCmd{siFactory, uuidGen, ui}
}

func (c CreateServiceBindingCmd) Run(opts CreateServiceBindingOpts) error {
	var resource *osb.BindResource

	if len(opts.Resource.Bytes) > 0 {
		err := yaml.Unmarshal(opts.Resource.Bytes, resource)
		if err != nil {
			return bosherr.WrapError(err, "Unmarshaling service binding resource")
		}
	}

	params, err := opts.ParamFlags.AsParams()
	if err != nil {
		return err
	}

	finder := ServiceInstanceFinder{
		ID:            opts.Args.ServiceInstanceID,
		ServiceID:     opts.ServiceID,
		ServicePlanID: opts.ServicePlanID,
	}

	si, err := c.siFactory.New(finder)
	if err != nil {
		return err
	}

	bindingID := opts.ID

	if len(bindingID) == 0 {
		uuid, err := c.uuidGen.Generate()
		if err != nil {
			return bosherr.WrapError(err, "Generating service binding ID")
		}

		bindingID = uuid
	}

	sb, err := si.Bind(bindingID, resource, params)
	if err != nil {
		return bosherr.WrapError(err, "Binding to instance")
	}

	info := boshtbl.Table{
		Header: []boshtbl.Header{
			boshtbl.NewHeader("ID"),
			boshtbl.NewHeader("Service Instance ID"),
			boshtbl.NewHeader("Credentials"),
			boshtbl.NewHeader("Syslog Drain URL"),
			boshtbl.NewHeader("Route Service URL"),
			boshtbl.NewHeader("Volume Mounts"),
		},
		Rows: [][]boshtbl.Value{
			{
				boshtbl.NewValueString(sb.ID()),
				boshtbl.NewValueString(sb.ServiceInstance().ID()),
				boshtbl.NewValueInterface(sb.Credentials()),
				boshtbl.NewValueString(sb.SyslogDrainURL()),
				boshtbl.NewValueString(sb.RouteServiceURL()),
				boshtbl.NewValueInterface(sb.VolumeMounts()),
			},
		},
		Transpose: true,
	}

	c.ui.PrintTable(info)

	return nil
}
