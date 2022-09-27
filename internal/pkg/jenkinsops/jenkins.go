package jenkinsops

import (
	"context"
	"encoding/xml"
	"strings"

	"github.com/bndr/gojenkins"
	"github.com/max-gui/logagent/pkg/confload"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/spells/internal/pkg/constset"
)

type jobInfo struct {
	Definition struct {
		Scm struct {
			UserRemoteConfigs struct {
				Hudson_plugins_git_UserRemoteConfig struct {
					Url string `xml:"url"`
				} `xml:"hudson.plugins.git.UserRemoteConfig"`
			} `xml:"userRemoteConfigs"`
			Branches struct {
				Hudson_plugins_git_BranchSpec struct {
					Name string `xml:"name"`
				} `xml:"hudson.plugins.git.BranchSpec"`
			} `xml:"branches"`
		} `xml:"scm"`
		ScriptPath string `xml:"scriptPath"`
	} `xml:"definition"`
}

func GetJenkins(envdc string, c context.Context) (*gojenkins.Jenkins, string, error) {

	log := logagent.InstArch(c)

	// confmap := map[string]interface{}{}
	confmap := confload.LoadEnv(envdc, c).(map[string]interface{})
	log.Print(confmap)
	jenkinsurl := confmap["af-arch"].(map[interface{}]interface{})["resource"].(map[interface{}]interface{})["jenkins"].(map[interface{}]interface{})["url"].(string)
	user := confmap["af-arch"].(map[interface{}]interface{})["resource"].(map[interface{}]interface{})["jenkins"].(map[interface{}]interface{})["user"].(string)
	pwd := confmap["af-arch"].(map[interface{}]interface{})["resource"].(map[interface{}]interface{})["jenkins"].(map[interface{}]interface{})["pwd"].(string)
	deskurl := confmap["af-arch"].(map[interface{}]interface{})["resource"].(map[interface{}]interface{})["jenkins"].(map[interface{}]interface{})["deskurl"].(string)
	// confmap := consulhelp.Getconfaml(*constset.ConfResPrefix, "deploy", "jenkins", envdc, c)
	jenkins := gojenkins.CreateJenkins(nil, jenkinsurl, user, pwd)

	_, err := jenkins.Init(c)
	if err != nil {
		log.Panic(err)
	}
	return jenkins, deskurl, err
}

func GetJobXML(appDescription, appname, appid string) string {
	scriptPath, giturl, gitbranch := genJobInfo(appid, appname)

	jenkinsJobxml := "<flow-definition plugin=\"workflow-job@2.32\">\n" +
		"<description>" + appDescription + "</description>\n" +
		"<keepDependencies>false</keepDependencies>\n" +
		"<definition class=\"org.jenkinsci.plugins.workflow.cps.CpsScmFlowDefinition\" plugin=\"workflow-cps@2.65\">\n" +
		"<scm class=\"hudson.plugins.git.GitSCM\" plugin=\"git@3.9.3\">\n" +
		"<configVersion>2</configVersion>\n" +
		"<userRemoteConfigs>\n" +
		"<hudson.plugins.git.UserRemoteConfig>\n" +
		"<url>" + giturl + "</url>\n" +
		"</hudson.plugins.git.UserRemoteConfig>\n" +
		"</userRemoteConfigs>\n" +
		"<branches>\n" +
		"<hudson.plugins.git.BranchSpec>\n" +
		"<name>" + gitbranch + "</name>\n" +
		"</hudson.plugins.git.BranchSpec>\n" +
		"</branches>\n" +
		"<doGenerateSubmoduleConfigurations>false</doGenerateSubmoduleConfigurations>\n" +
		"<submoduleCfg class=\"list\"/>\n" +
		"<extensions/>\n" +
		"</scm>\n" +
		"<scriptPath>" + scriptPath + "</scriptPath>\n" +
		"<lightweight>true</lightweight>\n" +
		"</definition>\n" +
		"<triggers/>\n" +
		"<disabled>false</disabled>\n" +
		"</flow-definition>"
	return jenkinsJobxml
}

func genJobInfo(appid, appname string) (string, string, string) {
	scriptPath := "app/iac/" + appid + "/" + appname + "/Jenkinsfile"
	giturl := constset.Codeurl + "/af-archops/jenkins-library.git"
	gitbranch := "*/" + *constset.IacBranch

	return scriptPath, giturl, gitbranch
}

func IsSameJob(jenkins *gojenkins.Jenkins, appname, appid string, c context.Context) bool {

	log := logagent.InstArch(c)
	bb := ""
	para := "/job/iac-" + appname + "/config.xml"
	_, err := jenkins.Requester.GetXML(c, para, &bb, nil)
	if err != nil {
		log.Panic(err)
	}
	bb = strings.Replace(bb, "version='1.1'", "version='1.0'", 1)
	// log.Print(bb)

	v := jobInfo{}
	err = xml.Unmarshal([]byte(bb), &v)
	if err != nil {
		log.Panic(err)
	}
	// _, o, _ := strings.Cut(bb, "<scriptPath>app/iac/")
	// // log.Print(o)
	// p, _, _ := strings.Cut(o, "/Jenkinsfile</scriptPath>")
	// log.Print(p)
	// log.Print(m)
	scriptPath, giturl, gitbranch := genJobInfo(appid, appname)
	return scriptPath == v.Definition.ScriptPath && giturl == v.Definition.Scm.UserRemoteConfigs.Hudson_plugins_git_UserRemoteConfig.Url && gitbranch == v.Definition.Scm.Branches.Hudson_plugins_git_BranchSpec.Name
}
