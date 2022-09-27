package dockfig

import (
	"context"
	"text/template"

	"github.com/max-gui/spells/internal/iac/archfig"
	"github.com/max-gui/spells/internal/iac/templ"
	"github.com/max-gui/spells/internal/iac/valfig"
	"github.com/max-gui/spells/internal/pkg/constset"
)

func MakeDockemple(apptype string, isinstall bool, c context.Context) *template.Template {
	name := "Dockerfile." + apptype
	dirPth := constset.Templepath + name
	//docker templ
	res := templ.GetemplFrom(dirPth, name, isinstall, c)

	return res

}

func GenDocfile(appconf archfig.Arch_config, c context.Context) string {
	// dirPth := orgconfigPth + pthSep + "arch.yaml"
	// str, _ := ioutil.ReadFile(dirPth)
	// a := Arch_config{}
	// yaml.Unmarshal(str, &a)
	// t.Log(a)
	// name := "Dockerfile." + appconf.Application.Type
	// result := templ.GemplFrom(name, appconf, c)

	name := "Dockerfile"
	result := templ.GemplFromType(name, appconf.Application.Type, appconf, c)

	return result
}

func GenRuntimeDocfile(appconf archfig.Arch_config, valconf valfig.ValuesInfo, c context.Context) string {
	// dirPth := orgconfigPth + pthSep + "arch.yaml"
	// str, _ := ioutil.ReadFile(dirPth)
	// a := Arch_config{}
	// yaml.Unmarshal(str, &a)
	// t.Log(a)
	dockerfile := GenDocfile(appconf, c)
	res := templ.Getempl(dockerfile, "", false, c) // GetemplFrom(dirPth, templname, true)
	result := templ.Gempl(res, valconf, c)
	// name := "Dockerfile." + appconf.Application.Type
	// result := templ.GemplFrom(name, appconf)

	return result
}
