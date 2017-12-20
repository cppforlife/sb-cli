package x

import (
	"fmt"
	"strings"
)

type Variable struct {
	DependencyNames []string
	Name            string
	Value           string
}

func NewVarsFromArgs(args []string) ([]Variable, []string, error) {
	var vars []Variable
	var unusedArgs []string

	nextIsVar := false

	for _, arg := range args {
		if arg == "-v" {
			nextIsVar = true
		} else if nextIsVar {
			nextIsVar = false

			pieces := strings.SplitN(arg, "=", 2)
			if len(pieces) != 2 {
				return nil, nil, fmt.Errorf("Expected variable '%s' to have format 'name=value'", arg)
			}

			name := pieces[0]
			val := pieces[1]

			depNames := strings.Split(name, ":")
			name = depNames[len(depNames)-1]
			depNames = depNames[:len(depNames)-1]

			vars = append(vars, Variable{DependencyNames: depNames, Name: name, Value: val})
		} else {
			unusedArgs = append(unusedArgs, arg)
		}
	}

	return vars, unusedArgs, nil
}
