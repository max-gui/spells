package deploy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"

	"sync"
	"time"

	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/logagent/pkg/logsets"
	"github.com/max-gui/spells/internal/githelp"
	"github.com/max-gui/spells/internal/iac"
	"github.com/max-gui/spells/internal/iac/altconfig"
	"github.com/max-gui/spells/internal/iac/archfig"
	"github.com/max-gui/spells/internal/iac/dockfig"
	"github.com/max-gui/spells/internal/iac/valfig"
	"github.com/max-gui/spells/internal/pkg/constset"
	"github.com/max-gui/spells/internal/pkg/jenkinsops"
)

type DeployInfo struct {
	Appname string
	Env     string
	Region  string
	Branch  string
	Team    string
	Proj    string
	Dc      string
}

type dcenvdeploy func(archfig.Arch_config, string, string, string, string, string, bool, context.Context) [][]map[string]string
type opsmethod func(map[string]string, bool, bool, context.Context) Deployresult

// func FlashNdeploy4Strtegy(branch, env, dc, appname, team, proj, region string, c context.Context, genconf bool) []deployresult {
// 	return FlashNOps(branch, env, dc, appname, team, proj, region, deploy4strategy, genconf, c)
// }

//	func FlashNdeploy4Target(branch, env, dc, appname, team, proj, region string, c context.Context, genconf bool) []deployresult {
//		return FlashNOps(branch, env, dc, appname, team, proj, region, deploy4target, genconf, c)
//	}
func StrtegyFlashDeploy(branch, env, dc, appname, team, proj, region string, c context.Context) []Deployresult {
	return FlashNOps(branch, env, dc, "", appname, team, proj, region, true, deploy4strategy, deploySingle, true, c)
}

func TargetFlashDeploy(branch, env, dc, appname, team, proj, region string, c context.Context) []Deployresult {
	return FlashNOps(branch, env, dc, "", appname, team, proj, region, true, deploy4target, deploySingle, true, c)
}

func StrtegyFlashRelease(branch, env, dc, appname, team, proj, region string, c context.Context) []Deployresult {
	return FlashNOps(branch, env, dc, "", appname, team, proj, region, false, deploy4strategy, releaseSingle, false, c)
}

func TargetFlashRelease(branch, env, dc, appname, team, proj, region string, c context.Context) []Deployresult {
	return FlashNOps(branch, env, dc, "", appname, team, proj, region, false, deploy4target, releaseSingle, false, c)
}

func TargetFlashMigrate(branch, env, dc, sourcedc, appname, team, proj, region string, c context.Context) []Deployresult {
	return FlashNOps(branch, env, dc, sourcedc, appname, team, proj, region, false, deploy4target, migrateSingle, false, c)
}

func FlashNOps(branch, env, dc, sourcedc, appname, team, proj, region string, updateRepo bool, deployhelp dcenvdeploy, opsingl opsmethod, genconf bool, c context.Context) []Deployresult {

	log := logagent.InstPlatform(c)

	if iac.IsBranchNameIllegal(branch, env) {
		log.Panic("branch name is ilegal")
	}
	var (
		isupdate bool
		appconf  archfig.Arch_config
	)
	if updateRepo {
		// appconf := archfig.GetAppconfig(appname, team, proj)
		gitres := githelp.UpdateAll(c)
		iacr := gitres[constset.Iacname].Repo
		// _, err := githelp.CloneGetrepo(*constset.Archurl, constset.Archpath)
		// if err != nil && err != git.NoErrAlreadyUpToDate {
		// 	log.Panic(err)
		// }
		// _, err = githelp.CloneGetrepo(*constset.Templurl, constset.Templepath)
		// if err != nil && err != git.NoErrAlreadyUpToDate {
		// 	log.Panic(err)
		// }
		// iacr, err := githelp.CloneGetrepo(*constset.IacUrl, constset.Iacpath)
		// if err != nil && err != git.NoErrAlreadyUpToDate {
		// 	log.Panic(err)
		// }

		// var isupdate = false
		// var regenflag = false
		// if archerr != git.NoErrAlreadyUpToDate || templerr != git.NoErrAlreadyUpToDate {
		// for _, v := range gitres {
		// 	if v.Isupdate {
		// 		regenflag = true
		// 		break
		// 	}
		// }
		iac.Clearcachelocal()

		appconf = archfig.GetAppconfigFull(appname, true, c)

		filesinfo, _ := altconfig.ArchAltGenWithChanges(appconf, true, false, c)
		appconf.Install(c)
		isupdate = githelp.CommitPushFiles(filesinfo, iacr, constset.Iacpath, c)
		// }
	} else {
		appconf = archfig.GetAppconfigFull(appname, true, c)
		isupdate = false
	}

	deploymap := deployhelp(appconf, env, dc, sourcedc, region, branch, genconf, c)
	// var deploymap []map[string]string

	resurl := opsMultiply(deploymap, isupdate, opsingl, c)
	return resurl
}

