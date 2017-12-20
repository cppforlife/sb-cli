package x

import (
	"strings"
)

type OpsFile struct {
	DependencyNames []string
	Path            string
}

func NewOpsFilesFromArgs(args []string) ([]OpsFile, []string, error) {
	var opsFiles []OpsFile
	var unusedArgs []string

	nextIsFile := false

	for _, arg := range args {
		if arg == "-o" {
			nextIsFile = true
		} else if nextIsFile {
			nextIsFile = false

			depNames := strings.Split(arg, ":")
			arg = depNames[len(depNames)-1]
			depNames = depNames[:len(depNames)-1]

			opsFiles = append(opsFiles, OpsFile{DependencyNames: depNames, Path: arg})
		} else {
			unusedArgs = append(unusedArgs, arg)
		}
	}

	return opsFiles, unusedArgs, nil
}
