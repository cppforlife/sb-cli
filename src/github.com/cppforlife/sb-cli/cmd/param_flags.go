package cmd

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"gopkg.in/yaml.v2"
)

type ParamFlags struct {
	Params     []ParamArg   `long:"param"  short:"p" value-name:"NAME=VALUE" description:"Set parameter (value is interpreted as YAML)"`
	ParamsFile FileBytesArg `long:"params"           value-name:"PATH"       description:"Path to a YAML file with params"`
}

func (f ParamFlags) AsParams() (map[string]interface{}, error) {
	params := map[string]interface{}{}

	if len(f.ParamsFile.Bytes) > 0 {
		err := yaml.Unmarshal(f.ParamsFile.Bytes, &params)
		if err != nil {
			return nil, bosherr.WrapError(err, "Unmarshaling service instance params")
		}
	}

	for _, p := range f.Params {
		params[p.Name] = p.Value
	}

	return params, nil
}
