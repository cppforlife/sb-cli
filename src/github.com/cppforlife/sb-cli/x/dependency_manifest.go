package x

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"gopkg.in/yaml.v2"
)

type deploymentManifest struct {
	Dependencies []DependencyManifest
}

type DependencyManifest struct {
	Name string
	URL  string
}

func NewDependenciesManifestFromBytes(bytes []byte) ([]DependencyManifest, error) {
	manifest := deploymentManifest{}

	err := yaml.Unmarshal(bytes, &manifest)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Unmarshaling dependencies manifest")
	}

	return manifest.Dependencies, nil
}