func FlashNdeploy(branch, env, dc, sourcedc, appname, team, proj, region string, deployhelp dcenvdeploy, c context.Context) []Deployresult {

	log := logagent.InstPlatform(c)

	if iac.IsBranchNameIllegal(branch, env) {
		log.Panic("branch name is ilegal")
	}
	// appconf := archfig.GetAppconfig(appname, team, proj)
	gitres := githelp.UpdateAll(c)
	iacr := gitres[constset.Iacname].Repo
	// _, err := githelp.CloneGetrepo(*constset.Archurl, constset.Archpath)
	// if err != nil && err != git.NoErrAlreadyUpToDate {
	// 	log.Panic(err)
	// }
	// _, err = githelp.CloneGetrepo(*constset.Templurl, constset.Templepath)
	// if err != nil && err != git.NoErrAlreadyUpToDate {
	// 	log.Panic(err)
	// }
	// iacr, err := githelp.CloneGetrepo(*constset.IacUrl, constset.Iacpath)
	// if err != nil && err != git.NoErrAlreadyUpToDate {
	// 	log.Panic(err)
	// }

	// var isupdate = false
	// var regenflag = false
	// if archerr != git.NoErrAlreadyUpToDate || templerr != git.NoErrAlreadyUpToDate {
	//
	//	for _, v := range gitres {
	//		if v.Isupdate {
	//			regenflag = true
	//			break
	//		}
	//	}
	iac.Clearcachelocal()
	appconf := archfig.GetAppconfigFull(appname, true, c)

	filesinfo, _ := altconfig.ArchAltGenWithChanges(appconf, true, false, c)
	appconf.Install(c)
	isupdate := githelp.CommitPushFiles(filesinfo, iacr, constset.Iacpath, c)

	deploymap := deployhelp(appconf, env, dc, sourcedc, region, branch, true, c)
	// var deploymap []map[string]string

	resurl := opsMultiply(deploymap, isupdate, deploySingle, c)
	return resurl
}

func FlashNrelease(branch, env, dc, sourcedc, appname, team, proj, region string, updateRepo bool, deployhelp dcenvdeploy, opsingl opsmethod, c context.Context) []Deployresult {

	log := logagent.InstPlatform(c)

	if iac.IsBranchNameIllegal(branch, env) {
		log.Panic("branch name is ilegal")
	}
	var (
		isupdate bool
		appconf  archfig.Arch_config
	)
	if updateRepo {
		// appconf := archfig.GetAppconfig(appname, team, proj)
		gitres := githelp.UpdateAll(c)
		iacr := gitres[constset.Iacname].Repo
		// _, err := githelp.CloneGetrepo(*constset.Archurl, constset.Archpath)
		// if err != nil && err != git.NoErrAlreadyUpToDate {
		// 	log.Panic(err)
		// }
		// _, err = githelp.CloneGetrepo(*constset.Templurl, constset.Templepath)
		// if err != nil && err != git.NoErrAlreadyUpToDate {
		// 	log.Panic(err)
		// }
		// iacr, err := githelp.CloneGetrepo(*constset.IacUrl, constset.Iacpath)
		// if err != nil && err != git.NoErrAlreadyUpToDate {
		// 	log.Panic(err)
		// }

		// var isupdate = false
		// var regenflag = false
		// if archerr != git.NoErrAlreadyUpToDate || templerr != git.NoErrAlreadyUpToDate {
		// for _, v := range gitres {
		// 	if v.Isupdate {
		// 		regenflag = true
		// 		break
		// 	}
		// }
		iac.Clearcachelocal()

		appconf = archfig.GetAppconfigFull(appname, true, c)

		filesinfo, _ := altconfig.ArchAltGenWithChanges(appconf, true, false, c)
		appconf.Install(c)
		isupdate = githelp.CommitPushFiles(filesinfo, iacr, constset.Iacpath, c)
		// }
	} else {
		appconf = archfig.GetAppconfigFull(appname, true, c)
		isupdate = false
	}

	deploymap := deployhelp(appconf, env, dc, sourcedc, region, branch, false, c)
	// var deploymap []map[string]string

	resurl := opsMultiply(deploymap, isupdate, opsingl, c)
	return resurl
}

func deploy4target(appconf archfig.Arch_config,
	env, dc, sourcedc, region, branch string, genconf bool, c context.Context) [][]map[string]string {
	// valstring, dockerstring := getvalues(appconf, BuildEnv, prfileActive)
	// deploymap = append(deploymap, map[string]string{
	// 	"BuildEnv":         BuildEnv,
	// 	"jenkinsNode":      "build-slave-node2",
	// 	"GitRepositoryURL": appconf.Application.Repositry,
	// 	"realseName":       realseName,
	// 	"AfVersion":        region,
	// 	"GitBranch":        branch,
	// 	"prfileActive":     prfileActive,
	// 	"dc":               dc,
	// 	"Valstring":        valstring,
	// 	"Dockerstring":     dockerstring,
	// 	"Description":      appconf.Application.Description,
	// 	"Appname":          appname,

	// 	"Appid": appconf.Application.Appid,
	// })

	// var depmap map[string]string

	depmap := genjenkloymap(appconf,
		archfig.EnvInfo{
			Env:      env,
			Dc:       dc,
			SourceDc: sourcedc,
		}, region, branch, c)
	// realseName, depmap = fn2(appconf, BuildEnv, envstr, realseName, region, branch)
	// depmap["BuildEnv"] = env
	// depmap["prfileActive"] = envdc
	// depmap["dc"] = dc
	deploymapseq := [][]map[string]string{}
	deploymapseq = append(deploymapseq, []map[string]string{depmap})

	confAppend(genconf, appconf, c)

	return deploymapseq
}

