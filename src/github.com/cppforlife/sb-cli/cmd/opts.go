package cmd

import (
	"time"
	// boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
)

type SBOpts struct {
	VersionOpt func() error `long:"version" short:"v" description:"Show CLI version"`

	JSONOpt    bool `long:"json"     description:"Output as JSON"`
	TTYOpt     bool `long:"tty"      description:"Force TTY-like output"`
	NoColorOpt bool `long:"no-color" description:"Toggle colorized output"`

	Help HelpOpts `command:"help" description:"Show this help message"`

	URL      string        `long:"broker-url"      value-name:"URL"      description:"Broker URL"                          env:"SB_BROKER_URL"`
	Username string        `long:"broker-username" value-name:"USERNAME" description:"Broker username"                     env:"SB_BROKER_USERNAME"`
	Password string        `long:"broker-password" value-name:"PASSWORD" description:"Broker password"                     env:"SB_BROKER_PASSWORD"`
	CACert   FileBytesArg  `long:"broker-ca-cert"  value-name:"PATH"     description:"Path to a file with CA certificates" env:"SB_BROKER_CA_CERT"`
	Timeout  time.Duration `long:"broker-timeout"  value-name:"DURATION" description:"Timeout for individual HTTP requests" default:"30s"`

	Services ServicesOpts `command:"services" alias:"ss" description:"List services"`

	ServiceInstances      ServiceInstancesOpts      `command:"service-instances"       alias:"sis" description:"List service instances"`
	CreateServiceInstance CreateServiceInstanceOpts `command:"create-service-instance" alias:"csi" description:"Create service instance"`
	UpdateServiceInstance UpdateServiceInstanceOpts `command:"update-service-instance" alias:"usi" description:"Update service instance"`
	DeleteServiceInstance DeleteServiceInstanceOpts `command:"delete-service-instance" alias:"dsi" description:"Delete service instance"`

	CreateServiceBinding CreateServiceBindingOpts `command:"create-service-binding" alias:"csb" description:"Create service binding"`
	DeleteServiceBinding DeleteServiceBindingOpts `command:"delete-service-binding" alias:"dsb" description:"Delete service binding"`
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

	ParamFlags

	cmd
}

type CreateServiceInstanceArgs struct {
	ID string `positional-arg-name:"SERVICE-INSTANCE-ID" description:"Service instance ID"`
}

type UpdateServiceInstanceOpts struct {
	Args UpdateServiceInstanceArgs `positional-args:"true" required:"true"`

	ServiceID     string `long:"service"      value-name:"SERVICE-ID" description:"Service ID"`
	ServicePlanID string `long:"service-plan" value-name:"PLAN-ID"    description:"Service plan ID"`

	ParamFlags

	cmd
}

type UpdateServiceInstanceArgs struct {
	ID string `positional-arg-name:"SERVICE-INSTANCE-ID" description:"Service instance ID"`
}

type DeleteServiceInstanceOpts struct {
	Args DeleteServiceInstanceArgs `positional-args:"true" required:"true"`

	ServiceID     string `long:"service"      value-name:"SERVICE-ID" description:"Service ID"`
	ServicePlanID string `long:"service-plan" value-name:"PLAN-ID"    description:"Service plan ID"`

	cmd
}

type DeleteServiceInstanceArgs struct {
	ID string `positional-arg-name:"SERVICE-INSTANCE-ID" description:"Service instance ID"`
}

type CreateServiceBindingOpts struct {
	Args CreateServiceBindingArgs `positional-args:"true" required:"true"`

	ID string `long:"id" value-name:"SERVICE-BINDING-ID" description:"Service binding ID (automatically generated if not provided)"`

	ServiceID     string `long:"service"      value-name:"SERVICE-ID" description:"Service ID"`
	ServicePlanID string `long:"service-plan" value-name:"PLAN-ID"    description:"Service plan ID"`

	Resource FileBytesArg `long:"resource" value-name:"PATH" description:"Path to a YAML file with resource definition"`

	ParamFlags

	Timeout time.Duration `long:"timeout" value-name:"DURATION" description:"Timeout for binding operation" default:"90s"`

	cmd
}

type CreateServiceBindingArgs struct {
	ServiceInstanceID string `positional-arg-name:"SERVICE-INSTANCE-ID" description:"Service instance ID"`
}

type DeleteServiceBindingOpts struct {
	Args DeleteServiceBindingArgs `positional-args:"true" required:"true"`

	ServiceID     string `long:"service"      value-name:"SERVICE-ID" description:"Service ID"`
	ServicePlanID string `long:"service-plan" value-name:"PLAN-ID"    description:"Service plan ID"`

	Timeout time.Duration `long:"timeout" value-name:"DURATION" description:"Timeout for unbinding operation" default:"90s"`

	cmd
}

type DeleteServiceBindingArgs struct {
	ID                string `positional-arg-name:"SERVICE-BINDING-ID"  description:"Service binding ID"`
	ServiceInstanceID string `positional-arg-name:"SERVICE-INSTANCE-ID" description:"Service instance ID"`
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
