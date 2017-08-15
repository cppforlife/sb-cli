package cmd

import (
	"fmt"

	boshui "github.com/cloudfoundry/bosh-cli/ui"
)

type ServiceInstancesCmd struct {
	siFactory ServiceInstanceFactory
	ui        boshui.UI
}

func NewServiceInstancesCmd(siFactory ServiceInstanceFactory, ui boshui.UI) ServiceInstancesCmd {
	return ServiceInstancesCmd{siFactory, ui}
}

func (c ServiceInstancesCmd) Run() error {
	return fmt.Errorf("Not implemented")
}
