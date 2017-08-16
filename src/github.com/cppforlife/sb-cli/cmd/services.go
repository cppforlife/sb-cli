package cmd

import (
	"bytes"
	"fmt"
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
			boshtbl.NewHeader("Description"),
			boshtbl.NewHeader("Tags"),
			boshtbl.NewHeader("Requires"),
			boshtbl.NewHeader("Bindable"),
			boshtbl.NewHeader("Metadata"),
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
			boshtbl.NewValueString(serv.Description),
			boshtbl.NewValueStrings(serv.Tags),
			boshtbl.NewValueStrings(serv.Requires),
			boshtbl.NewValueBool(serv.Bindable),
			boshtbl.NewValueInterface(serv.Metadata),
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
			boshtbl.NewHeader("Description"),
			boshtbl.NewHeader("Free"),
			boshtbl.NewHeader("Bindable"),
			boshtbl.NewHeader("Metadata"),
			boshtbl.NewHeader("Instance Create Schema"),
			boshtbl.NewHeader("Instance Update Schema"),
			boshtbl.NewHeader("Binding Create Schema"),
		},

		Transpose: true,
	}

	for _, plan := range plans {
		view := planView{plan}

		table.Rows = append(table.Rows, []boshtbl.Value{
			boshtbl.NewValueString(plan.ID),
			boshtbl.NewValueString(plan.Name),
			boshtbl.NewValueString(plan.Description),
			ValueBoolOptional{plan.Free},
			ValueBoolOptional{plan.Bindable},
			boshtbl.NewValueInterface(plan.Metadata),
			boshtbl.NewValueInterface(view.InstanceCreateSchema()),
			boshtbl.NewValueInterface(view.InstanceUpdateSchema()),
			boshtbl.NewValueInterface(view.BindingCreateSchema()),
		})
	}

	buf := bytes.NewBufferString("")
	table.Print(buf)

	return strings.Trim(buf.String(), "\n")
}

type ValueBoolOptional struct{ B *bool }

func (t ValueBoolOptional) String() string {
	if t.B == nil {
		return ""
	}
	return fmt.Sprintf("%t", *t.B)
}

func (t ValueBoolOptional) Value() boshtbl.Value            { return t }
func (t ValueBoolOptional) Compare(other boshtbl.Value) int { panic("unreachable") }
