package valfig

import (
	"context"
	"strconv"
	"strings"

	"github.com/max-gui/consulagent/pkg/consulhelp"
	"github.com/max-gui/fileconvagt/pkg/convertops"
	"github.com/max-gui/spells/internal/iac/archfig"
	"github.com/max-gui/spells/internal/iac/templ"
	"github.com/max-gui/spells/internal/pkg/constset"
)

type ValuesInfo struct {
	Env        string
	Dc         string
	MemMi      int
	Replica    int
	PrePackage string
	Unsafe     bool
	Expovice   string
	Expopath   string
	ExpoSufix  string
	Manual     []struct {
		PrefixPath string `yaml:"prefixPath,omitempty"`
		SurfixPath string `yaml:"surfixPath,omitempty"`
	} `yaml:"menual,omitempty"`
	Cpu        string
	Mem        string
	Port       string
	Tags       map[string]string
	Hosts      []map[string]interface{}
	Rtargs     string
	Volumns    []map[string]interface{}
	Detectorip map[string]interface{}
	Resource   map[string][]interface{}
	Appconf    archfig.Arch_config
}

func (v *ValuesInfo) GenPort(lens int) string {
	if len(v.Appconf.Environment.Port) <= 0 {
		head := convertops.RndRangestr(1, 3, 6)
		tail := convertops.RndRangestr(4, 0, 9)
		v.Port = head + tail

	} else {
		v.Port = v.Appconf.Environment.Port
	}
	return v.Port
}

func (v *ValuesInfo) GenMemMi() int {
	v.MemMi, _ = strconv.Atoi(strings.TrimSuffix(v.Mem, "Mi"))

	return v.MemMi
}

// func MakeValuemple(apptype string, isinstall bool) *template.Template {
// 	name := "values.yaml"
// 	// templtmp := templ.GetemplFrom(Templepath()+name, name)
// 	dirPth := constset.Templepath + name
// 	//docker templ
// 	res := templ.GetemplFrom(dirPth, name, isinstall)

// 	return res

// }

func GenValfile(valfig ValuesInfo, c context.Context) string {
	// dirPth := orgconfigPth + pthSep + "arch.yaml"
	// str, _ := ioutil.ReadFile(dirPth)
	// a := Arch_config{}
	// yaml.Unmarshal(str, &a)
	// t.Log(a)
	name := "values.yaml"
	result := templ.GemplFrom(name, valfig, c)

	return result
}

