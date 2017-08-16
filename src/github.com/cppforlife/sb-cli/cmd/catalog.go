package cmd

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

type Catalog struct {
	client osb.Client
}

func (c Catalog) Services() ([]osb.Service, error) {
	resp, err := c.client.GetCatalog()
	if err != nil {
		return nil, bosherr.WrapError(err, "Fetching catalog")
	}

	return resp.Services, nil
}
