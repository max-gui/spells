package archfig

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/max-gui/consulagent/pkg/consulhelp"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/spells/internal/iac/defig"
	"github.com/max-gui/spells/internal/pkg/constset"
	"gopkg.in/yaml.v2"
)

type FileContentInfo struct {
	Path    string
	Status  string //a,d,m
	Content string
}

type EnvInfo struct {
	Env string
	Dc  string
}
type BlackInfo struct {
	Visible   bool     `yaml:"visible,omitempty"`
	Blacklist []string `yaml:"blacklist,omitempty"`
}

type ExposeInfo struct {
	Internet   BlackInfo `yaml:"internet,omitempty"`
	Intranet   BlackInfo `yaml:"intranet,omitempty"`
	Clusternet WhiteInfo `yaml:"clusternet,omitempty"`
	Ptrnet     WhiteInfo `yaml:"ptrnet,omitempty"`
	// Black      struct {
	// 	Internet bool `yaml:"internet,omitempty"`
	// 	Intrenet bool `yaml:"intrenet,omitempty"`
	// } `yaml:"black,omitempty"`
	Unsafe bool `yaml:"unsafe,omitempty"`
	// Service string`yaml:"omitempty"`
	// Secgwon bool`yaml:"omitempty"`
	Host       string `yaml:"host,omitempty"`
	Expovice   string `yaml:"expovice,omitempty"`
	PrefixPath string `yaml:"prefixPath,omitempty"`
	SurfixPath string `yaml:"surfixPath,omitempty"`
	Manual     []struct {
		Host       string `yaml:"host,omitempty"`
		PrefixPath string `yaml:"prefixPath,omitempty"`
		SurfixPath string `yaml:"surfixPath,omitempty"`
	} `yaml:"menual,omitempty"`
}

// type ExposeInfoAll struct {
// 	Internet   BlackInfo `yaml:"internet,omitempty"`
// 	Intranet   BlackInfo `yaml:"intranet,omitempty"`
// 	Clusternet WhiteInfo `yaml:"clusternet,omitempty"`
// 	Ptrnet     WhiteInfo `yaml:"ptrnet,omitempty"`
// 	// Black      struct {
// 	// 	Internet bool `yaml:"internet,omitempty"`
// 	// 	Intrenet bool `yaml:"intrenet,omitempty"`
// 	// } `yaml:"black,omitempty"`
// 	Unsafe bool `yaml:"unsafe,omitempty"`
// 	// Service string`yaml:"omitempty"`
// 	// Secgwon bool`yaml:"omitempty"`
// 	Expovice   string `yaml:"expovice,omitempty"`
// 	PrefixPath string `yaml:"prefixPath,omitempty"`
// 	SurfixPath string `yaml:"surfixPath,omitempty"`
// 	Manual     []struct {
// 		PrefixPath string `yaml:"prefixPath,omitempty"`
// 		SurfixPath string `yaml:"surfixPath,omitempty"`
// 	} `yaml:"menual,omitempty"`
// 	Appname string `yaml:"appname,omitempty"`
// }

type WhiteInfo struct {
	Open bool `yaml:"open,omitempty"`
}

type Arch_config struct {
	Application struct {
		Name        string            `yaml:"name,omitempty"`
		Description string            `yaml:"description,omitempty"`
		Appid       string            `yaml:"appid,omitempty"`
		Type        string            `yaml:"type,omitempty"`
		Language    string            `yaml:"language,omitempty"`
		Langval     string            `yaml:"langval,omitempty"`
		NoSource    bool              `yaml:"nosource,omitempty"`
		Repositry   string            `yaml:"repositry,omitempty"`
		Resource    map[string]string `yaml:"resource,omitempty"`
		Service     []string          `yaml:"service,omitempty"`
		Extservice  []string          `yaml:"extservice,omitempty"`
		Team        string            `yaml:"team,omitempty"`
		Project     string            `yaml:"project,omitempty"`
		Ungenfig    bool              `yaml:"ungenfig,omitempty"`
		Allpath     bool              `yaml:"allpath,omitempty"`
	}
	Environment struct {
		NodeSelector map[string]string `yaml:"nodeSelector,omitempty"`
		// K8snode      string            `yaml:"k8snode,omitempty"`
		Strategy map[string]struct {
			Capacity string `yaml:"capacity,omitempty"`
			Cpu      string `yaml:"cpu,omitempty"`
			Mem      string `yaml:"mem,omitempty"`
			Replica  int    `yaml:"replica,omitempty"`
		} `yaml:"strategy,omitempty"`
		Resource map[string][]string `yaml:"resource,omitempty"`
		// Volumns []string
		Tag            map[string]string `yaml:"tag,omitempty"`
		Port           string            `yaml:"port,omitempty"`
		EnHostportable bool              `yaml:"enhostportable,omitempty"`
		Hostport       string            `yaml:"hostport,omitempty"`
		IsHostNetwork  bool              `yaml:"ishostnetwork,omitempty"`
		// Hosts   []string
		Expose ExposeInfo `yaml:"expose,omitempty"`
	} `yaml:"environment,omitempty"`
	Deploy struct {
		Limited  map[string]struct{} `yaml:"limited,omitempty"`
		Strategy []struct {
			Flow string   `yaml:"flow,omitempty"`
			Env  []string `yaml:"env,omitempty"`
		} `yaml:"strategy,omitempty"`
		Stratail map[string][][]EnvInfo `yaml:"stratail,omitempty"`
		Runtime  struct {
			Args []string            `yaml:"args,omitempty"`
			Ign  map[string][]string `yaml:"ign,omitempty"`
		} `yaml:"runtime,omitempty"`
		Build struct {
			// Appyml  string `yaml:"appyml,omitempty"`
			Cmd     string `yaml:"cmd,omitempty"`
			Args    string `yaml:"args,omitempty"`
			Pkgconf string `yaml:"pkgconf,omitempty"`
			Output  string `yaml:"output,omitempty"`
			// OutputPath string`yaml:"omitempty"`
			// PrePackage string`yaml:"omitempty"`
			Jenkignor []string `yaml:"jenkignor,omitempty"`
			Jenkexec  []string `yaml:"jenkexec,omitempty"`
		}
		Sidecar struct {
			Neighbour []string            `yaml:"neighbour,omitempty"`
			Ign       map[string][]string `yaml:"ign,omitempty"`
		} `yaml:"sidecar,omitempty"`
	}
	Docker struct {
		From   string   `yaml:"from,omitempty"`
		Append []string `yaml:"append,omitempty"`
		Cmd    string   `yaml:"cmd,omitempty"`
	} `yaml:"docker,omitempty"`
}