// deployhelp(appconf, env, dc, region, branch)
func deploy4strategy(appconf archfig.Arch_config, env, dc, sourcedc, region, branch string, genconf bool, c context.Context) [][]map[string]string {
	deploymapseq := [][]map[string]string{}
	// envstrlist := []string{}
	for _, v := range appconf.Deploy.Stratail[env] {
		deploymaps := []map[string]string{}
		for _, envinfo := range v {
			// envstr := envinfo.Env + "-" + envinfo.Dc
			// envstr = strings.TrimSuffix(envstr, "-LFB")
			// "BuildEnv":         env,
			// "prfileActive":     envstr,
			// "dc":               envinfo.Dc,
			// var depmap map[string]string
			if envinfo.Env != "prod" && region == "rc" {
				continue
			}
			envinfo.SourceDc = sourcedc
			depmap := genjenkloymap(appconf, envinfo, region, branch, c)
			// depmap["BuildEnv"] = envinfo.Env
			// depmap["prfileActive"] = envdc
			// depmap["dc"] = envinfo.Dc
			deploymaps = append(deploymaps, depmap)
			// envstrlist = append(envstrlist, depmap["prfileActive"])
			// envstrlist = append(envstrlist,depmap["prfileActive"]}
			// deploymap = append(deploymap, map[string]string{
			// 	"BuildEnv":         env,
			// 	"jenkinsNode":      "build-slave-node2",
			// 	"GitRepositoryURL": appconf.Application.Repositry,
			// 	"realseName":       realseName,
			// 	"AfVersion":        region,
			// 	"GitBranch":        branch,
			// 	"prfileActive":     envstr,
			// 	"dc":               envinfo.Dc,
			// 	"Valstring":        valstring,
			// 	"Dockerstring":     dockerstring,
			// 	"Description":      appconf.Application.Description,
			// 	"Appname":          appname,

			// 	"Appid": appconf.Application.Appid,
			// })
		}
		deploymapseq = append(deploymapseq, deploymaps)
	}

	confAppend(genconf, appconf, c)

	// var realseName string
	// if envinfo.Env == "prod" || envinfo.Env == "dr" {
	// 	realseName = appconf.Application.Name + ""
	// } else {
	// 	realseName = appconf.Application.Name + "-v" + region

	return deploymapseq
}

func confAppend(genconf bool, appconf archfig.Arch_config, c context.Context) {
	if !genconf {
		return
	}
	// defig.
	// var genEnvstrlist = func(appconfig archfig.Arch_config) []string {
	// 	envdcstrmap := map[string]interface{}{}
	// 	for _, v := range appconfig.Deploy.Stratail {
	// 		for _, vele := range v {
	// 			for _, envdcinfo := range vele {
	// 				envdcstrmap[envdcinfo.Env+envdcinfo.Dc] = struct{}{}
	// 			}
	// 		}
	// 	}

	// 	envstrlist := []string{}
	// 	for k, _ := range envdcstrmap {
	// 		envstrlist = append(envstrlist, k)
	// 	}
	// 	return envstrlist
	// }

	// envstrlist := genEnvstrlist(appconf)
	genConfAppend(appconf, appconf.Deploy.Dcenv, c)
}

// func genjenkloymapold(appconf archfig.Arch_config,
// 	env, dc, region, branch string) (string, map[string]string) {
// 	envstr := env + "-" + dc
// 	envstr = strings.TrimSuffix(envstr, "-LFB")
// 	valstring, dockerstring := getDeployconf(appconf, env, envstr)

// 	log.Print(valstring)
// 	var realseName string
// 	if env == "prod" || env == "dr" {
// 		realseName = appconf.Application.Name + ""
// 	} else {
// 		realseName = appconf.Application.Name + "v" + region
// 	}

// 	depmap := map[string]string{

// 		"jenkinsNode":      "build-slave-node2",
// 		"GitRepositoryURL": appconf.Application.Repositry,
// 		"realseName":       realseName,
// 		"AfVersion":        region,
// 		"GitBranch":        branch,

// 		"Valstring":    valstring,
// 		"Dockerstring": dockerstring,
// 		"Description":  appconf.Application.Description,
// 		"Appname":      appconf.Application.Name,

// 		"Appid": appconf.Application.Appid,
// 	}

// 	return envstr, depmap
// }

type ConfsolverIac struct {
	Data    map[string]map[string]string
	Error   string
	Message string
}

