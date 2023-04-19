package constset

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/max-gui/consulagent/pkg/consulsets"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/logagent/pkg/logsets"
	"github.com/max-gui/redisagent/pkg/redisops"
	"gopkg.in/yaml.v2"
)

const PthSep = string(os.PathSeparator)

var (
	Consolvername, RedisPwd, IacBranch, Gitname, Gitemail                                                                                                *string
	Repopathname, envset, Commitmsg, ConfArchPrefix, DeployInfoPrefix, ConfOrgPrefix, ConfResPrefix, ConfTeamProjPrefix                                  *string
	ConfWatchPrefix, ConfbalckPrefix, ConfwhitePrefix, ConfFabioPrefix, ConfmanBalckPrefix, ConfTraefikPrefix                                            *string
	EnvSet                                                                                                                                               []string
	Reppath, Archpath, Iacpath, Templepath, DbPath, Defconfpath, Archname, Iacname, Templname, Dbname, Sshkey, Archurl, IacUrl, Dburl, Codeurl, Templurl string
)

func StartupInit(bytes []byte, c context.Context) {

	EnvSet = strings.Split(*envset, ",")
	Reppath = *logsets.Apppath + PthSep + *Repopathname + PthSep
	Archpath = Reppath + "arch" + PthSep
	Iacpath = Reppath + "iac" + PthSep
	Templepath = Reppath + "temple" + PthSep
	DbPath = Reppath + "db" + PthSep
	Defconfpath = Templepath + "defaultconfig.yaml"
	Archname = "Archrepo"
	Iacname = "Iacrepo"
	Templname = "Templrepo"
	Dbname = "Dbrepo"
	readkey, err := ioutil.ReadFile(*logsets.Apppath + string(os.PathSeparator) + "code_key")
	log := logagent.InstPlatform(c)
	if err != nil {
		log.Panic(err)
	}
	Sshkey = string(readkey)

	seckey, err := ioutil.ReadFile(*logsets.Apppath + string(os.PathSeparator) + "sec.json")
	if err != nil {
		log.Panic(err)
	}
	secmap := make(map[string]string)
	json.Unmarshal(seckey, &secmap)
	IacUrl = secmap["iacurl"]
	Templurl = secmap["templurl"]
	Archurl = secmap["archurl"]
	Dburl = secmap["dburl"]
	Codeurl = secmap["codeurl"]
	// bytes, err := os.ReadFile(*Apppath + string(os.PathSeparator) + "application-" + *Appenv + ".yml")
	// if err != nil {
	// 	log.Panic(err)
	// }
	confmap := map[string]interface{}{}
	yaml.Unmarshal(bytes, confmap)
	*consulsets.Acltoken = confmap["af-arch"].(map[interface{}]interface{})["resource"].(map[interface{}]interface{})["private"].(map[interface{}]interface{})["acl-token"].(string)
	consulsets.StartupInit(*consulsets.Acltoken)
	redisopsUrl := confmap["af-arch"].(map[interface{}]interface{})["resource"].(map[interface{}]interface{})["redis-cluster-predixy"].(map[interface{}]interface{})["url"].(string)
	// redisops.Url = confmap["url"].(string)
	redisopsPwd := confmap["af-arch"].(map[interface{}]interface{})["resource"].(map[interface{}]interface{})["redis-cluster-predixy"].(map[interface{}]interface{})["password"].(string)
	// redisops.Pwd = confmap["password"].(string)
	redisops.StartupInit(redisopsUrl, redisopsPwd)
}

func init() {

	Consolvername = flag.String("consolvername", "127.0.0.1:8181", "confsolver's service name")
	IacBranch = flag.String("iacbranch", "iac", "iac git branch")
	envset = flag.String("envset", "sit,uat,prod", "envset spilt by ','")
	// Apppath = flag.String("apppath", "/Users/jimmy/Downloads/spells", "app root path")
	Repopathname = flag.String("repo", "repo", "repo path name")
	Gitname = flag.String("gitname", "", "git user name")
	Gitemail = flag.String("gitemail", "", "git user email")
	Commitmsg = flag.String("commitmsg", "STRUC#1742", "commit msg") //"FLS-AFLM-YW#6", "commit msg")
	ConfArchPrefix = flag.String("confArchPrefix", "ops/iac/arch/", "arch prefix for consul")
	DeployInfoPrefix = flag.String("srv_intention", "ops/iac/deploy_info/", "deploy info prefix for consul")
	ConfOrgPrefix = flag.String("confOrgPrefix", "ops/iac/org/", "arch prefix for consul")
	ConfResPrefix = flag.String("ConfResPrefix", "ops/resource/", "resource prefix for consul")
	ConfTeamProjPrefix = flag.String("ConfTeamProjPrefix", "ops/iac/team-proj", "teamproj prefix for consul")
	ConfWatchPrefix = flag.String("ConfWatchPrefix", "ops/", "watch prefix for consul")
	ConfbalckPrefix = flag.String("ConfbalckPrefix", "ops/iac/blacklist/", "black prefix for consul")
	ConfwhitePrefix = flag.String("ConfwhitePrefix", "ops/iac/whitelist/", "white prefix for consul")
	ConfFabioPrefix = flag.String("ConfFabioPrefix", "ops/iac/fabio", "fabio prefix for consul")
	ConfmanBalckPrefix = flag.String("ConfmanBalckPrefix", "ops/blacklist", "manual black list prefix for consul")
	ConfTraefikPrefix = flag.String("ConfTraefikPrefix", "ops/traefik", "traefik prefix for consul")

}

// var Reppath = func() string {
// 	return *Apppath + PthSep + *Repopathname + PthSep
// }
