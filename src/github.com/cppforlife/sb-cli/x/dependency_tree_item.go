package x

import (
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type FSDependencyTreeItem struct {
	Name         string
	ManifestPath string

	Manifest []byte

	ParentItem *FSDependencyTreeItem
	Items      []*FSDependencyTreeItem

	cli CLIImpl
	fs  boshsys.FileSystem
}

func (d *FSDependencyTreeItem) Build() error {
	manifest, err := d.fs.ReadFile(d.ManifestPath)
	if err != nil {
		return bosherr.WrapErrorf(err, "Reading base manifest")
	}

	d.Manifest = manifest

	// todo support ops files for changing depedencies
	depsManifest, err := NewDependenciesManifestFromBytes(manifest)
	if err != nil {
		return err
	}

	for _, man := range depsManifest {
		item := &FSDependencyTreeItem{
			Name:         man.Name,
			ManifestPath: strings.TrimPrefix(man.URL, "file://"), // todo support more schemas?

			ParentItem: d,

			cli: d.cli,
			fs:  d.fs,
		}

		err := item.Build()
		if err != nil {
			return err
		}

		d.Items = append(d.Items, item)
	}

	return nil
}

func (d FSDependencyTreeItem) names() []string {
	if d.ParentItem != nil {
		return append(d.ParentItem.names(), d.Name)
	}
	return []string{d.Name}
}

func (d FSDependencyTreeItem) AsDependency() *Dependency {
	deps := []*Dependency{}

	for _, item := range d.Items {
		deps = append(deps, item.AsDependency())
	}

	return &Dependency{
		Name:     d.Name,
		Manifest: d.Manifest,

		DeploymentName: strings.Join(d.names(), "-"),

		Dependencies: deps,

		cli: d.cli,
	}
}