func genjenkloymap(appconf archfig.Arch_config,
	envinfo archfig.EnvInfo, region, branch string, c context.Context) map[string]string {
	// envstr := envinfo.Env + "-" + envinfo.Dc

	envstr := getenvstr(envinfo)

	// var realseName string
	// if envinfo.Env == "prod" || envinfo.Env == "dr" {
	// 	realseName = appconf.Application.Name + ""
	// } else {
	// 	realseName = appconf.Application.Name + "-v" + region

	//curl http://arch-spells/generateConfig/iac/
	// apolloEnv
	// apolloCluster
	// "apolloEnv":     apolloEnv,
	// "apolloCluster": apolloCluster,
	return genjenkinsmap(appconf, envinfo, envstr, region, branch, c)
}

func genjenkinsmap(appconf archfig.Arch_config, envinfo archfig.EnvInfo, envstr string, region string, branch string, c context.Context) map[string]string {

	var valstring, dockerstring string

	log := logagent.InstPlatform(c)

	// }
	realseName := appconf.Application.Name + "-" + envinfo.Env + "-" + region

	if !appconf.Application.Ungenfig {

		valstring, dockerstring = getDeployconf(appconf, envinfo, envstr, c)
		// tags := map[string]interface{}{}
		// err := yaml.Unmarshal([]byte(valstring), tags)
		// if err != nil {
		// 	log.Panic(err)
		// }

		// rediscli := redisops.Pool().Get()

		// defer rediscli.Close()
		// jsonbs, _ := json.Marshal(tags["podtags"])

		// _, err = rediscli.Do("SETEX", "spells-"+realseName+"-tags", 600, jsonbs)
		// if err != nil {
		// 	log.Panic(err)
		// }
	}
	log.Print(valstring)
	bytes, _ := json.MarshalIndent(appconf, "", "    ")
	depmap := map[string]string{

		"jenkinsNode":      "build-slave-node2",
		"GitRepositoryURL": appconf.Application.Repositry,
		"realseName":       realseName,
		"AfVersion":        region,
		"GitBranch":        branch,

		"Valstring":    valstring,
		"Dockerstring": dockerstring,
		"Description":  appconf.Application.Description,
		"Appname":      appconf.Application.Name,
		"team":         appconf.Application.Team,
		"Langval":      appconf.Application.Langval,

		"Appid": appconf.Application.Appid,

		"BuildEnv":     envinfo.Env,
		"prfileActive": envstr,
		"dc":           envinfo.Dc,
		"sourcedc":     envinfo.SourceDc,
		"iacenv":       *logsets.Appenv,
		"isupdate":     "false",
		"archfig":      string(bytes),
		"Unbuild":      strconv.FormatBool(envinfo.Dc != *logsets.Appdc),
	}

	return depmap
}

func genConfAppend(appconf archfig.Arch_config, envstrlist []string, c context.Context) {
	log := logagent.InstPlatform(c)

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	// values := []string{""}
	json_data, err := json.Marshal(envstrlist)

	if err != nil {
		log.Panic(err)
	}

	// resp, err := http.Post("https://httpbin.org/post", "application/json",
	// 	bytes.NewBuffer(json_data))

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var res map[string]interface{}

	// json.NewDecoder(resp.Body).Decode(&res)

	// fmt.Println(res["json"])

	response, err := netClient.Post("http://"+*constset.Consolvername+"/conf/gen/"+appconf.Application.Name, "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		log.Panic(err)
	}
	resbody, err := io.ReadAll(response.Body)
	if err != nil || response.StatusCode != 200 {
		log.Panic("consolver gen :" + string(resbody))
	}
	var resjson = ConfsolverIac{}
	err = json.Unmarshal(resbody, &resjson)
	if err != nil {
		log.Panic(err)
	}
	if resjson.Error != "" {
		log.Panic(err)
	}
	if val, ok := resjson.Data[envstrlist[0]]; ok {
		for k, v := range val {
			appconf.Environment.Tag[k] = v
		}
	}

}

func getenvstr(envinfo archfig.EnvInfo) string {
	envstr := envinfo.Env + envinfo.Dc
	// envstr = strings.TrimSuffix(envstr, "LFB")
	return envstr
}

type Deployresult struct {
	Resulturl string `json:"resulturl"`
	JobName   string `json:"jobName"`
	TaskIndex int64  `json:"taskIndex"`
	Dcenv     string `json:"dcenv"`
	Dc        string `json:"dc"`
	Status    bool   `json:"status"`
	Msg       string `json:"msg"`
}

// func deploySingle(deployinfo map[string]string, isupdate, flag bool, c chan deployresult, wg *sync.WaitGroup) {
// 	// var resurl deployresult
// 	depres := deployresult{Dcenv: deployinfo["BuildEnv"], Dc: deployinfo["dc"]}
// 	if !flag {
// 		depres.Status = false
// 		depres.Msg = "deploy in front is error"
// 		c <- depres
// 		wg.Done()
// 		return
// 	}

// 	// restmp := deployresult{status: false, msg: fmt.Sprintf("%v", e)}
// 	ctx := context.Background()
// 	jenkins, jurl, _ := jenkinsops.GetJenkins(ctx, deployinfo["BuildEnv"])

