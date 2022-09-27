package iac

import (
	"bufio"
	"context"
	"os"
	"strings"
	"time"

	"github.com/max-gui/consulagent/pkg/consulhelp"
	"github.com/max-gui/spells/internal/githelp"
	"github.com/max-gui/spells/internal/iac/archfig"
	"github.com/max-gui/spells/internal/iac/defig"
	"github.com/max-gui/spells/internal/iac/templ"
	"github.com/max-gui/spells/internal/pkg/constset"
	"gopkg.in/yaml.v2"

	"log"
)

// type AppconfDist map[string]Arch_config

// var app_arch_map = make(map[string]Arch_config)

func isBranchName(branch string) bool {
	// stringTime := "20210930"

	loc, _ := time.LoadLocation("Asia/Shanghai")

	_, err := time.ParseInLocation("20060102", strings.ReplaceAll(branch, "release", ""), loc)
	if err != nil {
		log.Print(err)
		return false
	}

	return true
}

func IsBranchNameIllegal(branch string, env string) bool {
	if !isBranchName(branch) && env == "prod" {
		return true
	}
	return false
}

func GetTeamApps(team string, c context.Context) map[string][]string {
	pairs := consulhelp.GetConfigs(*constset.ConfArchPrefix, team, c)
	projapps := make(map[string][]string, len(pairs))
	for _, v := range pairs {
		infos := strings.Split(strings.TrimPrefix(v.Key, "ops/resource/"), "/")
		proj := infos[1]
		appname := infos[2]
		if projapps[proj] == nil {
			projapps[proj] = []string{}
		}
		// appconf := GetAppconfig(appname, team, proj)
		projapps[proj] = append(projapps[proj], appname)
	}

	return projapps
}

func ClsAppfig() {

	consulhelp.ClsConfig()
	// app_arch_map = make(map[string]Arch_config)
	// consulhelp.DelConfig
	// rediscli := redisops.Pool().Get()

	// defer rediscli.Close()

	// rediscli.Do("DEL", "arch-spell-appconfig") //, app_conf.Application.Name, jsonbs)

	// rediscli.Do("HSET", "arch-spell-projteam", app_conf.Application.Project, app_conf.Application.Team)
}

func Clearcachelocal() {
	ClsAppfig()
	templ.ClsGempl()
}

func ClearcacheAll(c context.Context) {
	githelp.UpdateAll(c)
	ClsAppfig()
	templ.ClsGempl()
	defig.ClearDefconfig()
}

func GetKeyres(keyenvs []struct {
	Key string
	Env string
}, c context.Context) []map[string]string {
	kvs := consulhelp.GetConfigs("ops/resource", "", c)

	keymaps := []map[string]string{}
	for _, v := range kvs {
		// hostkey := strings.Split(strings.TrimPrefix(v.Key, "ops/resource"), "/")[0]
		if strings.Contains(v.Key, "LogConfig") || strings.Contains(v.Key, "bootload") {
			continue
		}
		valuestr := string(v.Value)

		for _, keyenv := range keyenvs {
			if strings.Contains(valuestr, keyenv.Key) {
				keymap := map[string]string{}
				keymap["key"] = v.Key
				keymap["env"] = keyenv.Env
				keymap["value"] = valuestr
				keymaps = append(keymaps, keymap)
			}
		}
		// if _, ok := hostmaptmp["real-id"]; !ok {
		// 	hostmap[hostmaptmp["host"]] = hostkey
		// }
	}
	return keymaps
}

func getHosts(c context.Context) map[string]string {
	kvs := consulhelp.GetConfigs("ops/resource/hostAlias", "", c)

	hostmap := map[string]string{}
	for _, v := range kvs {
		hostkey := strings.Split(strings.TrimPrefix(v.Key, "ops/resource/hostAlias/"), "/")[0]
		hostmaptmp := make(map[string]string)
		yaml.Unmarshal(v.Value, hostmaptmp)

		if hostv, ok := hostmaptmp["host"]; ok {
			hostmap[hostv] = hostkey
		}
		// if _, ok := hostmaptmp["real-id"]; !ok {
		// 	hostmap[hostmaptmp["host"]] = hostkey
		// }
	}
	return hostmap
}

func getVolumns(c context.Context) map[string]string {
	kvs := consulhelp.GetConfigs("ops/resource/volumn", "", c)

	volumnmap := map[string]string{}
	for _, v := range kvs {
		volumnname := strings.Split(strings.TrimPrefix(v.Key, "ops/resource/volumn/"), "/")[0]
		volumnmap[volumnname] = volumnname

	}
	return volumnmap
}

