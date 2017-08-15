package cmd

import (
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
	siFactory := NewServiceInstanceFactory(c.osbClient(), deps.Time, deps.UI)

	switch opts := c.Opts.(type) {
	case *ServiceInstancesOpts:
		return NewServiceInstancesCmd(siFactory, deps.UI).Run()

	case *CreateServiceInstanceOpts:
		return NewCreateServiceInstanceCmd(siFactory, deps.UI).Run(*opts)

	case *DeleteServiceInstanceOpts:
		return NewDeleteServiceInstanceCmd(siFactory, deps.UI).Run(*opts)

	case *CreateServiceBindingOpts:
		return fmt.Errorf("Not implemented yet")

	case *DeleteServiceBindingOpts:
		return fmt.Errorf("Not implemented yet")

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
	config := osb.DefaultClientConfiguration()
	config.URL = c.SBOpts.URL
	config.CAData = c.SBOpts.CACert.Bytes
	config.Insecure = false

	config.AuthConfig = &osb.AuthConfig{
		BasicAuthConfig: &osb.BasicAuthConfig{
			Username: c.SBOpts.Username,
			Password: c.SBOpts.Password,
		},
	}

	client, err := osb.NewClient(config)
	c.panicIfErr(err)

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