// 	job, err := jenkins.GetJob(ctx, "iac-"+deployinfo["Appname"])
// 	if err != nil {
// 		log.Print(err)
// 		job, err = jenkins.CreateJob(ctx, jenkinsops.GetJobXML(deployinfo["Description"], deployinfo["Appname"], deployinfo["Appid"]), "iac-"+deployinfo["Appname"])
// 		if err != nil {
// 			log.Print(err)

// 			depres.Status = false
// 			depres.Msg = fmt.Sprintf("%v", err)
// 			c <- depres
// 			wg.Done()
// 			return
// 			// return depres
// 		}
// 	}
// 	log.Print(job)

// 	if isupdate {
// 		err = job.UpdateConfig(ctx, jenkinsops.GetJobXML(deployinfo["Description"], deployinfo["Appname"], deployinfo["Appid"]))
// 		if err != nil {
// 			log.Print(err)
// 			depres.Status = false
// 			depres.Msg = fmt.Sprintf("%v", err)
// 			c <- depres
// 			wg.Done()
// 			return
// 		}
// 	}

// 	queid, err := job.InvokeSimple(ctx, deployinfo)
// 	if err != nil {

// 		log.Print(err)

// 		depres.Status = false
// 		depres.Msg = fmt.Sprintf("%v", err)
// 		c <- depres
// 		wg.Done()
// 		return
// 	}

// 	build, err := job.Jenkins.GetBuildFromQueueID(ctx, queid)
// 	if err != nil {
// 		log.Print(err)

//			depres.Status = false
//			depres.Msg = fmt.Sprintf("%v", err)
//			c <- depres
//			wg.Done()
//			return
//		}
//		buildno := build.GetBuildNumber()
//		depres.Status = true
//		depres.Resulturl = jurl + "job/" + "iac-" + deployinfo["Appname"] + "/" + fmt.Sprint(buildno) + "/console"
//		depres.JobName = "iac-" + deployinfo["Appname"]
//		depres.TaskIndex = buildno
//		c <- depres
//		wg.Done()
//	}
func migrateSingle(releaseinfo map[string]string, isupdate, flag bool, c context.Context) Deployresult {
	// var resurl deployresult
	log := logagent.InstPlatform(c)

	depres := Deployresult{Dcenv: releaseinfo["BuildEnv"], Dc: releaseinfo["dc"]}

	ctx := context.Background()
	// restmp := deployresult{status: false, msg: fmt.Sprintf("%v", e)}
	jenkins, jurl, _ := jenkinsops.GetJenkins(releaseinfo["prfileActive"], c)

	jobname := "devops-migrate"
	job, err := jenkins.GetJob(ctx, jobname)
	// bytes, err := json.MarshalIndent(job.Raw.Scm, "", "    ")
	// str, err := job.GetConfig(ctx)
	// if err != nil {
	// 	log.Print(err)
	// }

	// log.Print(str)
	// log.Print(string(bytes))

	// jenkins.Requester.GetXML()

	if err != nil {
		log.Panic("devops-deleteEnv doesnt exist")
	}
	log.Print(job)

	queid, err := job.InvokeSimple(ctx, releaseinfo)
	if err != nil {

		log.Print(err)

		depres.Status = false
		depres.Msg = fmt.Sprintf("%v", err)

		return depres
	}

	build, err := job.Jenkins.GetBuildFromQueueID(ctx, queid) //job.GetLastBuild(ctx) //
	if err != nil {
		log.Print(err)

		depres.Status = false
		depres.Msg = fmt.Sprintf("%v", err)

		return depres
	}
	buildno := build.GetBuildNumber()
	depres.Status = true
	depres.Resulturl = jurl + "job/" + jobname + "/" + fmt.Sprint(buildno) + "/console"
	depres.JobName = jobname
	depres.TaskIndex = buildno

	return depres
}

func releaseSingle(releaseinfo map[string]string, isupdate, flag bool, c context.Context) Deployresult {
	// var resurl deployresult
	log := logagent.InstPlatform(c)

	depres := Deployresult{Dcenv: releaseinfo["BuildEnv"], Dc: releaseinfo["dc"]}

	ctx := context.Background()
	// restmp := deployresult{status: false, msg: fmt.Sprintf("%v", e)}
	jenkins, jurl, _ := jenkinsops.GetJenkins(releaseinfo["prfileActive"], c)

	job, err := jenkins.GetJob(ctx, "devops-deleteEnv")
	// bytes, err := json.MarshalIndent(job.Raw.Scm, "", "    ")
	// str, err := job.GetConfig(ctx)
	// if err != nil {
	// 	log.Print(err)
	// }

	// log.Print(str)
	// log.Print(string(bytes))

	// jenkins.Requester.GetXML()

	if err != nil {
		log.Panic("devops-deleteEnv doesnt exist")
	}
	log.Print(job)

	queid, err := job.InvokeSimple(ctx, releaseinfo)
	if err != nil {

		log.Print(err)

		depres.Status = false
		depres.Msg = fmt.Sprintf("%v", err)

		return depres
	}

	build, err := job.Jenkins.GetBuildFromQueueID(ctx, queid) //job.GetLastBuild(ctx) //
	if err != nil {
		log.Print(err)

		depres.Status = false
		depres.Msg = fmt.Sprintf("%v", err)

		return depres
	}
	buildno := build.GetBuildNumber()
	depres.Status = true
	depres.Resulturl = jurl + "job/devops-deleteEnv/" + fmt.Sprint(buildno) + "/console"
	depres.JobName = "devops-deleteEnv"
	depres.TaskIndex = buildno

	return depres
}

