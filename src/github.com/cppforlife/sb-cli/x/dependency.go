package x

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

type Dependency struct {
	Name     string
	Manifest []byte

	DeploymentName string

	// todo
	// Variables []Variable
	// OpsFiles  []OpsFile

	Dependencies []*Dependency

	cli CLIImpl
}

func (d *Dependency) manifestWithAdjustedName() ([]byte, error) {
	var manifest map[interface{}]interface{}

	err := yaml.Unmarshal(d.Manifest, &manifest)
	if err != nil {
		return nil, fmt.Errorf("Unmarshaling manifest: %s", err)
	}

	manifest["name"] = d.DeploymentName

	// Update all cross deployment links with their dependency deployment names
	d.typedArrayOfMaps(manifest["instance_groups"], func(ig map[interface{}]interface{}) {
		d.typedArrayOfMaps(ig["jobs"], func(job map[interface{}]interface{}) {
			d.typedMapOfMaps(job["consumes"], func(link map[interface{}]interface{}) {
				if deploymentName, ok := link["deployment"].(string); ok {
					for _, dep := range d.Dependencies {
						if dep.Name == deploymentName {
							link["deployment"] = dep.DeploymentName
						}
					}
				}
			})
		})
	})

	manifestBytes, err := yaml.Marshal(manifest)
	if err != nil {
		return nil, fmt.Errorf("Marshaling manifest: %s", err)
	}

	return manifestBytes, nil
}

func (Dependency) typedArrayOfMaps(obj interface{}, iter func(map[interface{}]interface{})) {
	if array, ok := obj.([]interface{}); ok {
		for _, item := range array {
			if typedItem, ok := item.(map[interface{}]interface{}); ok {
				iter(typedItem)
			}
		}
	}
}

func (Dependency) typedMapOfMaps(obj interface{}, iter func(map[interface{}]interface{})) {
	if map_, ok := obj.(map[interface{}]interface{}); ok {
		for _, item := range map_ {
			if typedItem, ok := item.(map[interface{}]interface{}); ok {
				iter(typedItem)
			}
		}
	}
}

// todo get rid of nasties
func (d *Dependency) applicableVars(vars []Variable) ([]Variable, error) {
	var applicable []Variable
	var unknownDepNames []string

	for _, v1 := range vars {
		if len(v1.DependencyNames) == 1 {
			if v1.DependencyNames[0] == d.Name {
				applicable = append(applicable, v1)
			} else {
				unknownDepNames = append(unknownDepNames, v1.DependencyNames[0])
			}
		}
	}

	if len(unknownDepNames) > 0 {
		return nil, fmt.Errorf("Expected dependency names '%s' to exist for variables", strings.Join(unknownDepNames, "', '"))
	}

	return applicable, nil
}

func (d *Dependency) applicableVars2(vars []Variable) []Variable {
	var applicable []Variable

	for _, v1 := range vars {
		if len(v1.DependencyNames) > 1 && v1.DependencyNames[0] == d.Name {
			v1.DependencyNames = v1.DependencyNames[1:]
			applicable = append(applicable, v1)
		}
	}

	return applicable
}

func (d *Dependency) applicableOps(opsFiles []OpsFile) ([]OpsFile, error) {
	var applicable []OpsFile
	var unknownDepNames []string

	for _, op := range opsFiles {
		if len(op.DependencyNames) == 1 {
			if op.DependencyNames[0] == d.Name {
				applicable = append(applicable, op)
			} else {
				unknownDepNames = append(unknownDepNames, op.DependencyNames[0])
			}
		}
	}

	if len(unknownDepNames) > 0 {
		return nil, fmt.Errorf("Expected dependency names '%s' to exist for ops files", strings.Join(unknownDepNames, "', '"))
	}

	return applicable, nil
}

func (d *Dependency) applicableOps2(opsFiles []OpsFile) []OpsFile {
	var applicable []OpsFile

	for _, op := range opsFiles {
		if len(op.DependencyNames) > 1 && op.DependencyNames[0] == d.Name {
			op.DependencyNames = op.DependencyNames[1:]
			applicable = append(applicable, op)
		}
	}

	return applicable
}
