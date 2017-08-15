package cmd

import (
	"time"

	"code.cloudfoundry.org/clock"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

type ServiceInstanceFinder struct {
	Name            string
	ServiceName     string
	ServicePlanName string
}

type ServiceInstanceFactory struct {
	client osb.Client
	clock  clock.Clock
	ui     boshui.UI
}

func NewServiceInstanceFactory(client osb.Client, clock clock.Clock, ui boshui.UI) ServiceInstanceFactory {
	return ServiceInstanceFactory{client, clock, ui}
}

func (f ServiceInstanceFactory) New(id ServiceInstanceFinder) *ServiceInstance {
	return &ServiceInstance{id, f.client, f.clock, f.ui, ""}
}

type ServiceInstance struct {
	id ServiceInstanceFinder

	client osb.Client
	clock  clock.Clock
	ui     boshui.UI

	dashboardURL string
}

func (c ServiceInstance) Name() string            { return c.id.Name }
func (c ServiceInstance) ServiceName() string     { return c.id.ServiceName }
func (c ServiceInstance) ServicePlanName() string { return c.id.ServicePlanName }
func (c ServiceInstance) DashboardURL() string    { return c.dashboardURL }

func (c *ServiceInstance) Provision(params map[string]interface{}) error {
	req := &osb.ProvisionRequest{
		InstanceID: c.id.Name,
		ServiceID:  c.id.ServiceName,
		PlanID:     c.id.ServicePlanName,

		AcceptsIncomplete: true,

		OrganizationGUID: "unused", // todo
		SpaceGUID:        "unused", // todo

		Parameters: params,
	}

	resp, err := c.client.ProvisionInstance(req)
	if err != nil {
		return bosherr.WrapError(err, "Provisioning instance")
	}

	if resp.DashboardURL != nil {
		c.dashboardURL = *resp.DashboardURL
	}

	if resp.Async {
		opReq := &osb.LastOperationRequest{
			InstanceID:   req.InstanceID,
			ServiceID:    &req.ServiceID,
			PlanID:       &req.PlanID,
			OperationKey: resp.OperationKey,
		}

		err := c.waitForCompletion(opReq, false)
		if err != nil {
			return bosherr.WrapError(err, "Polling instance")
		}
	}

	return nil
}

func (c ServiceInstance) Deprovision() error {
	req := &osb.DeprovisionRequest{
		InstanceID: c.id.Name,
		ServiceID:  c.id.ServiceName,
		PlanID:     c.id.ServicePlanName,

		AcceptsIncomplete: true,
	}

	resp, err := c.client.DeprovisionInstance(req)
	if err != nil {
		return bosherr.WrapError(err, "Deprovisioning instance")
	}

	if resp.Async {
		opReq := &osb.LastOperationRequest{
			InstanceID:   req.InstanceID,
			ServiceID:    &req.ServiceID,
			PlanID:       &req.PlanID,
			OperationKey: resp.OperationKey,
		}

		err := c.waitForCompletion(opReq, true)
		if err != nil {
			return bosherr.WrapError(err, "Polling instance")
		}
	}

	return nil
}

func (c ServiceInstance) waitForCompletion(req *osb.LastOperationRequest, deleting bool) error {
	var lastDesc string

	for {
		resp, err := c.client.PollLastOperation(req)
		if err != nil {
			if osb.IsGoneError(err) && deleting {
				return nil
			}

			return bosherr.WrapError(err, "Polling last operation")
		}

		switch resp.State {
		case osb.StateInProgress:
			if resp.Description != nil {
				if lastDesc != *resp.Description {
					c.ui.PrintLinef("In progress state: %s", *resp.Description)
				} else {
					c.ui.BeginLinef(".")
				}
				lastDesc = *resp.Description
			}

		case osb.StateSucceeded:
			return nil

		case osb.StateFailed:
			description := ""
			if resp.Description != nil {
				description = *resp.Description
			}

			return bosherr.Errorf("Failed state: %q", description)

		default:
			return bosherr.Errorf("Invalid state: %q", resp.State)
		}

		c.clock.Sleep(500 * time.Millisecond)
	}
}