func deploySingle(deployinfo map[string]string, isupdate, flag bool, c context.Context) Deployresult {
	// var resurl deployresult
	log := logagent.InstPlatform(c)

	depres := Deployresult{Dcenv: deployinfo["BuildEnv"], Dc: deployinfo["dc"]}
	if !flag {
		depres.Status = false
		depres.Msg = "deploy in front is error"

		return depres
	}

	ctx := context.Background()
	// restmp := deployresult{status: false, msg: fmt.Sprintf("%v", e)}
	jenkins, jurl, _ := jenkinsops.GetJenkins(deployinfo["prfileActive"], c)

	job, err := jenkins.GetJob(ctx, "iac-"+deployinfo["Appname"])
	// bytes, err := json.MarshalIndent(job.Raw.Scm, "", "    ")
	// str, err := job.GetConfig(ctx)
	// if err != nil {
	// 	log.Print(err)
	// }

	// log.Print(str)
	// log.Print(string(bytes))

	// jenkins.Requester.GetXML()

	if err != nil {
		log.Print(err)
		job, err = jenkins.CreateJob(ctx, jenkinsops.GetJobXML(deployinfo["Description"], deployinfo["Appname"], deployinfo["Appid"]), "iac-"+deployinfo["Appname"])
		if err != nil {
			log.Print(err)

			depres.Status = false
			depres.Msg = fmt.Sprintf("%v", err)

			return depres
			// return depres
		}
	}
	log.Print(job)

	jenkinsUpdated := !jenkinsops.IsSameJob(jenkins, deployinfo["Appname"], deployinfo["Appid"], ctx)
	if jenkinsUpdated {
		deployinfo["isupdate"] = "true"
		err = job.UpdateConfig(ctx, jenkinsops.GetJobXML(deployinfo["Description"], deployinfo["Appname"], deployinfo["Appid"]))
		if err != nil {
			log.Print(err)
			depres.Status = false
			depres.Msg = fmt.Sprintf("%v", err)

			return depres
		}
	}

	queid, err := job.InvokeSimple(ctx, deployinfo)
	if err != nil {

		log.Print(err)

		depres.Status = false
		depres.Msg = fmt.Sprintf("%v", err)

		return depres
	}

	build, err := job.Jenkins.GetBuildFromQueueID(ctx, queid) //job.GetLastBuild(ctx) //
	if err != nil {
		log.Print(err)

		depres.Status = false
		depres.Msg = fmt.Sprintf("%v", err)

		return depres
	}
	buildno := build.GetBuildNumber()
	depres.Status = true
	depres.Resulturl = jurl + "job/" + "iac-" + deployinfo["Appname"] + "/" + fmt.Sprint(buildno) + "/console"
	depres.JobName = "iac-" + deployinfo["Appname"]
	depres.TaskIndex = buildno

	return depres
}

// func deploy(deployinfoseq [][]map[string]string, isupdate bool, c context.Context) []deployresult {
// 	var resurl []deployresult

// 	log := logagent.Inst(c)

// 	flag := true
// 	for _, deployinfos := range deployinfoseq {
// 		var wg sync.WaitGroup
// 		ch := make(chan deployresult, len(deployinfos))
// 		for _, deployinfo := range deployinfos {
// 			// resurl = append(resurl, ) deploySingle(deployinfo, isupdate, flag)

// 			wg.Add(1)
// 			go func(depinfo map[string]string) {
// 				defer func() {
// 					if e := recover(); e != nil {

// 						ch <- deployresult{Dcenv: depinfo["BuildEnv"], Dc: depinfo["dc"], Status: false, Msg: fmt.Sprint(e)}
// 						wg.Done()
// 					}
// 				}()
// 				// deploySingle(deployinfo, isupdate, flag, c, &wg)

// 				ch <- deploySingle(depinfo, isupdate, flag, c)
// 				wg.Done()
// 			}(deployinfo)
// 		}
// 		wg.Wait()
// 		close(ch)
// 		for v := range ch {
// 			resurl = append(resurl, v)
// 			flag = v.Status
// 		}
// 	}
// 	log.Print(resurl)
// 	// STATUS_FAIL           = "FAIL"
// 	// STATUS_ERROR          = "ERROR"
// 	// STATUS_ABORTED        = "ABORTED"
// 	// STATUS_REGRESSION     = "REGRESSION"
// 	// STATUS_SUCCESS        = "SUCCESS"
// 	// STATUS_FIXED          = "FIXED"
// 	// STATUS_PASSED         = "PASSED"
// 	// RESULT_STATUS_FAILURE = "FAILURE"
// 	// RESULT_STATUS_FAILED  = "FAILED"
// 	// RESULT_STATUS_SKIPPED = "SKIPPED"