func getContentInfo(fcinfo FileContentInfo, c context.Context) []string {
	log := logagent.InstArch(c)

	log.Print(fcinfo.Content)
	tailpath := strings.Split(fcinfo.Path, "arch/")[1]
	if strings.Contains(tailpath, "/") {
		log.Panicf("arch file should be in arch folder,get %s", fcinfo.Path)
	}
	projinfo := strings.Split(fcinfo.Path, constset.PthSep)
	return projinfo
}

func GetArchfigByGitContentSin(app_conf Arch_config, appname string, install bool, c context.Context) Arch_config {

	appconf := GenArchConfigSinFrominst(app_conf, appname, install, c) // GenArchConfigSin([]byte(content), appname, install, c)
	return appconf
}

func GetArchfigByGitContent(fcinfo FileContentInfo, install bool, c context.Context) Arch_config {
	projinfo := getContentInfo(fcinfo, c)

	appconf := GenArchConfig([]byte(fcinfo.Content), projinfo[0], projinfo[1], strings.TrimSuffix(projinfo[2], ".yaml"), install, c)
	return appconf
}

func GetArchfig(fcinfo FileContentInfo, c context.Context) Arch_config {
	projinfo := getContentInfo(fcinfo, c)

	appconf := GetAppconfig(strings.TrimSuffix(projinfo[2], ".yaml"), projinfo[0], projinfo[1], c)
	return appconf
}

func GetArchfigSin(appname string, c context.Context) Arch_config {

	appconf := Arch_config{}
	defer func() {
		if e := recover(); e != nil {

			logger := logagent.InstArch(c)
			logger.WithField("misarch", appname).
				Info("miss iac data")
		}
	}()
	appbytes := consulhelp.GetConfigFull(*constset.ConfArchPrefix+appname, c)
	json.Unmarshal(appbytes, &appconf)
	return appconf
	// if appbytes != nil {
	// 	json.Unmarshal(appbytes, &appconf)
	// 	return appconf
	// } else {
	// 	// log.Panic("app arch does not exist")
	// 	// return GenArchConfigFromSin(appname+".yaml", appname, true)
	// }
}

func GetAppconfigOnline(appname string, c context.Context) Arch_config {
	appconf := Arch_config{}
	appbytes := consulhelp.GetConfigFull(*constset.ConfArchPrefix+appname, c)
	json.Unmarshal(appbytes, &appconf)
	return appconf
}
func GetAppconfig(appname, team, proj string, c context.Context) Arch_config {
	appconf := Arch_config{}
	appbytes := consulhelp.GetConfigFull(*constset.ConfArchPrefix+appname, c)
	if appbytes != nil {
		json.Unmarshal(appbytes, &appconf)
		return appconf
	} else {
		return GenArchConfigFrom(team+constset.PthSep+proj+constset.PthSep+appname+".yaml", team, proj, appname, true, c)
	}
	// if v, ok := app_arch_map[appname]; ok {
	// 	return v
	// } else {
	// 	rediscli := redisops.Pool().Get()

	// 	defer rediscli.Close()
	// 	appconf := Arch_config{}
	// 	mm, err := redis.Bytes(rediscli.Do("HGET", "arch-spell-appconfig", appname))
	// 	json.Unmarshal(mm, &appconf)

	// 	if err != nil {
	// 		return GenArchConfigFrom(team+constset.PthSep+proj+constset.PthSep+appname+".yaml", team, proj, appname, true)
	// 	} else {
	// 		app_arch_map[appname] = appconf
	// 		return appconf
	// 	}

	// 	//	rediscli.Do("HSET", "arch-spell-appconfig", app_conf.Application.Name, jsonbs)
	// }
}

