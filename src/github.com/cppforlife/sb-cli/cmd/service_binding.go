package cmd

type ServiceBinding struct {
	id              string
	serviceInstance *ServiceInstance

	credentials map[string]interface{}

	syslogDrainURL  *string
	routeServiceURL *string
	volumeMounts    []interface{}
}

func (b ServiceBinding) ID() string                        { return b.id }
func (b ServiceBinding) ServiceInstance() *ServiceInstance { return b.serviceInstance }

func (b ServiceBinding) Credentials() map[string]interface{} { return b.credentials }

func (b ServiceBinding) SyslogDrainURL() string {
	if b.syslogDrainURL == nil {
		return ""
	}
	return *b.syslogDrainURL
}

func (b ServiceBinding) RouteServiceURL() string {
	if b.routeServiceURL == nil {
		return ""
	}
	return *b.routeServiceURL
}

func (b ServiceBinding) VolumeMounts() []interface{} { return b.volumeMounts }