// 	return resurl
// }

func opsMultiply(deployinfoseq [][]map[string]string, isupdate bool, opsingle opsmethod, c context.Context) []Deployresult {
	var resurl []Deployresult

	log := logagent.InstPlatform(c)

	flag := true
	for _, deployinfos := range deployinfoseq {
		var wg sync.WaitGroup
		ch := make(chan Deployresult, len(deployinfos))
		for _, deployinfo := range deployinfos {
			// resurl = append(resurl, ) deploySingle(deployinfo, isupdate, flag)

			wg.Add(1)
			go func(depinfo map[string]string) {
				defer func() {
					if e := recover(); e != nil {

						ch <- Deployresult{Dcenv: depinfo["BuildEnv"], Dc: depinfo["dc"], Status: false, Msg: fmt.Sprint(e)}
						wg.Done()
					}
				}()
				// deploySingle(deployinfo, isupdate, flag, c, &wg)

				ch <- opsingle(depinfo, isupdate, flag, c)
				wg.Done()
			}(deployinfo)
		}
		wg.Wait()
		close(ch)
		for v := range ch {
			resurl = append(resurl, v)
			flag = v.Status
		}
	}
	log.Print(resurl)
	// STATUS_FAIL           = "FAIL"
	// STATUS_ERROR          = "ERROR"
	// STATUS_ABORTED        = "ABORTED"
	// STATUS_REGRESSION     = "REGRESSION"
	// STATUS_SUCCESS        = "SUCCESS"
	// STATUS_FIXED          = "FIXED"
	// STATUS_PASSED         = "PASSED"
	// RESULT_STATUS_FAILURE = "FAILURE"
	// RESULT_STATUS_FAILED  = "FAILED"
	// RESULT_STATUS_SKIPPED = "SKIPPED"

	return resurl
}

func getDeployconf(appconf archfig.Arch_config, envinfo archfig.EnvInfo, envdc string, c context.Context) (string, string) {
	valconfig := valfig.GenValfig(appconf, envinfo, envdc, c)

	valfile := valfig.GenValfile(valconfig, c)

	// dockerfile := dockfig.GenRuntimeDocfile(appconf, valconfig)
	dockerfile := dockfig.GenDocfile(appconf, c)

	return valfile, dockerfile
}

// func deployold(deployinfoseq [][]map[string]string, isupdate bool) []deployresult {
// 	var resurl []deployresult
// 	flag := true
// 	for _, deployinfos := range deployinfoseq {
// 		for _, deployinfo := range deployinfos {

// 			depres := deployresult{Dcenv: deployinfo["BuildEnv"], Dc: deployinfo["dc"]}
// 			if !flag {
// 				depres.Status = false
// 				depres.Msg = "deploy in front is error"
// 				resurl = append(resurl, depres)
// 				continue
// 			}
// 			// restmp := deployresult{status: false, msg: fmt.Sprintf("%v", e)}
// 			ctx := context.Background()
// 			jenkins, jurl, _ := jenkinsops.GetJenkins(ctx, deployinfo["BuildEnv"])

// 			job, err := jenkins.GetJob(ctx, "iac-"+deployinfo["Appname"])
// 			if err != nil {
// 				log.Print(err)
// 				job, err = jenkins.CreateJob(ctx, jenkinsops.GetJobXML(deployinfo["Description"], deployinfo["Appname"], deployinfo["Appid"]), "iac-"+deployinfo["Appname"])
// 				if err != nil {
// 					log.Print(err)
// 					flag = false
// 					depres.Status = false
// 					depres.Msg = fmt.Sprintf("%v", err)
// 					resurl = append(resurl, depres)
// 					// resurl = append(resurl, deployresult{status: false, msg: fmt.Sprintf("%v", err)})
// 					continue
// 					// resurl = append(resurl, restmp)
// 				}
// 			}
// 			log.Print(job)

// 			if isupdate {
// 				err = job.UpdateConfig(ctx, jenkinsops.GetJobXML(deployinfo["Description"], deployinfo["Appname"], deployinfo["Appid"]))
// 				if err != nil {
// 					log.Print(err)
// 					flag = false
// 					depres.Status = false
// 					depres.Msg = fmt.Sprintf("%v", err)
// 					resurl = append(resurl, depres)
// 					// resurl = append(resurl, deployresult{status: false, msg: fmt.Sprintf("%v", err)})
// 					continue
// 					// resurl = append(resurl, restmp)
// 				}
// 			}
// 			// job.UpdateConfig(ctx, jenkinsops.GetJobXML(deployinfo["Description"], deployinfo["Appname"], deployinfo["Appid"]))

// 			// _, err = jenkins.GetJob(ctx, deployinfo["Appname"])
// 			// if err != nil {
// 			// 	log.Print(err)
// 			// }
// 			queid, err := job.InvokeSimple(ctx, deployinfo)
// 			if err != nil {