func GenArchConfigFromSin(appconfpath, appfname string, isinstall bool, c context.Context) Arch_config {

	log := logagent.InstArch(c)

	dirPth := constset.Archpath + appconfpath
	str, err := os.ReadFile(dirPth)
	if err != nil {
		log.Panic(err)
	}
	return GenArchConfigSin(str, appfname, isinstall, c)
}

func GenArchConfigFrom(appconfpath, team, proj, appfname string, isinstall bool, c context.Context) Arch_config {

	log := logagent.InstArch(c)

	dirPth := constset.Archpath + appconfpath
	str, err := os.ReadFile(dirPth)
	if err != nil {
		log.Panic(err)
	}
	return GenArchConfig(str, team, proj, appfname, isinstall, c)
}

func GetAppconfigFull(appname string, regenflag bool, c context.Context) Arch_config {
	if regenflag {
		return ReGenArchConfig(appname, c)
	} else {

		return GetAppconfigOnline(appname, c)
	}
}

func ReGenArchConfig(appfname string, c context.Context) Arch_config {
	orgconf := Arch_config{}
	appbytes := consulhelp.GetConfigFull(*constset.ConfOrgPrefix+appfname, c)
	json.Unmarshal(appbytes, &orgconf)
	// return appconf
	// confstr := GetArchfigSin(appfname)
	return GenArchConfigSinFrominst(orgconf, appfname, true, c)
}

func GenArchConfig(appconf []byte, team, proj, appfname string, isinstall bool, c context.Context) Arch_config {

	log := logagent.InstArch(c)

	app_conf := GenArchConfigSin(appconf, appfname, isinstall, c)

	app_conf.Application.Team = team
	app_conf.Application.Project = proj
	// app_conf.Deploy.Runtime.Args += "-Dspring.application.name=" + app_conf.Application.Name
	if !app_conf.Application.Ungenfig && (app_conf.Application.Team == "" || app_conf.Application.Project == "") {

		log.Panicf("team or project should not be empty, get:\nteam:%s project:%s", app_conf.Application.Team, app_conf.Application.Project)
	}

	if !(strings.HasPrefix(app_conf.Application.Repositry, "ssh") || strings.HasPrefix(app_conf.Application.Repositry, "git")) {
		log.Panicf("repo shoud use ssh or git, get:%s", app_conf.Application.Repositry)
	}
	return app_conf
}

func GenArchConfigSin(appconf []byte, appfname string, isinstall bool, c context.Context) Arch_config {
	// log := logagent.Inst(c)

	orgconf := GenArchConfFromBytes(appconf, c)
	app_conf := Arch_config{}
	yaml.Unmarshal(appconf, &app_conf)

	archconf := GenArchConfigSinFrominst(app_conf, appfname, isinstall, c)

	// bytes, _ := json.MarshalIndent(orgconf, "", "    ")

	// consulhelp.PutConfig(*constset.ConfArchPrefix, v.Application.Team, v.Application.Project, v.Application.Name, bytes)
	// log.Print(string(bytes))
	if isinstall {
		orgconf.OrgInstall(c)
		// consulhelp.PutConfigFull(*constset.ConfOrgPrefix+orgconf.Application.Name, bytes, c)
	}
	return archconf
}

func GenArchConfFromBytes(appconf []byte, c context.Context) Arch_config {
	log := logagent.InstArch(c)

	orgconf := Arch_config{}
	log.Print(string(appconf))
	err := yaml.Unmarshal(appconf, &orgconf)
	if err != nil {
		log.Panic(err)
	}
	return orgconf
}

