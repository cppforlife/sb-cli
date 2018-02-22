package cmd

import (
	"fmt"

	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
)

type ValueBoolOptional struct{ B *bool }

func (t ValueBoolOptional) String() string {
	if t.B == nil {
		return ""
	}
	return fmt.Sprintf("%t", *t.B)
}

func (t ValueBoolOptional) Value() boshtbl.Value            { return t }
func (t ValueBoolOptional) Compare(other boshtbl.Value) int { panic("unreachable") }
