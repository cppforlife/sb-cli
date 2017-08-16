package cmd

import (
	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
)

type SBOpts struct {
	VersionOpt func() error `long:"version" short:"v" description:"Show CLI version"`

	JSONOpt    bool `long:"json"     description:"Output as JSON"`
	TTYOpt     bool `long:"tty"      description:"Force TTY-like output"`
	NoColorOpt bool `long:"no-color" description:"Toggle colorized output"`

	Help HelpOpts `command:"help" description:"Show this help message"`

	URL      string               `long:"broker-url"      value-name:"URL"      description:"Broker URL"                          env:"SB_BROKER_URL"`
	Username string               `long:"broker-username" value-name:"USERNAME" description:"Broker username"                     env:"SB_BROKER_USERNAME"`
	Password string               `long:"broker-password" value-name:"PASSWORD" description:"Broker password"                     env:"SB_BROKER_PASSWORD"`
	CACert   boshcmd.FileBytesArg `long:"broker-ca-cert"  value-name:"PATH"     description:"Path to a file with CA certificates" env:"SB_BROKER_CA_CERT"`

	Services ServicesOpts `command:"services" alias:"ss" description:"List services"`

	ServiceInstances      ServiceInstancesOpts      `command:"service-instances"       alias:"sis" description:"List service instances"`
	CreateServiceInstance CreateServiceInstanceOpts `command:"create-service-instance" alias:"cis" description:"Create service instance"`
	DeleteServiceInstance DeleteServiceInstanceOpts `command:"delete-service-instance" alias:"dis" description:"Delete service instance"`

	CreateServiceBinding CreateServiceBindingOpts `command:"create-service-binding" alias:"csb" description:""`
	DeleteServiceBinding DeleteServiceBindingOpts `command:"delete-service-binding" alias:"dsb" description:""`
}

type HelpOpts struct {
	cmd
}

type ServicesOpts struct {
	cmd
}

type ServiceInstancesOpts struct {
	cmd
}

type CreateServiceInstanceOpts struct {
	Args CreateServiceInstanceArgs `positional-args:"true" required:"true"`

	ServiceID     string `long:"service"      value-name:"SERVICE-ID" description:"Service ID"`
	ServicePlanID string `long:"service-plan" value-name:"PLAN-ID"    description:"Service plan ID"`

	Params boshcmd.FileBytesArg `long:"params" value-name:"PATH" description:"Path to a file with options"`
	cmd
}

type CreateServiceInstanceArgs struct {
	ID string `positional-arg-name:"ID" description:"Service instance ID"`
}

type DeleteServiceInstanceOpts struct {
	Args DeleteServiceInstanceArgs `positional-args:"true" required:"true"`

	ServiceID     string `long:"service"      value-name:"SERVICE-ID" description:"Service ID"`
	ServicePlanID string `long:"service-plan" value-name:"PLAN-ID"    description:"Service plan ID"`

	cmd
}

type DeleteServiceInstanceArgs struct {
	ID string `positional-arg-name:"ID" description:"Service instance ID"`
}

type CreateServiceBindingOpts struct {
	cmd
}

type DeleteServiceBindingOpts struct {
	cmd
}

// MessageOpts is used for version and help flags
type MessageOpts struct {
	Message string
}

type cmd struct{}

// Execute is necessary for each command to be goflags.Commander
func (c cmd) Execute(_ []string) error {
	panic("Unreachable")
}