func GenArchConfigSinFrominst(app_conf Arch_config, appfname string, isinstall bool, c context.Context) Arch_config {
	// dirPth := appconfpath
	// str, _ := ioutil.ReadFile(dirPth)

	log := logagent.InstArch(c)

	defconf := defig.GetDefconfig(c)
	log.Print(app_conf)

	// if !app_conf.Application.Ungenfig && (app_conf.Application.Language == "") {
	// 	log.Panicf("Language should not be empty, get:\n%s", app_conf.Application.Language)
	// }
	// if !app_conf.Application.Ungenfig && (app_conf.Application.Description == "") {
	// 	log.Panicf("Description should not be empty, get:\n%s", app_conf.Application.Description)
	// }
	if !app_conf.Application.Ungenfig && (app_conf.Application.Team == "" || app_conf.Application.Project == "") {

		log.Panicf("team or project should not be empty, get:\nteam:%s project:%s", app_conf.Application.Team, app_conf.Application.Project)
	}
	if app_conf.Application.Name != appfname {
		log.Panicf("application name shoud be equal with filename, get:\napplication name:%s appfilename:%s", app_conf.Application.Name, appfname)
	}
	if len(app_conf.Application.Name) > 35 {
		log.Panicf("application name shoud be less than 35, get:%s", app_conf.Application.Name)
	}
	bytes := consulhelp.GetConfigFull(*constset.ConfTeamProjPrefix, c)
	teamproj := map[string]map[string]string{}
	yaml.Unmarshal(bytes, &teamproj)
	if app_conf.Application.Language != "" {
		if lang, ok := teamproj["language"][app_conf.Application.Language]; !ok {
			log.Panic("this language is not existed")
		} else {
			app_conf.Application.Langval = lang
		}
	}
	if _, ok := teamproj["team"][app_conf.Application.Team]; !ok {
		log.Panic("this team is not existed")

	}
	if _, ok := teamproj["proj"][app_conf.Application.Project]; !ok {

		log.Panic("this proj is not existed")
	}

	if app_conf.Application.Appid == "" || app_conf.Application.Appid != strings.ReplaceAll(strings.ToLower(app_conf.Application.Appid), " ", "") {
		log.Panicf("application Appid shoud not be empty,captitalism or with space, get:\napplication Appid:%s", app_conf.Application.Appid)
	}
	if app_conf.Application.Type == "" || app_conf.Application.Type != strings.ReplaceAll(strings.ToLower(app_conf.Application.Type), " ", "") {
		log.Panicf("application Type shoud not be empty,captitalism or with space, get:\n%s", app_conf.Application.Type)
	}
	// app_conf.Application.Appid = strings.ToLower(app_conf.Application.Appid)
	// app_conf.Application.Name = strings.ToLower(app_conf.Application.Name)

	if app_conf.Environment.Expose.PrefixPath == "" && (app_conf.Environment.Expose.Unsafe || app_conf.Environment.Expose.Expovice != "") {
		log.Panicf("Expose.Path can't be empty if you want to expose a service, get:%s", app_conf.Environment.Expose.PrefixPath)
	}

	for _, v := range app_conf.Environment.Expose.Internet.Blacklist {
		if !strings.HasPrefix(v, app_conf.Environment.Expose.PrefixPath) {
			log.Panicf("Internet.Blacklist must start with Expose.PrefixPath, get black:%s;get prefix:%s", v, app_conf.Environment.Expose.PrefixPath)
		}
	}

	for _, v := range app_conf.Environment.Expose.Intranet.Blacklist {
		if !strings.HasPrefix(v, app_conf.Environment.Expose.PrefixPath) {
			log.Panicf("Intranet.Blacklist must start with Expose.PrefixPath, get black:%s;get prefix:%s", v, app_conf.Environment.Expose.PrefixPath)
		}
	}

	if strings.HasPrefix(app_conf.Environment.Expose.PrefixPath, "/") {
		log.Panic("app_conf.Environment.Expose.Path should not start with /")
	}

	for _, v := range app_conf.Application.Service {
		if strings.ToLower(v) != v {
			log.Panicf("app' service must be lower-case,get:%s", v)
		}
	}
	//get default capacity temple
	capkeys := defconf.GetCapacity()
	//gen default capacity temple
	for k, v := range defconf.Defualtinfo.Capacity {
		if v.Capacity != "" {
			v.Cpu = defconf.Capacitylable[v.Capacity].Cpu
			v.Mem = defconf.Capacitylable[v.Capacity].Mem
			defconf.Defualtinfo.Capacity[k] = v
		}
	}
	//gen env range from default capacity temple
	envkeys := defconf.GetEnvKeys()

	if app_conf.Environment.Strategy == nil {
		app_conf.Environment.Strategy = make(map[string]struct {
			Capacity string `yaml:"capacity,omitempty"`
			Cpu      string `yaml:"cpu,omitempty"`
			Mem      string `yaml:"mem,omitempty"`
			Replica  int    `yaml:"replica,omitempty"`
		})
	}
	app_conf.Deploy.Limited = map[string]struct{}{}
	for _, v := range defconf.Defualtinfo.Deploy.Limited {
		app_conf.Deploy.Limited[v] = struct{}{}
	}

	//gen app deploy env capacity
	for env, defconv := range defconf.Defualtinfo.Capacity {

		//app with self setting or use default setting
		if appconf, ok := app_conf.Environment.Strategy[env]; ok {
			//app with cpu mem setting
			//app with capacity setting
			if appconf.Capacity == "" {
				tmpfloat, err := strconv.ParseFloat(appconf.Cpu, 32)
				if err != nil || tmpfloat <= 0 {
					log.Panicf("%s cpu config is wrong, get:%s", env, appconf.Cpu)
				}
				if !strings.HasSuffix(appconf.Mem, "Mi") {
					log.Panicf("%s mem config is wrong, should with suffix 'Mi' get:%s", env, appconf.Mem)
				}
				tmp, err := strconv.Atoi(strings.TrimSuffix(appconf.Mem, "Mi"))
				if err != nil || tmp <= 0 {
					log.Panicf("%s mem config is wrong, get:%s", env, appconf.Mem)
				}
			} else if v, ok := defconf.Capacitylable[appconf.Capacity]; ok {
				appconf.Cpu = v.Cpu
				appconf.Mem = v.Mem
			} else {

				log.Panicf("capacity should be in range of %v, get:%s", capkeys, appconf.Capacity)
			}

			if appconf.Replica <= 0 {
				appconf.Replica = defconv.Replica
			}

			app_conf.Environment.Strategy[env] = appconf
		} else {
			app_conf.Environment.Strategy[env] = struct {
				Capacity string `yaml:"capacity,omitempty"`
				Cpu      string `yaml:"cpu,omitempty"`
				Mem      string `yaml:"mem,omitempty"`
				Replica  int    `yaml:"replica,omitempty"`
			}{Cpu: defconv.Cpu, Mem: defconv.Mem, Replica: defconv.Replica}
		}
	}

	for k := range app_conf.Environment.Strategy {
		if _, ok := defconf.Defualtinfo.Capacity[k]; !ok {
			log.Panicf("env should be in range of %v, get:%s", envkeys, k)
		}
	}

	if app_conf.Environment.Resource == nil {
		app_conf.Environment.Resource = make(map[string][]string)
	}
	if app_conf.Environment.Tag == nil {
		app_conf.Environment.Tag = make(map[string]string)
	}

	for k, v := range app_conf.Environment.Tag {
		if strings.Contains(k, ".") {
			log.Panicf("tag's key should not with '.' of %s, get:%s", k, v)
		}
	}

	if app_conf.Environment.Port == "" {
		app_conf.Environment.Port = defconf.Defualtinfo.Port
	}

	if app_conf.Environment.EnHostportable {
		app_conf.Environment.Hostport = app_conf.Environment.Port
	}

	//add default resources to app resoureces slice
	for k, v := range defconf.Defualtinfo.Resource {
		app_conf.Environment.Resource[k] = append(app_conf.Environment.Resource[k], v...)
	}
	//remove duplicate resources
	for k, v := range app_conf.Environment.Resource {
		result := make([]string, 0, len(v))
		temp := map[string]struct{}{}
		rmkeys := map[string]struct{}{}
		// rmlisst := []string{}
		for _, item := range v {
			if strings.HasPrefix(item, "-") {
				itemrm := strings.TrimPrefix(item, "-")
				// rmlisst = append(rmlisst, itemrm)
				delete(temp, itemrm)
				rmkeys[itemrm] = struct{}{}
			} else if _, ok := temp[item]; !ok {
				if _, ok := rmkeys[item]; !ok {
					temp[item] = struct{}{}
				}
				// result = append(result, item)
			}
		}
		for k := range temp {
			result = append(result, k)
		}

		app_conf.Environment.Resource[k] = result
	}

	app_conf.Deploy.Sidecar.Neighbour = append(app_conf.Deploy.Sidecar.Neighbour, defconf.Defualtinfo.Sidecar.Neighbour[app_conf.Application.Type]...)
	app_conf.Deploy.Sidecar.Neighbour = removeDuplicateElement(app_conf.Deploy.Sidecar.Neighbour)
	// app_conf.Environment.Volumns = append(app_conf.Environment.Volumns, defconf.Defualtinfo.Volumes...)

	if app_conf.Deploy.Sidecar.Ign == nil {
		app_conf.Deploy.Sidecar.Ign = make(map[string][]string)
	}
	// merge app runtime Ign with default runtime Ign;k is env
	for k := range defconf.Defualtinfo.Capacity {
		app_conf.Deploy.Sidecar.Ign[k] = append(app_conf.Deploy.Sidecar.Ign[k], defconf.Defualtinfo.Sidecar.Ign[k]...)
		app_conf.Deploy.Sidecar.Ign[k] = removeDuplicateElement(app_conf.Deploy.Sidecar.Ign[k])
	}

	// merge app runtime args with default certain type runtime args
	app_conf.Deploy.Runtime.Args = append(app_conf.Deploy.Runtime.Args, defconf.Defualtinfo.Cmdarg.Args[app_conf.Application.Type]...)
	app_conf.Deploy.Runtime.Args = removeDuplicateElement(app_conf.Deploy.Runtime.Args)
	// app_conf.Deploy.Runtime.Args = strings.ReplaceAll(app_conf.Deploy.Runtime.Args, defconf.Defualtinfo.Cmdarg.Args[app_conf.Application.Type], "")
	// app_conf.Deploy.Runtime.Args = app_conf.Deploy.Runtime.Args + " " + defconf.Defualtinfo.Cmdarg.Args[app_conf.Application.Type]

	if app_conf.Deploy.Runtime.Ign == nil {
		app_conf.Deploy.Runtime.Ign = make(map[string][]string)
	}
	// merge app runtime Ign with default runtime Ign;k is env
	for k := range defconf.Defualtinfo.Capacity {
		app_conf.Deploy.Runtime.Ign[k] = append(app_conf.Deploy.Runtime.Ign[k], defconf.Defualtinfo.Cmdarg.Ign[k]...)
		app_conf.Deploy.Runtime.Ign[k] = removeDuplicateElement(app_conf.Deploy.Runtime.Ign[k])
	}
	// for k, v := range app_conf.Deploy.Runtime.Ign {
	// 	app_conf.Deploy.Runtime.Ign[k] = append(v, defconf.Defualtinfo.Cmdarg.Ign[k]...)
	// }

	//use app cmd,args,pkgconf settings if existed
	//or else use default config
	// if app_conf.Deploy.Build.Appyml == "" {
	// 	log.Panic("Deploy.Build.Appyml can't be empty")
	// }
	if app_conf.Deploy.Build.Cmd == "" {
		app_conf.Deploy.Build.Cmd = defconf.Defualtinfo.Build[app_conf.Application.Type].Cmd
	}
	if app_conf.Deploy.Build.Args == "" {
		app_conf.Deploy.Build.Args = defconf.Defualtinfo.Build[app_conf.Application.Type].Arg
	}
	if app_conf.Deploy.Build.Pkgconf == "" {
		app_conf.Deploy.Build.Pkgconf = defconf.Defualtinfo.Build[app_conf.Application.Type].Config
	}

	if app_conf.Deploy.Build.Output == "" {
		log.Panic("Deploy.Build.Output should not be empty")
	}

	if strings.HasPrefix(app_conf.Deploy.Build.Output, "/") {
		log.Panic("Deploy.Build.Output should not start with /")
	}

	for _, v := range defconf.Defualtinfo.Build[app_conf.Application.Type].Jenkignor {
		flag := false
		for _, exv := range app_conf.Deploy.Build.Jenkexec {
			if v == exv {
				flag = true
				break
			}
		}
		if flag {
			continue
		}
		app_conf.Deploy.Build.Jenkignor = append(app_conf.Deploy.Build.Jenkignor, v)
	}

	// pos := strings.LastIndex(app_conf.Deploy.Build.Output, "/")
	// if pos < 0 {

	// } else {

	// }
	// app_conf.Deploy.Build.OutputPath = app_conf.Deploy.Build.Output[0:pos]
	// outputtail := app_conf.Deploy.Build.Output[pos+1 : len(app_conf.Deploy.Build.Output)]
	// if !strings.Contains(outputtail, ".") {
	// 	log.Panic("Deploy.Build.Output should contain file extension")
	// }

	app_conf.Deploy.Stratail = make(map[string][][]EnvInfo)
	for _, v := range defconf.Defualtinfo.Deploy.Strategy {
		for _, env := range v.Env {
			if _, ok := app_conf.Deploy.Stratail[env]; ok {
				log.Panicf("default deploy strategy should not with duplicated env config for %s", env)
			}
			// ff := [][]string{}
			for _, sflow := range strings.Split(v.Flow, "->") {
				sflows := []EnvInfo{}
				for _, flow := range strings.Split(sflow, "|") {

					flowenv := strings.Split(strings.TrimSuffix(flow, ")"), "(")
					if len(flowenv) > 1 {
						sflows = append(sflows, struct {
							Env string
							Dc  string
						}{
							Env: flowenv[1],
							Dc:  flowenv[0],
						})
						// envstr = flowenv[1] + "-" + flowenv[0]
					} else {
						sflows = append(sflows, struct {
							Env string
							Dc  string
						}{
							Env: env,
							Dc:  flowenv[0],
						})
						// envstr = k + "-" + flowenv[0]
					}
					// envstr = strings.TrimSuffix(envstr, "-LFB")
					// env = append(env, envstr)
				}
				app_conf.Deploy.Stratail[env] = append(app_conf.Deploy.Stratail[env], sflows)
				// asyncflow := strings.Split(sflow, "|")
				// app_conf.Deploy.Stratail[env] = append(app_conf.Deploy.Stratail[env], strings.Split(sflow, "|"))
			}
			// app_conf.Deploy.Stratail[env] = ff //struct{ Flow string }{Flow: v.Flow}
		}
	}

	for _, v := range app_conf.Deploy.Strategy {
		for _, env := range v.Env {

			// if _, ok := app_conf.Deploy.Stratail[env]; ok {
			// 	log.Panicf("app deploy strategy should not with duplicated env config for %s", env)
			// }
			// ff := [][]string{}
			app_conf.Deploy.Stratail[env] = [][]EnvInfo{}
			for _, sflow := range strings.Split(v.Flow, "->") {
				for _, flow := range strings.Split(sflow, "|") {
					sflows := []EnvInfo{}
					flowenv := strings.Split(strings.TrimSuffix(flow, ")"), "(")
					if len(flowenv) > 1 {
						sflows = append(sflows, struct {
							Env string
							Dc  string
						}{
							Env: flowenv[1],
							Dc:  flowenv[0],
						})
						// envstr = flowenv[1] + "-" + flowenv[0]
					} else {
						sflows = append(sflows, struct {
							Env string
							Dc  string
						}{
							Env: env,
							Dc:  flowenv[0],
						})
						// envstr = k + "-" + flowenv[0]
					}
					app_conf.Deploy.Stratail[env] = append(app_conf.Deploy.Stratail[env], sflows)
					// envstr = strings.TrimSuffix(envstr, "-LFB")
					// env = append(env, envstr)
				}
				// asyncflow := strings.Split(sflow, "|")
				// app_conf.Deploy.Stratail[env] = append(app_conf.Deploy.Stratail[env], strings.Split(sflow, "|"))
			}
			// app_conf.Deploy.Stratail[env] = ff //struct{ Flow string }{Flow: v.Flow}
		}
	}

	if isinstall {

		// app_arch_map[app_conf.Application.Name] = app_conf

		// rediscli := redisops.Pool().Get()

		// defer rediscli.Close()
		// rediscli.Do("HSET", "arch-spell-appconfig", app_conf.Application.Name, jsonbs)
		// rediscli.Do("HSET", "arch-spell-projteam-"+app_conf.Application.Team, app_conf.Application.Name, app_conf.Application.Project)
		app_conf.Install(c)

	}
	// fileops.Write()
	return app_conf
}

