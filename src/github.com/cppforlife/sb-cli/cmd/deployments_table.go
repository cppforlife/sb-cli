package cmd

import (
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"

	boshx "github.com/cppforlife/sb-cli/x"
)

type DeploymentsTable struct {
	dep *boshx.Dependency
	ui  boshui.UI
}

func (t DeploymentsTable) Print() {
	table := boshtbl.Table{
		Content: "Deployments",

		Header: []boshtbl.Header{
			boshtbl.NewHeader("Name"),
			boshtbl.NewHeader("Dependencies"),
		},
	}

	allDeps := map[*boshx.Dependency]struct{}{}

	t.collectAllDeps(t.dep, allDeps)

	for dep, _ := range allDeps {
		depNames := []string{}

		for _, dep1 := range dep.Dependencies {
			depNames = append(depNames, dep1.Name)
		}

		table.Rows = append(table.Rows, []boshtbl.Value{
			boshtbl.NewValueString(dep.Name),
			boshtbl.NewValueStrings(depNames),
		})
	}

	t.ui.PrintTable(table)
}

func (t DeploymentsTable) collectAllDeps(dep *boshx.Dependency, allDeps map[*boshx.Dependency]struct{}) {
	allDeps[dep] = struct{}{}
	for _, dep1 := range dep.Dependencies {
		t.collectAllDeps(dep1, allDeps)
	}
}
