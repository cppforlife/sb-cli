package cmd

import (
	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

type planView struct {
	Plan osb.Plan
}

func (v planView) InstanceCreateSchema() interface{} {
	schemas := v.Plan.AlphaParameterSchemas
	if schemas == nil {
		return nil
	}
	si := schemas.ServiceInstances
	if si == nil {
		return nil
	}
	params := si.Create
	if params == nil {
		return nil
	}
	return params.Parameters
}

func (v planView) InstanceUpdateSchema() interface{} {
	schemas := v.Plan.AlphaParameterSchemas
	if schemas == nil {
		return nil
	}
	si := schemas.ServiceInstances
	if si == nil {
		return nil
	}
	params := si.Update
	if params == nil {
		return nil
	}
	return params.Parameters
}

func (v planView) BindingCreateSchema() interface{} {
	schemas := v.Plan.AlphaParameterSchemas
	if schemas == nil {
		return nil
	}
	si := schemas.ServiceBindings
	if si == nil {
		return nil
	}
	params := si.Create
	if params == nil {
		return nil
	}
	return params.Parameters
}