// func setNetworkListEachApp(oldlist []string, flag, isclear bool, listinfo []string, appname, confkey string, c context.Context) map[string]struct{} {
// 	log := logagent.Inst(c)
// 	maplist := map[string]struct{}{}
// 	if !isclear {

// 		bytes := consulhelp.GetConfigFull(confkey, c)
// 		err := json.Unmarshal(bytes, &maplist)
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 	}

// 	for _, v := range oldlist {
// 		delete(maplist, v)
// 	}
// 	delete(maplist, appname)

// 	if flag || len(listinfo) > 0 {
// 		if _, ok := maplist[appname]; !ok && flag {
// 			maplist[appname] = struct{}{}
// 		}
// 		// } else if len(listinfo) > 0 {
// 		for _, v := range listinfo {
// 			if _, ok := maplist[v]; !ok {
// 				maplist[v] = struct{}{}
// 			}
// 		}

// 		bytes, err := json.Marshal(maplist)
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 		consulhelp.PutConfigFull(confkey, bytes, c)
// 	} else if isclear {
// 		consulhelp.DelConfigFull(confkey, c)
// 	}
// 	return maplist
// }

// func setNetworkList(oldlist []string, flag, isclear bool, listinfo []string, appname, confkey string, c context.Context) map[string]struct{} {
// 	log := logagent.Inst(c)
// 	maplist := map[string]struct{}{}
// 	if !isclear {

