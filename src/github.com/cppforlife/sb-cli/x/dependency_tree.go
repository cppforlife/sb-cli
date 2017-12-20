package x

import (
	"fmt"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type DependencyTree struct {
	cli CLIImpl
	fs  boshsys.FileSystem
}

func NewDependencyTree(cli CLIImpl, fs boshsys.FileSystem) DependencyTree {
	return DependencyTree{cli, fs}
}

func (d DependencyTree) Build(deploymentGroupName, manifestPath string) (*Dependency, error) {
	if len(deploymentGroupName) == 0 {
		return nil, fmt.Errorf("Expected deployment group name to be non-empty")
	}

	if len(manifestPath) == 0 {
		return nil, fmt.Errorf("Expected manifest path to be non-empty")
	}

	rootItem := &FSDependencyTreeItem{
		Name:         deploymentGroupName,
		ManifestPath: manifestPath,

		cli: d.cli,
		fs:  d.fs,
	}

	err := rootItem.Build()
	if err != nil {
		return nil, err
	}

	return rootItem.AsDependency(), nil
}
