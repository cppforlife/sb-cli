package cmd

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

func (c *ServiceInstance) Bind(bindingID string, resource *osb.BindResource, params map[string]interface{}) (ServiceBinding, error) {
	req := &osb.BindRequest{
		BindingID:  bindingID,
		InstanceID: c.id.ID,

		ServiceID: c.id.ServiceID,
		PlanID:    c.id.ServicePlanID,

		BindResource: resource,
		Parameters:   params,
	}

	resp, err := c.client.Bind(req)
	if err != nil {
		return ServiceBinding{}, bosherr.WrapError(err, "Binding to instance")
	}

	return ServiceBinding{
		id:              req.BindingID,
		serviceInstance: c,

		credentials: resp.Credentials,

		syslogDrainURL:  resp.SyslogDrainURL,
		routeServiceURL: resp.RouteServiceURL,
		volumeMounts:    resp.VolumeMounts,
	}, nil
}

func (c *ServiceInstance) Unbind(bindingID string) error {
	req := &osb.UnbindRequest{
		BindingID:  bindingID,
		InstanceID: c.id.ID,

		ServiceID: c.id.ServiceID,
		PlanID:    c.id.ServicePlanID,
	}

	_, err := c.client.Unbind(req)
	if err != nil {
		return bosherr.WrapError(err, "Unbinding from instance")
	}

	return nil
}