// 		bytes := consulhelp.GetConfigFull(confkey, c)
// 		err := json.Unmarshal(bytes, &maplist)
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 	}

// 	for _, v := range oldlist {
// 		delete(maplist, v)
// 	}
// 	delete(maplist, appname)

// 	if flag || len(listinfo) > 0 {
// 		if _, ok := maplist[appname]; !ok && flag {
// 			maplist[appname] = struct{}{}
// 		}
// 		// } else if len(listinfo) > 0 {
// 		for _, v := range listinfo {
// 			if _, ok := maplist[v]; !ok {
// 				maplist[v] = struct{}{}
// 			}
// 		}

// 		bytes, err := json.Marshal(maplist)
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 		consulhelp.PutConfigFull(confkey, bytes, c)
// 	} else if isclear {
// 		consulhelp.DelConfigFull(confkey, c)
// 	}
// 	return maplist
// }

func setWhiteList(open bool, appname, confkey string, c context.Context) map[string]struct{} {
	log := logagent.InstArch(c)
	maplist := map[string]struct{}{}

	if open {
		maplist[appname] = struct{}{}

		bytes, err := json.Marshal(maplist)
		if err != nil {
			log.Panic(err)
		}
		consulhelp.PutConfigFull(confkey, bytes, c)
	} else {
		consulhelp.DelConfigFull(confkey, c)
	}
	return maplist
}

