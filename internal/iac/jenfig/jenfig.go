package jenfig

import (
	"context"
	"strings"
	"text/template"

	"github.com/max-gui/spells/internal/iac/archfig"
	"github.com/max-gui/spells/internal/iac/templ"
	"github.com/max-gui/spells/internal/pkg/constset"
)

type JenkinsInfo struct {
	Cmd        string
	CmdArgs    string
	Pom        string
	Classes    string
	OutputPath string
	Output     string
	AppId      string
	HelmDir    string
	AppName    string
	Jenkignor  []string
	Type       string
	Repositry  string
}

func GenJenfig(app_conf archfig.Arch_config) JenkinsInfo {

	temparg := JenkinsInfo{
		Cmd:     app_conf.Deploy.Build.Cmd,
		CmdArgs: app_conf.Deploy.Build.Args,
		Pom:     app_conf.Deploy.Build.Pkgconf,
		// Classes:    app_conf.Deploy.Build.OutputPath + "/classes",
		// OutputPath: app_conf.Deploy.Build.OutputPath,
		Output:    app_conf.Deploy.Build.Output,
		AppId:     app_conf.Application.Appid,
		HelmDir:   "/data/helm/af-hercules/" + strings.TrimPrefix(app_conf.Application.Appid, "fls-"),
		AppName:   app_conf.Application.Name,
		Jenkignor: app_conf.Deploy.Build.Jenkignor,
		Type:      app_conf.Application.Type,
		Repositry: app_conf.Application.Repositry,
	}

	return temparg
}

func MakeJenkimple(apptype string, isinstall bool, c context.Context) *template.Template {
	name := "jenkins." + apptype
	// templtmp := templ.GetemplFrom(Templepath()+name, name)
	dirPth := constset.Templepath + name
	//docker templ
	res := templ.GetemplFrom(dirPth, name, isinstall, c)

	return res

}

func GenJenfile(jenfig JenkinsInfo, c context.Context) string {
	// dirPth := orgconfigPth + pthSep + "arch.yaml"
	// str, _ := ioutil.ReadFile(dirPth)
	// a := Arch_config{}
	// yaml.Unmarshal(str, &a)
	// t.Log(a)
	name := "jenkins." + jenfig.Type
	result := templ.GemplFrom(name, jenfig, c)

	return result
}
