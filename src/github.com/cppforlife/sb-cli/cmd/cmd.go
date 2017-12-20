package cmd

import (
	"errors"
	"fmt"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

type Cmd struct {
	SBOpts SBOpts
	Opts   interface{}

	deps BasicDeps
}

func NewCmd(sbOpts SBOpts, opts interface{}, deps BasicDeps) Cmd {
	return Cmd{sbOpts, opts, deps}
}

type cmdConveniencePanic struct {
	Err error
}

func (c Cmd) Execute() (cmdErr error) {
	// Catch convenience panics from panicIfErr
	defer func() {
		if r := recover(); r != nil {
			if cp, ok := r.(cmdConveniencePanic); ok {
				cmdErr = cp.Err
			} else {
				panic(r)
			}
		}
	}()

	c.configureUI()
	c.configureFS()

	deps := c.deps
	var client osb.Client

	switch c.Opts.(type) {
	case *XDeployOpts:
	case *XInterpolateOpts:
	case *MessageOpts:
	default:
		client = c.osbClient()
	}

	catalog := Catalog{client}
	siFactory := NewServiceInstanceFactory(client, catalog, deps.Time, deps.UI)

	switch opts := c.Opts.(type) {
	case *ServicesOpts:
		return NewServicesCmd(catalog, deps.UI).Run()

	case *ServiceInstancesOpts:
		return NewServiceInstancesCmd(siFactory, deps.UI).Run()

	case *CreateServiceInstanceOpts:
		return NewCreateServiceInstanceCmd(siFactory, deps.UI).Run(*opts)

	case *UpdateServiceInstanceOpts:
		return NewUpdateServiceInstanceCmd(siFactory).Run(*opts)

	case *DeleteServiceInstanceOpts:
		return NewDeleteServiceInstanceCmd(siFactory, deps.UI).Run(*opts)

	case *CreateServiceBindingOpts:
		return NewCreateServiceBindingCmd(siFactory, deps.UUIDGen, deps.UI).Run(*opts)

	case *DeleteServiceBindingOpts:
		return NewDeleteServiceBindingCmd(siFactory).Run(*opts)

	case *XDeployOpts:
		return NewXDeployCmd(deps.UI, deps.FS).Run(*opts)

	case *XDeleteOpts:
		return NewXDeleteCmd(deps.UI, deps.FS).Run(*opts)

	case *XInterpolateOpts:
		return NewXInterpolateCmd(deps.UI, deps.FS).Run(*opts)

	case *MessageOpts:
		deps.UI.PrintBlock(opts.Message)
		return nil

	default:
		return fmt.Errorf("Unhandled command: %#v", c.Opts)
	}
}

func (c Cmd) configureUI() {
	c.deps.UI.EnableTTY(c.SBOpts.TTYOpt)

	if !c.SBOpts.NoColorOpt {
		c.deps.UI.EnableColor()
	}

	if c.SBOpts.JSONOpt {
		c.deps.UI.EnableJSON()
	}

	c.deps.UI.EnableNonInteractive()
}

func (c Cmd) osbClient() osb.Client {
	if len(c.SBOpts.URL) == 0 {
		c.panicIfErr(errors.New("Expected service broker URL"))
	}

	if len(c.SBOpts.Username) == 0 {
		c.panicIfErr(errors.New("Expected service broker username"))
	}

	if len(c.SBOpts.Password) == 0 {
		c.panicIfErr(errors.New("Expected service broker password"))
	}

	config := osb.DefaultClientConfiguration()
	config.URL = c.SBOpts.URL
	config.CAData = c.SBOpts.CACert.Bytes
	config.Insecure = false
	config.TimeoutSeconds = int(c.SBOpts.Timeout) / 1000000000

	config.AuthConfig = &osb.AuthConfig{
		BasicAuthConfig: &osb.BasicAuthConfig{
			Username: c.SBOpts.Username,
			Password: c.SBOpts.Password,
		},
	}

	client, err := osb.NewClient(config)
	c.panicIfErr(err)

	if len(config.URL) > 0 {
		c.deps.UI.PrintLinef("Using service broker '%s' as user '%s'",
			config.URL, config.AuthConfig.BasicAuthConfig.Username)
	}

	return client
}

func (c Cmd) configureFS() {
	tmpDirPath, err := c.deps.FS.ExpandPath("~/.sb-cli/tmp")
	c.panicIfErr(err)

	err = c.deps.FS.ChangeTempRoot(tmpDirPath)
	c.panicIfErr(err)
}

func (c Cmd) panicIfErr(err error) {
	if err != nil {
		panic(cmdConveniencePanic{err})
	}
}
