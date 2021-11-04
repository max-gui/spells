package jenkinsops

import (
	"context"

	"github.com/bndr/gojenkins"
	"github.com/max-gui/consulagent/pkg/consulhelp"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/spells/internal/pkg/constset"
)

func GetJenkins(env string, c context.Context) (*gojenkins.Jenkins, string, error) {

	log := logagent.Inst(c)

	confmap := consulhelp.Getconfaml(*constset.ConfResPrefix, "deploy", "jenkins", env, c)
	jenkins := gojenkins.CreateJenkins(nil, confmap["url"].(string), confmap["user"].(string), confmap["pwd"].(string))

	_, err := jenkins.Init(c)
	if err != nil {
		log.Panic(err)
	}
	return jenkins, confmap["deskurl"].(string), err
}

func GetJobXML(appDescription, appname, appid string) string {
	jenkinsJobxml := "<flow-definition plugin=\"workflow-job@2.32\">\n" +
		"<description>" + appDescription + "</description>\n" +
		"<keepDependencies>false</keepDependencies>\n" +
		"<definition class=\"org.jenkinsci.plugins.workflow.cps.CpsScmFlowDefinition\" plugin=\"workflow-cps@2.65\">\n" +
		"<scm class=\"hudson.plugins.git.GitSCM\" plugin=\"git@3.9.3\">\n" +
		"<configVersion>2</configVersion>\n" +
		"<userRemoteConfigs>\n" +
		"<hudson.plugins.git.UserRemoteConfig>\n" +
		"<url>" + constset.Codeurl + "/af-archops/jenkins-library.git" + "</url>\n" +
		"</hudson.plugins.git.UserRemoteConfig>\n" +
		"</userRemoteConfigs>\n" +
		"<branches>\n" +
		"<hudson.plugins.git.BranchSpec>\n" +
		"<name>*/" + *constset.IacBranch + "</name>\n" +
		"</hudson.plugins.git.BranchSpec>\n" +
		"</branches>\n" +
		"<doGenerateSubmoduleConfigurations>false</doGenerateSubmoduleConfigurations>\n" +
		"<submoduleCfg class=\"list\"/>\n" +
		"<extensions/>\n" +
		"</scm>\n" +
		"<scriptPath>" + "app/iac/" + appid + "/" + appname + "/Jenkinsfile" + "</scriptPath>\n" +
		"<lightweight>true</lightweight>\n" +
		"</definition>\n" +
		"<triggers/>\n" +
		"<disabled>false</disabled>\n" +
		"</flow-definition>"
	return jenkinsJobxml
}
