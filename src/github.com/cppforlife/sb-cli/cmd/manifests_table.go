package cmd

import (
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"

	boshx "github.com/cppforlife/sb-cli/x"
)

type ManifestsTable struct {
	manifests []boshx.Manifest
	ui        boshui.UI
}

func (t ManifestsTable) Print() {
	table := boshtbl.Table{
		Content: "Manifests",

		Header: []boshtbl.Header{
			boshtbl.NewHeader("Name"),
			boshtbl.NewHeader("Dependencies"),
			boshtbl.NewHeader("Manifest"),
		},

		// todo Transpose: true,
	}

	for _, manifest := range t.manifests {
		depNames := []string{}

		for _, dep1 := range manifest.Dependency.Dependencies {
			depNames = append(depNames, dep1.Name)
		}

		table.Rows = append(table.Rows, []boshtbl.Value{
			boshtbl.NewValueString(manifest.Dependency.Name),
			boshtbl.NewValueStrings(depNames),
			boshtbl.NewValueString(string(manifest.Bytes)),
		})
	}

	t.ui.PrintTable(table)
}
