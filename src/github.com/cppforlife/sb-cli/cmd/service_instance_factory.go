package cmd

import (
	"code.cloudfoundry.org/clock"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

type ServiceInstanceFinder struct {
	ID            string
	ServiceID     string
	ServicePlanID string
}

type ServiceInstanceFactory struct {
	client  osb.Client
	catalog Catalog
	clock   clock.Clock
	ui      boshui.UI
}

func NewServiceInstanceFactory(client osb.Client, catalog Catalog, clock clock.Clock, ui boshui.UI) ServiceInstanceFactory {
	return ServiceInstanceFactory{client, catalog, clock, ui}
}

func (f ServiceInstanceFactory) New(id ServiceInstanceFinder) (*ServiceInstance, error) {
	id, err := f.backfillID(id)
	return &ServiceInstance{id, f.client, f.clock, f.ui, ""}, err
}

func (f ServiceInstanceFactory) backfillID(id ServiceInstanceFinder) (ServiceInstanceFinder, error) {
	if len(id.ServiceID) == 0 || len(id.ServicePlanID) == 0 {
		service, err := f.findService(id.ServiceID)
		if err != nil {
			return id, err
		}

		id.ServiceID = service.ID

		if len(id.ServicePlanID) == 0 {
			if len(service.Plans) == 0 {
				return id, bosherr.Errorf("Expected to find at least one service plan for service ID '%s'", service.ID)
			}

			id.ServicePlanID = service.Plans[0].ID
		}
	}

	return id, nil
}

func (f ServiceInstanceFactory) findService(serviceID string) (osb.Service, error) {
	services, err := f.catalog.Services()
	if err != nil {
		return osb.Service{}, bosherr.WrapError(err, "Fetching catalog")
	}

	if len(serviceID) == 0 {
		if len(services) == 0 {
			return osb.Service{}, bosherr.Errorf("Expected to find at least one service in the catalog")
		}

		return services[0], nil
	}

	for _, serv := range services {
		if serv.ID == serviceID {
			return serv, nil
		}
	}

	return osb.Service{}, bosherr.Errorf("Expected to find service ID '%s'", serviceID)
}