func GenValfig(appconf archfig.Arch_config, envinfo archfig.EnvInfo, env_dc string, c context.Context) ValuesInfo {

	// env := strings.Split(env_dc, "-")[0]
	// var runargs = strings.Split(appconf.Deploy.Runtime.Args, " ")
	// lens := len(runargs)

	//remove ign runtime args
	var argstr string
	argstr = appconf.Deploy.Runtime.Args
	for _, v := range appconf.Deploy.Runtime.Ign[envinfo.Env] {
		argstr = strings.Replace(appconf.Deploy.Runtime.Args, v, "", -1)

		// app_conf.Deploy.Runtime.Args = strings.Replace(app_conf.Deploy.Runtime.Args, "  ", " ", -1)
	}
	// var argstr = appconf.Deploy.Runtime.Args
	// if len(appconf.Deploy.Runtime.Ign[env]) > 0 {
	// 	for _, v := range appconf.Deploy.Runtime.Ign[env] {
	// 		for _, arg := range runargs {
	// 			if arg != v {
	// 				argstr = arg + " "
	// 				// runargs = append(runargs[:index], runargs[index+1:]...)
	// 			}
	// 		}
	// 	}
	// } else {
	// 	argstr = appconf.Deploy.Runtime.Args
	// }
	// argstr := strings.Join(runargs, " ")

	valuearg := ValuesInfo{
		Replica:    appconf.Environment.Strategy[envinfo.Env].Replica,
		PrePackage: appconf.Deploy.Build.Output, // appconf.Deploy.Build.OutputPath,
		Unsafe:     appconf.Environment.Expose.Unsafe,
		Expovice:   appconf.Environment.Expose.Expovice,
		Expopath:   appconf.Environment.Expose.PrefixPath,
		ExpoSufix:  appconf.Environment.Expose.SurfixPath,
		Manual:     appconf.Environment.Expose.Manual,
		Cpu:        appconf.Environment.Strategy[envinfo.Env].Cpu,
		Mem:        appconf.Environment.Strategy[envinfo.Env].Mem,
		Tags:       appconf.Environment.Tag,
		Rtargs:     argstr,
		Env:        envinfo.Env,
		Dc:         envinfo.Dc,
		Resource:   make(map[string][]interface{}),
		Appconf:    appconf,
	}

	if _, ok := appconf.Deploy.Limited[envinfo.Dc]; ok {
		valuearg.Replica = 1
	}

	// if appconf.Environment.Expose.Secgwon {
	// 	valuearg.Expovice = valuearg.Expovice + "-vdefault"
	// } else {
	// 	valuearg.Expovice = realseName
	// }
	// if appconf.Environment.Expose.Unsafe {
	// 	valuearg.Expovice = valuearg.Expovice + "-vdefault"
	// }

	//valuearg.Env = "prod"
	// valuearg.Mem = "8192Mi"
	valuearg.GenMemMi()
	valuearg.GenPort(5)

	for typekey, ids := range appconf.Environment.Resource {
		//get config from key and env
		// var typeresources []
		for _, id := range ids {
			maptmp := consulhelp.Getconfaml(*constset.ConfResPrefix, typekey, id, env_dc, c)
			maptmp["id"] = id
			valuearg.Resource[typekey] = append(valuearg.Resource[typekey], maptmp)
			// log.Print(hostkey)
		}
	}

	iphostmap := make(map[string][]string)
	for _, v := range valuearg.Resource["hostAlias"] {
		iphost := v.(map[string]interface{})
		// ip := vv["ip"].(string)
		// hosts := vv["host"].([]string)
		// vvvvv := vv["host"].(interface{}).([]string)
		iphostmap[iphost["ip"].(string)] = append(iphostmap[iphost["ip"].(string)], iphost["host"].(string))
		// for _, iphostv := range iphost["host"].([]interface{}) {
		// 	iphostmap[iphost["ip"].(string)] = append(iphostmap[iphost["ip"].(string)], iphostv.(string))
		// }
		// for _, host := range iphost["host"].([]interface{}) {
		// 	iphostmap[iphost["ip"].(string)] = append(iphostmap[iphost["ip"].(string)], host.(string))
		// }

	}
	valuearg.Resource["hostAliasfinal"] = []interface{}{}
	for k, v := range iphostmap {
		iphostmap[k] = removeDuplicateElement(v)
		valuearg.Resource["hostAliasfinal"] = append(valuearg.Resource["hostAliasfinal"], map[string]interface{}{
			"ip":   k,
			"host": v,
		})
	}

	// detip := struct {
	// 	Ip string
	// }{}
	// detectortmp := make(map[string]interface{})
	// bytes := GetConfig("pinpoint", "detector", env_dc)
	// err := yaml.Unmarshal(bytes, &detectortmp)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// valuearg.Detectorip = detectortmp
	// log.Print(valuearg.Detectorip)

	// host := struct {
	// 	Ip   string
	// 	Host []string
	// }{}
	// for _, hostkey := range appconf.Environment.Hosts {
	// 	//get config from key and env
	// 	maptmp := make(map[string]interface{})
	// 	bytes = GetConfig("hostAlias", hostkey, env_dc)
	// 	err = yaml.Unmarshal(bytes, &maptmp)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	valuearg.Hosts = append(valuearg.Hosts, maptmp)
	// 	log.Print(hostkey)
	// }
	// vol := struct {
	// 	MountPath string `yaml:"mountPath"`
	// 	HostPath  string `yaml:"hostPath"`
	// 	Name      string
	// }{}
	// for _, volkey := range appconf.Environment.Volumns {
	// 	//get config from key and env
	// 	bytes = GetConfig("volumn", volkey, env_dc)
	// 	maptmp := make(map[string]interface{})
	// 	err = yaml.Unmarshal(bytes, &maptmp)
	// 	// vol.Name = volkey
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	valuearg.Volumns = append(valuearg.Volumns, maptmp)
	// 	log.Print(volkey)
	// }

	return valuearg

}

func removeDuplicateElement(target []string) []string {
	result := make([]string, 0, len(target))
	temp := map[string]struct{}{}
	for _, item := range target {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
