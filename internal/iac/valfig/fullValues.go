package valfig

import (
	"context"

	"github.com/max-gui/spells/internal/iac/archfig"
	"github.com/max-gui/spells/internal/iac/templ"
)

type FullValues struct {
	Appconf    archfig.Arch_config
	Values     map[string]interface{}
	Neighbours []string
}

func GenFullValues(infos map[string]interface{}, c context.Context) string {
	// dirPth := orgconfigPth + pthSep + "arch.yaml"
	// str, _ := ioutil.ReadFile(dirPth)
	// a := Arch_config{}
	// yaml.Unmarshal(str, &a)
	// t.Log(a)

	result := templ.GemplFromType("full.yml", "values", infos, c)

	return result
}