// 				log.Print(err)

// 				flag = false
// 				depres.Status = false
// 				depres.Msg = fmt.Sprintf("%v", err)
// 				resurl = append(resurl, depres)
// 				// resurl = append(resurl, deployresult{status: false, msg: fmt.Sprintf("%v", err)})
// 				continue
// 				// log.Panic(err)
// 			}

// 			build, err := job.Jenkins.GetBuildFromQueueID(ctx, queid)
// 			if err != nil {
// 				log.Print(err)

// 				flag = false
// 				depres.Status = false
// 				depres.Msg = fmt.Sprintf("%v", err)
// 				resurl = append(resurl, depres)
// 				// resurl = append(resurl, deployresult{status: false, msg: fmt.Sprintf("%v", err)})
// 				continue
// 				// resurl = append(resurl, restmp)
// 			}
// 			buildno := build.GetBuildNumber()
// 			depres.Status = true
// 			depres.Resulturl = jurl + "job/" + deployinfo["Appname"] + "/" + fmt.Sprint(buildno) + "/console"
// 			depres.JobName = "iac-" + deployinfo["Appname"]
// 			depres.TaskIndex = buildno
// 			resurl = append(resurl, depres)
// 			// resurl = append(resurl,
// 			// 	deployresult{resulturl: jurl + "job/" + deployinfo["Appname"] + "/" + fmt.Sprint(buildno) + "/console",
// 			// 		jobName:   "iac-" + deployinfo["Appname"],
// 			// 		taskIndex: buildno, status: true})
// 			// mm, _ := jenkins.GetBuild(ctx, deployinfo["Appname"], 3)
// 			// url := mm.GetUrl()
// 			// log.Println(mm.GetUrl()) //http://10.47.162.128:83/job/iactest/3/consoleFull
// 			// log.Println(mm.GetResult())
// 		}
// 	}
// 	log.Print(resurl)
// 	// STATUS_FAIL           = "FAIL"
// 	// STATUS_ERROR          = "ERROR"
// 	// STATUS_ABORTED        = "ABORTED"
// 	// STATUS_REGRESSION     = "REGRESSION"
// 	// STATUS_SUCCESS        = "SUCCESS"
// 	// STATUS_FIXED          = "FIXED"
// 	// STATUS_PASSED         = "PASSED"
// 	// RESULT_STATUS_FAILURE = "FAILURE"
// 	// RESULT_STATUS_FAILED  = "FAILED"
// 	// RESULT_STATUS_SKIPPED = "SKIPPED"
// 	return resurl
// }

type checkResult struct {
	Jobname   string
	Taskindex int64
	Result    string
	Dc        string
	Env       string
}

type DeploycheckInfo struct {
	Taskindex int64
	Dc        string
	Env       string
}

// taskindexs []int64, jobname string, envinfo archfig.EnvInfo
func Checkresults(checkinfos []DeploycheckInfo, jobname string, c context.Context) []checkResult {
	var results []checkResult
	var wg sync.WaitGroup
	chain := make(chan checkResult, len(checkinfos))
	for _, checkinfo := range checkinfos {
		wg.Add(1)
		go func(info DeploycheckInfo) {
			defer func() {
				if e := recover(); e != nil {

					chain <- checkResult{
						Jobname:   jobname,
						Taskindex: info.Taskindex,
						Result:    fmt.Sprint(e),
						Dc:        info.Dc,
						Env:       info.Env,
					}
					wg.Done()
				}
			}()
			res := Checkresult(info, jobname, c)
			chain <- res
			wg.Done()
		}(checkinfo)
		// go checksingle(taskindex, jobname, chain, &wg)
	}
	wg.Wait()
	close(chain)
	for v := range chain {
		results = append(results, v)
	}
	return results
}

// func checksingle(taskindex int64, jobname string, c chan deployResult, wg *sync.WaitGroup) deployResult {

// 	res := Checkresult(jobname, taskindex)
// 	c <- res
// 	wg.Done()

//		return res
//	}
//
// jobname string, envinfo archfig.EnvInfo, taskindex int64
func Checkresult(checkinfo DeploycheckInfo, jobname string, c context.Context) checkResult {

	log := logagent.InstPlatform(c)

	envinfo := archfig.EnvInfo{Dc: checkinfo.Dc, Env: checkinfo.Env}
	envstr := getenvstr(envinfo)
	ctx := context.Background()
	jenkins, _, err := jenkinsops.GetJenkins(envstr, c)
	if err != nil {
		log.Panic(err)
	}
	b, err := jenkins.GetBuild(ctx, jobname, checkinfo.Taskindex)

	if err != nil {
		log.Panic(err)
	}
	result := b.GetResult()
	if b.Info().Building {
		result = "Building"
	}

	res := checkResult{
		Jobname:   jobname,
		Taskindex: checkinfo.Taskindex,
		Result:    result,
		Dc:        envinfo.Dc,
		Env:       envinfo.Env,
	}
	return res
}
