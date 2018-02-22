package cmd

import (
	"bytes"
	"strings"

	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

type ServicesCmd struct {
	catalog Catalog
	ui      boshui.UI
}

func NewServicesCmd(catalog Catalog, ui boshui.UI) ServicesCmd {
	return ServicesCmd{catalog, ui}
}

func (c ServicesCmd) Run() error {
	table := boshtbl.Table{
		Content: "services",

		Header: []boshtbl.Header{
			boshtbl.NewHeader("ID"),
			boshtbl.NewHeader("Name"),
			boshtbl.NewHeader("Plans"),
		},

		Transpose: true,
	}

	services, err := c.catalog.Services()
	if err != nil {
		return err
	}

	for _, serv := range services {
		table.Rows = append(table.Rows, []boshtbl.Value{
			boshtbl.NewValueString(serv.ID),
			boshtbl.NewValueString(serv.Name),
			boshtbl.NewValueString(c.plansTableStr(serv.Plans)),
		})
	}

	c.ui.PrintTable(table)

	return nil
}

func (c ServicesCmd) plansTableStr(plans []osb.Plan) string {
	table := boshtbl.Table{
		Header: []boshtbl.Header{
			boshtbl.NewHeader("ID"),
			boshtbl.NewHeader("Name"),
			boshtbl.NewHeader("Bindable"),
		},

		// Transpose: true,
	}

	for _, plan := range plans {
		table.Rows = append(table.Rows, []boshtbl.Value{
			boshtbl.NewValueString(plan.ID),
			boshtbl.NewValueString(plan.Name),
			ValueBoolOptional{plan.Bindable},
		})
	}

	buf := bytes.NewBufferString("")
	table.Print(buf)

	return strings.Trim(buf.String(), "\n")
}