func (v *Arch_config) Install(c context.Context) {
	// oldfig := GetArchfigSin(v.Application.Name, c)
	// v.FireWallRefresh(oldfig, c)
	v.FireWallRefresh4Wthie(c)
	bytes, _ := json.MarshalIndent(v, "", "    ")

	// consulhelp.PutConfig(*constset.ConfArchPrefix, v.Application.Team, v.Application.Project, v.Application.Name, bytes)
	consulhelp.PutConfigFull(*constset.ConfArchPrefix+v.Application.Name, bytes, c)
}

func (v *Arch_config) OrgInstall(c context.Context) {
	// oldfig := GetArchfigSin(v.Application.Name, c)
	// v.FireWallRefresh(oldfig, c)
	v.FireWallRefresh4Wthie(c)
	bytes, _ := json.MarshalIndent(v, "", "    ")

	consulhelp.PutConfigFull(*constset.ConfOrgPrefix+v.Application.Name, bytes, c)
	// consulhelp.PutConfig(*constset.ConfArchPrefix, v.Application.Team, v.Application.Project, v.Application.Name, bytes)
	// consulhelp.PutConfigFull(*constset.ConfArchPrefix+v.Application.Name, bytes, c)
}

// func GenFWfile(templName string, data map[string]struct{}, c context.Context) string {
// 	// dirPth := orgconfigPth + pthSep + "arch.yaml"
// 	// str, _ := ioutil.ReadFile(dirPth)
// 	// a := Arch_config{}
// 	// yaml.Unmarshal(str, &a)
// 	// t.Log(a)
// 	name := "fabio.fw" // "jenkins." + jenfig.Type
// 	result := templ.GemplFrom(name, data, c)