func Gen4old(appid, appname string, c context.Context) archfig.Arch_config {
	hostmap := getHosts(c)
	volumnmap := getVolumns(c)
	iacPath := constset.Iacpath + "app" + constset.PthSep + appid + constset.PthSep + appname + constset.PthSep
	var archconf archfig.Arch_config
	// dir, _ = ioutil.ReadDir(iacDir)
	{
		iacDir, _ := os.ReadDir(iacPath)
		var AppId, giturl, Cmd, CmdArgs, Pom, Output, AppName, Env, PrePackage, Cpu, Mem, Rtargs string
		var Replica int
		var ExpoviceOk bool
		AppId = appid

		valuesm := make(map[string]interface{})
		stratgym := make(map[string]struct {
			Capacity string `yaml:"capacity,omitempty"`
			Cpu      string `yaml:"cpu,omitempty"`
			Mem      string `yaml:"mem,omitempty"`
			Replica  int    `yaml:"replica,omitempty"`
		})
		resources := make(map[string]map[string]struct{})
		for _, iacfi := range iacDir { //app

			log.Print("===================================filename===================================================")
			// vfilepath := dirPth + fi.Name() + string(os.PathSeparator) + subfi.Name() + string(os.PathSeparator) + iacfi.Name()
			iacfilepath := iacPath + iacfi.Name()
			log.Print(iacfilepath)

			if strings.Contains(iacfi.Name(), "values") {
				if iacfi.Name() == "values-prod.yaml" {
					Env = "prod"
				} else if iacfi.Name() == "values-test.yaml" {
					Env = "test"
				} else if iacfi.Name() == "values-uat.yaml" {
					Env = "uat"
				} else if iacfi.Name() == "values-dr.yaml" {
					Env = "dr"
				} else {
					continue
				}
				bytes, err := os.ReadFile(iacfilepath)
				// file, err := os.Open(dirPth + fi.Name() + string(os.PathSeparator) + subfi.Name() + string(os.PathSeparator) + iacfi.Name())
				if err != nil {
					log.Panic(err)
				}
				valuesm = make(map[string]interface{})
				yaml.Unmarshal(bytes, &valuesm)
				log.Print(valuesm)
				// valuearg := ValuesInfo{
				// 	Replica:    appconf.Environment.Strategy[env].Replica,
				// 	PrePackage: appconf.Deploy.Build.OutputPath,
				// 	Expovice:   appconf.Environment.Expose.Service,
				// 	Cpu:        appconf.Environment.Strategy[env].Cpu,
				// 	Mem:        appconf.Environment.Strategy[env].Mem,
				// 	Tags:       appconf.Environment.Tag,
				// 	Rtargs:     argstr,
				// 	Env:        env,
				// 	Resource:   make(map[string][]interface{}),
				// }
				Replica = valuesm["replicaCount"].(int)
				if val, ok := valuesm["webPreDir"]; ok {
					PrePackage = val.(string)
				}
				_, ExpoviceOk = valuesm["ingress"]
				if valuesm["resources"] == nil {
					// nonvalues = append(nonvalues, vfilepath)
				} else {
					if val, ok := valuesm["resources"].(map[interface{}]interface{})["limits"].(map[interface{}]interface{})["cpu"]; ok {
						Cpu = val.(string)
					} else {
						Cpu = "0.5"
					}

					// Cpu = valuesm["resources"].(map[interface{}]interface{})["limits"].(map[interface{}]interface{})["cpu"].(string)
					Mem = valuesm["resources"].(map[interface{}]interface{})["limits"].(map[interface{}]interface{})["memory"].(string)
					if val, ok := valuesm["java_opts"]; ok {
						Rtargs = val.(string)
					}
				}
				// 							hostAliases:
				//   - ip: "10.47.52.70"
				//     hostnames:
				stratgym[Env] = struct {
					Capacity string `yaml:"capacity,omitempty"`
					Cpu      string `yaml:"cpu,omitempty"`
					Mem      string `yaml:"mem,omitempty"`
					Replica  int    `yaml:"replica,omitempty"`
				}{
					Cpu:     Cpu,
					Mem:     Mem,
					Replica: Replica,
				}
				// resources = make(map[string][]string)
				// log.Print(valuesm["hostAliases"])
				if val, ok := valuesm["hostAliases"]; ok {
					for _, v := range val.([]interface{}) {

						// log.Print(v)
						vv := v.(map[interface{}]interface{})
						// ip := vv["ip"].(string)

						hostnames := vv["hostnames"].([]interface{})

						for _, hs := range hostnames {
							if hostval, hostvalok := hostmap[hs.(string)]; hostvalok {
								// resid := mmmttt[hs+ip]
								if resources["hostAlias"] == nil {
									resources["hostAlias"] = map[string]struct{}{}
								}
								resources["hostAlias"][hostval] = struct{}{} //append(resources["hostAlias"], hostval)
							} else {
								ip := vv["ip"].(string)
								log.Panicf("hostAlias does not exist! host:%s ip:%s env:%s", hs, ip, Env)
							}

						}
					}
				}

				if val, ok := valuesm["volumeMounts"]; ok {
					for _, v := range val.([]interface{}) {
						vv := v.(map[interface{}]interface{})
						volname := vv["name"].(string)
						if _, ok := volumnmap[volname]; !ok {
							log.Panicf("volumn does not exist! name:%s env:%s", volname, Env)
						}

						for _, vvv := range valuesm["volumes"].([]interface{}) {
							vvvv := vvv.(map[interface{}]interface{})
							if vvvv["name"].(string) == volname {
								// log.Print(vvvv)
								log.Print(vvvv)

								break
							}
						}
						// mountPath: /wls/wls81/logs
						// hostPath: /nfsc/cnas_csp_stg_fls_aflm_id9192_vol1003_stg/logs/test
						// name: logs
						// valuesm["volumes"]
						// mountPath := valuesm["volumeMounts"]["mountPath"].(string)
						if resources["volumn"] == nil {
							resources["volumn"] = map[string]struct{}{}
						}
						resources["volumn"][volname] = struct{}{} // = append(resources["volumn"], volname)
					}
				}
			} else if iacfi.Name() == "Dockerfile" {
				bytes, err := os.ReadFile(iacfilepath)
				if err != nil {
					log.Panic(err)
				}
				x0 := string(bytes)
				if strings.Contains(x0, "pazl-web.war") {
					Output = "pazl-web.war"
				} else if strings.Contains(string(bytes), "dist") {
					Output = "dist"
				} else if strings.Contains(string(bytes), "test-platform-server") {
					Output = "./"
				} else if strings.Contains(string(bytes), "javaconv.jar") {
					Output = "./"
				}
			} else if iacfi.Name() == "Jenkinsfile" {
				file, err := os.Open(iacfilepath)
				if err != nil {
					log.Panic(err)
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					str := scanner.Text()
					if strings.Contains(str, "GitRepositoryURL") {
						strarr := strings.Split(str, ",")
						for _, v := range strarr {
							if strings.Contains(v, "defaultValue") {
								gitv := strings.TrimPrefix(strings.TrimSpace(v), "defaultValue:")
								// gitv := strings.Split(v, ":")
								giturl = strings.Trim(strings.TrimSpace(gitv), "'")
								break
							}
						}
					} else {
						vars := strings.Split(str, "=")
						switch strings.TrimSpace(vars[0]) {
						case "buildType":
							Cmd = strings.TrimSpace(vars[1])
						case "buildShell":
							CmdArgs = strings.TrimSpace(vars[1])
						case "Pom":
							Pom = strings.Trim(strings.TrimSpace(vars[1]), "'")
						// case "classes":
						// 	break
						case "appJarDir":
							if strings.Trim(strings.TrimSpace(vars[1]), "'") != "" {
								Output = strings.Trim(strings.TrimSpace(vars[1]), "'")
							} else {
								log.Print(Output)
							}
						// case "appId":
						// 	break //AppId = strings.TrimSpace(vars[1])
						// case "HelmDir":
						// 	break
						case "ServiceName":
							AppName = strings.Trim(strings.TrimSpace(vars[1]), "'")
						}
					}
				}

				if err := scanner.Err(); err != nil {
					log.Fatal(err)
				}
			}

		}
		archconf = archfig.Arch_config{}
		archconf.Application.Appid = AppId
		archconf.Application.Name = AppName
		archconf.Application.Project = ""
		archconf.Application.Repositry = giturl
		archconf.Application.Team = ""
		if PrePackage == "" {
			archconf.Application.Type = "java"
		} else {
			archconf.Application.Type = "h5"

			if ExpoviceOk {
				archconf.Environment.Expose.Expovice = PrePackage
				archconf.Application.Name = PrePackage
			}
		}
		for rk, rv := range resources {
			var rvks = []string{}
			for rvk := range rv {
				rvks = append(rvks, rvk)
			}
			if archconf.Environment.Resource == nil {
				archconf.Environment.Resource = map[string][]string{}
			}
			archconf.Environment.Resource[rk] = rvks
		}
		// archfig.Environment.Resource = resources
		archconf.Environment.Strategy = stratgym
		archconf.Deploy.Build.Args = strings.Trim(CmdArgs, "'")
		archconf.Deploy.Build.Cmd = strings.Trim(Cmd, "'")
		archconf.Deploy.Build.Output = strings.Trim(Output, "'")
		archconf.Deploy.Build.Pkgconf = Pom
		archconf.Deploy.Runtime.Args = strings.Split(Rtargs, " ")
		archconf.Application.Ungenfig = true

		// archfigs = append(archfigs, archfig)
		// bbbb, _ := yaml.Marshal(archfig)
		// log.Print(string(bbbb))
		// archfigfull := GenArchConfig(bbbb, "unknown", "unknown", archfig.Application.Name, false)
		// GenJenfile(GenJenfig(archfigfull))
		// // log.Print(GenJenfile(GenJenfig(archfigfull)))
		// for _, env := range []string{"prod", "uat", "dr", "test"} {
		// 	GenValfile(GenValfig(archfigfull, env, env))
		// 	// log.Print(GenValfile(GenValfig(archfigfull, env, env)))
		// 	log.Print("++++++++++++++++++++++++++++++++++++++++++++++++++++")
		// }
		// rediscli.Do("SET", dirPth+fi.Name()+string(os.PathSeparator)+subfi.Name(), "skip")
	}

	return archconf
}