// 	return result
// }

// func FireWallFlush(fwlist map[string]map[string]struct{}, c context.Context) {

// 	for k, v := range fwlist {
// 		fabioConf := GenFWfile(k, v, c)
// 		log.Print(fabioConf)
// 		bytes := []byte(fabioConf)
// 		// bytes, _ := json.MarshalIndent(v, "", "    ")
// 		consulhelp.PutConfigFull(*constset.ConfFabioPrefix+"/"+k, bytes, c)
// 	}
// }

// func (v *Arch_config) FireWallRefresh(oldfig Arch_config, c context.Context) {

// 	fwList := map[string]map[string]struct{}{}
// 	// ml := map[string]struct{}{}
// 	mlInternet := setNetworkList(oldfig.Environment.Expose.Internet.Blacklist, v.Environment.Expose.Internet.Visible, false, v.Environment.Expose.Internet.Blacklist, v.Application.Name, *constset.ConfbalckPrefix+"/internet", c)
// 	fwList["internet"] = newFunction("internet", mlInternet, fwList, c)

// 	mlIntranet := setNetworkList(oldfig.Environment.Expose.Intranet.Blacklist, v.Environment.Expose.Intranet.Visible, false, v.Environment.Expose.Intranet.Blacklist, v.Application.Name, *constset.ConfbalckPrefix+"/intranet", c)
// 	fwList["intranet"] = newFunction("intranet", mlIntranet, fwList, c)

// 	// mlClusternet := setNetworkList(v.Environment.Expose.Clusternet.Open, []string{}, v.Application.Name, *constset.ConfwhitePrefix+"/clusternet", c)
// 	setNetworkList([]string{}, v.Environment.Expose.Clusternet.Open, true, []string{}, v.Application.Name, *constset.ConfwhitePrefix+"/clusternet/"+v.Application.Name, c)
// 	// fwList["clusternet"] = mlClusternet
// 	// mlPtrnet := setNetworkList(v.Environment.Expose.Ptrnet.Open, []string{}, v.Application.Name, *constset.ConfwhitePrefix+"/ptrnet", c)
// 	setNetworkList([]string{}, v.Environment.Expose.Ptrnet.Open, true, []string{}, v.Application.Name, *constset.ConfwhitePrefix+"/ptrnet/"+v.Application.Name, c)
// 	// fwList["ptrnet"] = mlPtrnet

// 	FireWallFlush(fwList, c)
// }

func (v *Arch_config) FireWallRefresh4Wthie(c context.Context) {
	setWhiteList(v.Environment.Expose.Clusternet.Open, v.Application.Name, *constset.ConfwhitePrefix+"clusternet/"+v.Application.Name, c)
	setWhiteList(v.Environment.Expose.Ptrnet.Open, v.Application.Name, *constset.ConfwhitePrefix+"ptrnet/"+v.Application.Name, c)
}

// func newFunction(key string, ml map[string]struct{}, fwList map[string]map[string]struct{}, c context.Context) map[string]struct{} {
// 	bytes := consulhelp.GetConfigFull(*constset.ConfmanBalckPrefix+"/"+key, c)
// 	maplist := map[string]struct{}{}
// 	err := json.Unmarshal(bytes, &maplist)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	resmap := map[string]struct{}{}
// 	for k := range maplist {
// 		resmap[k] = struct{}{}
// 	}
// 	for k := range ml {
// 		resmap[k] = struct{}{}
// 	}

// 	return resmap
// 	// fwList[key] = ml
// }

func AppExist(appname string, c context.Context) (bool, []string) {
	// appconf := Arch_config{}
	key := "ops/iac/arch/" + appname
	kvs := consulhelp.GetConfigs(*constset.ConfArchPrefix, "", c)

	apps := []string{}
	for _, v := range kvs {
		if v.Key == key {
			apps = append(apps, v.Key)
		}
	}

	return len(apps) > 0, apps
}

func (v *Arch_config) RMArch(c context.Context) {
	consulhelp.DelConfigFull(*constset.ConfArchPrefix+v.Application.Name, c)
	consulhelp.DelConfigFull(*constset.ConfOrgPrefix+v.Application.Name, c)
	// consulhelp.DelConfig(*constset.ConfArchPrefix, v.Application.Team, v.Application.Project, v.Application.Name)
}

func removeDuplicateElement(target []string) []string {
	result := make([]string, 0, len(target))
	temp := map[string]struct{}{}
	var tmpstr string
	for _, item := range target {
		if _, ok := temp[item]; !ok {
			tmpstr = strings.Trim(item, " ")
			temp[tmpstr] = struct{}{}
			result = append(result, tmpstr)
		}
	}
	return result
}
