package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/logagent/pkg/logsets"
	"github.com/max-gui/logagent/pkg/routerutil"
	"github.com/max-gui/spells/internal/githelp"
	"github.com/max-gui/spells/internal/iac"
	"github.com/max-gui/spells/internal/iac/altconfig"
	"github.com/max-gui/spells/internal/iac/archfig"
	"github.com/max-gui/spells/internal/iac/deploy"
	"github.com/max-gui/spells/internal/pkg/constset"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	// nethttp "net/http"
)

func SetupRouter() *gin.Engine {
	// gin.New()
	// r := gin.Default()
	if *logsets.Appenv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()                      //.Default()
	r.Use(routerutil.GinHeaderMiddle()) // ginHeaderMiddle())
	r.Use(routerutil.GinLogger())       //LoggerWithConfig())
	r.Use(routerutil.GinErrorMiddle())  //ginErrorMiddle())

	// r.Use(ginErrorMiddle())
	// r.POST("/hook/arch/commit/check", Arch_commit_check)
	// r.POST("/hook/archDef/commit/check", ArchDef_commit_check)
	// r.POST("/hook/arch/commit", archCommit)
	// r.POST("/hook/arch/makeup", archMakeup)
	v1 := r.Group("/hook/commit")
	{
		v1.Use(ginCommitMiddle())
		v1.POST("/check", Commit_hook)
		v1.POST("/clear/check", clear_hook)
		v1.POST("/separate/commit/check", separate_commit_hook)
	}
	// r.POST("/hook/archDef/commit", archDefCommit)
	// r.GET("/hook/archDef/makeup", archDefMakeup)
	// r.POST("/gen/values/:appname/:appid/:env/:envdc/:team/:proj", genvalues)
	v2 := r.Group("/apply/deploy")
	{
		v2.POST("/strategy/:appname/:appid/:env/:region/:branch/:team/:proj", applyDeploy)
		v2.POST("/direct/:appname/:appid/:dcenv/:dc/:region/:branch/:team/:proj", targetDeploy)
		v2.POST("/result/task/:env/:dc/:jobname/:taskindex", deployCheck)
		v2.POST("/result/tasks/:jobname", deployChecks)

		// v2.POST("/strategy/list", applyDeploy)
		// v2.POST("/direct/list", targetDeploy)
		// v2.POST("/result/task/list", deployCheck)
		// v2.POST("/result/tasks/list", deployChecks)
	}
	v5 := r.Group("/apply/resource")
	{
		// v5.POST("/free/:relasename", applyResFree)
		v5.POST("/free/strategy/:appname/:appid/:env/:region/:team/:proj", applyResFree)
		v5.POST("/free/direct/:appname/:appid/:dcenv/:dc/:region/:team/:proj", targetResFree)
	}
	v3 := r.Group("/info")
	{
		v3.POST("/proj/:team", getteamproj)
		v3.POST("/app/:proj/:team/:app", getteamprojapp)
		v3.GET("/app/:proj/:team/:app", getteamprojapp)
		v3.GET("/resource/:app", getappres)
	}
	v4 := r.Group("/Reverso")
	{
		v4.GET("/arch/:appid/:appname", reversoArch)
		v4.POST("/resource/key/env", getResourceinfo)
	}
	// v2.Group()
	// r.POST("/apply/deploy/strategy/:appname/:appid/:env/:region/:branch/:team/:proj", applyDeploy)
	// r.POST("/apply/deploy/direct/:appname/:appid/:dcenv/:dc/:region/:branch/:team/:proj", targetDeploy)
	// r.POST("/deploy/result/:jobname/:taskindex", deployCheck)
	r.POST("/apply/clear/cache", clearcache)
	r.POST("/arch/install", arch_install)
	// r.GET("/projinfo/:team", getteamproj)
	// r.GET("/appinfo/:proj/:team/:app", getteamprojapp)
	r.GET("/actuator/health", health)
	r.POST("/hook/mr", MrHook)
	return r
}

func ginCommitMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				logentry := e.(*logrus.Entry)
				c.JSON(http.StatusOK, gin.H{
					"ret": 1,
					"err": logentry.Message,
					"msg": logentry.Message,
				})
			}
		}()

		c.Next()
		// host := c.Request.Host
		// fmt.Printf("Before: %s\n", host)
		// c.Next()
		// fmt.Println("Next: ...")
	}
}

func health(c *gin.Context) {
	// logger := logagent.InstPlatform(c)
	// logger.Info("wonderful!")
	c.String(http.StatusOK, "online")
}

func clearcache(c *gin.Context) {
	iac.ClearcacheAll(c)

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
	})
	// c.JSON(http.StatusOK, nil)
}

func deployChecks(c *gin.Context) {
	// defer func() {
	// 	if e := recover(); e != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{
	// 			"msg": fmt.Sprint(e),
	// 		})
	// 	}
	// }()
	jobname := c.Param("jobname")
	// taskindex := []int64{}
	checkinfos := []deploy.DeploycheckInfo{}
	c.BindJSON(&checkinfos)
	// c.BindJSON(&checkinfos)

	// for _, val := range taskindex {
	// 	checkinfos = append(checkinfos, deploy.DeploycheckInfo{Dc: "LFB", Env: "test", Taskindex: val})
	// }

	log := logagent.InstPlatform(c).WithField("ops-method", "deployChecks").WithField("check-jobname", jobname)
	// deploymap = fn0(appconf, env, region, branch, appname)
	log.Print("start")

	log.Println(checkinfos)
	results := deploy.Checkresults(checkinfos, jobname, c)

	// result := checksingle(c, jobname)
	c.JSON(http.StatusOK, gin.H{
		"result": results,
	})
}

func deployCheck(c *gin.Context) {
	// defer func() {
	// 	if e := recover(); e != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{
	// 			"msg": fmt.Sprint(e),
	// 		})
	// 	}
	// }()

	jobname := c.Param("jobname")
	env := c.Param("env")
	dc := c.Param("dc")
	// envinfo := archfig.EnvInfo{Env: env, Dc: dc}
	// appid := c.Param("appid")
	log := logagent.InstPlatform(c).WithField("ops-method", "deployCheck").WithField("check-jobname", jobname)
	// deploymap = fn0(appconf, env, region, branch, appname)
	log.Print("start")

	taskindex, err := strconv.ParseInt(c.Param("taskindex"), 10, 64)
	checkinfo := deploy.DeploycheckInfo{Taskindex: taskindex, Dc: dc, Env: env}
	if err != nil {
		log.Panic(err)
	}
	result := deploy.Checkresult(checkinfo, jobname, c).Result
	// ctx := context.Background()
	// jenkins, _, _ := jenkinsops.GetJenkins(ctx, "test")
	// b, _ := jenkins.GetBuild(ctx, jobname, taskindex)
	// result := b.GetResult()
	// if b.Info().Building {
	// 	result = "Building"
	// }
	// log.Print(b.GetResult())
	// b.Raw.Building
	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func targetDeploy(c *gin.Context) {
	appname := c.Param("appname")
	// appid := c.Param("appid")
	region := c.Param("region")
	branch := c.Param("branch")
	team := c.Param("team")
	proj := c.Param("proj")
	env := c.Param("dcenv")
	dc := c.Param("dc")

	log := logagent.InstPlatform(c).WithField("ops-method", "targetDeploy").WithField("deploy-app", appname).WithField("deploy-branch", branch).WithField("deploy-region", region)
	// deploymap = fn0(appconf, env, region, branch, appname)
	log.Print("start")
	// resurl := deploy.FlashNdeploy4Target(branch, env, dc, appname, team, proj, region, c) //FlashNdeploy(branch, env, dc, appname, team, proj, region, deploy4target)
	resurl := deploy.TargetFlashDeploy(branch, env, dc, appname, team, proj, region, c)

	log.Print(resurl)
	c.JSON(http.StatusOK, gin.H{
		"result": resurl,
	})
}

func applyDeploy(c *gin.Context) {
	// defer func() {
	// 	if e := recover(); e != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{
	// 			"msg": fmt.Sprint(e),
	// 		})
	// 	}
	// }()

	appname := c.Param("appname")
	// appid := c.Param("appid")
	env := c.Param("env")
	region := c.Param("region")
	branch := c.Param("branch")
	team := c.Param("team")
	proj := c.Param("proj")
	dc := ""

	// var flag bool
	// var deployinfo = []Deployinfo{}
	// var version string
	// var realseName string
	// if BuildEnv == "" && prfileActive == "" && dc == "" {
	// 	//envinfo.Dc,
	// 	// "Projname":			proj,
	// 	deploymap = fn0(appconf, env, region, branch, appname)
	// } else {
	// 	if BuildEnv == "" || prfileActive == "" || dc == "" {
	// 		log.Panic("BuildEnv prfileActive dc should be empty together")
	// 	}
	// 	//envinfo.Dc,
	// 	// "Projname":			proj,
	// 	deploymap = fn1(appconf, BuildEnv, dc, region, branch, appname)
	// }
	log := logagent.InstPlatform(c).WithField("ops-method", "applyDeploy").WithField("deploy-app", appname).WithField("deploy-branch", branch).WithField("deploy-region", region)
	// deploymap = fn0(appconf, env, region, branch, appname)
	log.Print("start")
	// resurl := deploy.FlashNdeploy4Strtegy(branch, env, dc, appname, team, proj, region, c) //renew4deploy(branch, env, dc, appname, team, proj, region, deploy4strategy)
	resurl := deploy.StrtegyFlashDeploy(branch, env, dc, appname, team, proj, region, c)

	log.Print(resurl)
	c.JSON(http.StatusOK, gin.H{
		"result": resurl,
	})
}

func applyResFree(c *gin.Context) {
	// defer func() {
	// 	if e := recover(); e != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{
	// 			"msg": fmt.Sprint(e),
	// 		})
	// 	}
	// }()

	appname := c.Param("appname")
	// appid := c.Param("appid")
	env := c.Param("env")
	region := c.Param("region")
	branch := ""
	team := c.Param("team")
	proj := c.Param("proj")
	dc := ""

	// var flag bool
	// var deployinfo = []Deployinfo{}
	// var version string
	// var realseName string
	// if BuildEnv == "" && prfileActive == "" && dc == "" {
	// 	//envinfo.Dc,
	// 	// "Projname":			proj,
	// 	deploymap = fn0(appconf, env, region, branch, appname)
	// } else {
	// 	if BuildEnv == "" || prfileActive == "" || dc == "" {
	// 		log.Panic("BuildEnv prfileActive dc should be empty together")
	// 	}
	// 	//envinfo.Dc,
	// 	// "Projname":			proj,
	// 	deploymap = fn1(appconf, BuildEnv, dc, region, branch, appname)
	// }
	log := logagent.InstPlatform(c).WithField("ops-method", "applyResFree").WithField("release-app", appname).WithField("region", region)
	// deploymap = fn0(appconf, env, region, branch, appname)
	log.Print("start")
	resurl := deploy.StrtegyFlashRelease(branch, env, dc, appname, team, proj, region, c) //renew4deploy(branch, env, dc, appname, team, proj, region, deploy4strategy)

	log.Print(resurl)
	c.JSON(http.StatusOK, gin.H{
		"result": resurl,
	})
}

func targetResFree(c *gin.Context) {
	appname := c.Param("appname")
	// appid := c.Param("appid")
	region := c.Param("region")
	branch := ""
	team := c.Param("team")
	proj := c.Param("proj")
	env := c.Param("dcenv")
	dc := c.Param("dc")

	log := logagent.InstPlatform(c).WithField("ops-method", "targetResFree").WithField("release-app", appname).WithField("region", region)
	// deploymap = fn0(appconf, env, region, branch, appname)
	log.Print("start")
	resurl := deploy.TargetFlashRelease(branch, env, dc, appname, team, proj, region, c) //FlashNdeploy(branch, env, dc, appname, team, proj, region, deploy4target)

	log.Print(resurl)
	c.JSON(http.StatusOK, gin.H{
		"result": resurl,
	})
}

// func genvalues(c *gin.Context) {
// 	appname := c.Param("appname")
// 	appid := c.Param("appid")
// 	env := c.Param("env")
// 	envdc := c.Param("envdc")
// 	team := c.Param("team")
// 	proj := c.Param("proj")

// 	// bytes := []byte(valfile)
// 	bytes := []byte(getvalues(iac.GetAppconfig(appname, team, proj, appid), env, envdc))

// 	fileContentDisposition := "attachment;filename=\"" + "vaules.yaml" + "\""
// 	c.Header("Content-Type", "application/yml") // 这里是压缩文件类型 .zip
// 	c.Header("Content-Disposition", fileContentDisposition)
// 	c.Data(http.StatusOK, "application/yml", bytes)
// }

func ArchDef_commit_check(c *gin.Context) {

	commitcheck := altconfig.CommitCheckHookInfo{}
	// json := make(map[string]string) //注意该结构接受的内容
	// mm := make(map[string]interface{})
	c.BindJSON(&commitcheck)
	log := logagent.InstPlatform(c)
	log.Println(commitcheck)

	// f0:=func(filename )

	// c.Path
	// if a == "defaultconfig.yaml" {
	// 	iac.GenDefig(iac.Templepath() + a)
	// 	break
	// }
	altconfig.Archdef_commit(commitcheck, false, c)

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
	})
}

func Arch_commit_check(c *gin.Context) {

	commitcheck := altconfig.CommitCheckHookInfo{}

	c.BindJSON(&commitcheck)
	log := logagent.InstPlatform(c)
	log.Println(commitcheck)
	c.JSON(http.StatusOK, gin.H{
		"ret": 1,
		"err": "test error msg",
	})
}

func getteamprojapp(c *gin.Context) {
	team := c.Param("team")
	proj := c.Param("proj")
	app := c.Param("app")

	githelp.UpdateAll(c)
	acrchconf := archfig.GetAppconfig(app, team, proj, c)

	// rediscli := redisops.Pool().Get()

	c.JSON(http.StatusOK, gin.H{
		app: acrchconf,
	})
}

func getResourceinfo(c *gin.Context) {

	keyenv := []struct {
		Key string
		Env string
	}{}
	c.BindJSON(&keyenv)

	keyres := iac.GetKeyres(keyenv, c)

	var res string
	res = "key,env,value\n"
	for _, v := range keyres {
		res += fmt.Sprintf("%s,%s,\"%s\"\n", v["key"], v["env"], v["value"])
	}
	fileContentDisposition := "attachment;filename=\"" + "res.csv" + "\""
	c.Header("Content-Type", "application/csv") // 这里是压缩文件类型 .zip
	c.Header("Content-Disposition", fileContentDisposition)
	c.Data(http.StatusOK, "application/yaml", []byte(res))
}

func reversoArch(c *gin.Context) {

	githelp.UpdateAll(c)
	appid := c.Param("appid")
	appname := c.Param("appname")

	archconf := iac.Gen4old(appid, appname, c)
	bytes, _ := yaml.Marshal(archconf)
	//projapps := iac.GetTeamApps(team)

	// rediscli := redisops.Pool().Get()

	// defer rediscli.Close()

	// m, err := redis.StringMap(rediscli.Do("hgetall", "arch-spell-projteam-"+team))

	// n := make(map[string][]string, len(m))
	// for k, v := range m {
	// 	if n[v] == nil {
	// 		n[v] = []string{}
	// 	}
	// 	n[v] = append(n[v], k)
	// }

	// if err != nil {
	// 	log.Panic(err)
	// }
	fileContentDisposition := "attachment;filename=\"" + "arch.yaml" + "\""
	c.Header("Content-Type", "application/yaml") // 这里是压缩文件类型 .zip
	c.Header("Content-Disposition", fileContentDisposition)
	c.Data(http.StatusOK, "application/yaml", bytes)
}

func getappres(c *gin.Context) {
	app := c.Param("app")

	githelp.UpdateAll(c)
	acrchconf := archfig.GetAppconfigOnline(app, c)

	// rediscli := redisops.Pool().Get()

	c.JSON(http.StatusOK, acrchconf.Application.Resource)
}

func getteamproj(c *gin.Context) {

	team := c.Param("team")

	projapps := iac.GetTeamApps(team, c)

	// rediscli := redisops.Pool().Get()

	// defer rediscli.Close()

	// m, err := redis.StringMap(rediscli.Do("hgetall", "arch-spell-projteam-"+team))

	// n := make(map[string][]string, len(m))
	// for k, v := range m {
	// 	if n[v] == nil {
	// 		n[v] = []string{}
	// 	}
	// 	n[v] = append(n[v], k)
	// }

	// if err != nil {
	// 	log.Panic(err)
	// }
	c.JSON(http.StatusOK, projapps)
}

// func GetProjWithTeam(dirPth string) map[string]map[string]string {
// 	var resmap = make(map[string]map[string]string)
// 	var subdir, archdir []fs.FileInfo
// 	dir, err := ioutil.ReadDir(dirPth)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	rediscli := redisops.Pool().Get()
// 	defer rediscli.Close()

// 	for _, fi := range dir {
// 		if fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") { // 目录, 递归遍历team
// 			// var projtmp []string
// 			resmap[fi.Name()] = make(map[string]string)
// 			subdir, _ = ioutil.ReadDir(dirPth + fi.Name())
// 			for _, subfi := range subdir { //proj
// 				if subfi.IsDir() && !strings.HasPrefix(subfi.Name(), ".") {
// 					archdir, _ = ioutil.ReadDir(dirPth + fi.Name() + string(os.PathSeparator) + subfi.Name())
// 					for _, archfi := range archdir { //app
// 						if !strings.HasPrefix(archfi.Name(), ".") && strings.Split(archfi.Name(), ".")[1] == "yaml" {
// 							resmap[fi.Name()][strings.TrimSuffix(archfi.Name(), ".yaml")] = subfi.Name()
// 							// resmap[subfi.Name()] = fi.Name()
// 						}
// 						// projtmp = append(projtmp, subfi.Name())
// 					}
// 				}
// 			}
// 			if len(resmap[fi.Name()]) > 0 {
// 				rediscli.Do("hmset", redis.Args{}.Add("arch-spell-projteam-"+fi.Name()).AddFlat(resmap[fi.Name()])...)
// 			}
// 		}
// 	}

// 	// rediscli.Do("hmset", redis.Args{}.Add("arch-spell-projteam").AddFlat(resmap)...)

//		return resmap
//	}
func clear_hook(c *gin.Context) {
	commitcheck := altconfig.CommitCheckHookInfo{}
	// json := make(map[string]string) //注意该结构接受的内容
	// mm := make(map[string]interface{})
	c.BindJSON(&commitcheck)
	bs, err := json.Marshal(commitcheck)
	log := logagent.InstPlatform(c).WithField("ops-method", "clear_hook")
	if len(commitcheck.CommitIds) > 0 {
		log = log.WithField("commitid", commitcheck.CommitIds[0])
	} else {
		log = log.WithField("commitid", "0000")
	}

	if err != nil {
		log.Print(err)
	}
	log.Println(string(bs))

	var installflag = false
	if commitcheck.Branch == "master" {
		installflag = true
	}

	switch commitcheck.RepositoryName {
	case "af-db":
		//Sql_commits(commitcheck, c)
	case "templeinfo":
		if installflag {
			iac.ClearcacheAll(c)
		}
		//altconfig.Archdef_commit(commitcheck, installflag, c)
	case "archinfo":
		// if installflag {
		//altconfig.Arch_commits(commitcheck, installflag, c)
		// } else {
		// 	arch_cc(commitcheck)
		// }
	// case "jenkins-library":

	default:
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
	})
}

func separate_commit_hook(c *gin.Context) {
	commitcheck := altconfig.CommitCheckHookInfo{}
	// json := make(map[string]string) //注意该结构接受的内容
	// mm := make(map[string]interface{})
	c.BindJSON(&commitcheck)
	bs, err := json.Marshal(commitcheck)

	log := logagent.InstPlatform(c).WithField("ops-method", "separate_commit_hook")
	if len(commitcheck.CommitIds) > 0 {

		log = log.WithField("commitid", commitcheck.CommitIds[0])
		// deploymap = fn0(appconf, env, region, branch, appname)
		log.Print("start")
	} else {
		log = log.WithField("commitid", "0000")
	}

	if err != nil {
		log.Print(err)
	}

	log.Println(string(bs))

	var installflag = false
	if commitcheck.Branch == "master" {
		installflag = true
	}

	var iacfileinfo []githelp.Writeinfo
	var dbfileinfo []githelp.Writeinfo
	var altinfos []altconfig.Archalt
	repourl := fmt.Sprintf(constset.Codeurl+"/%s/%s.git", commitcheck.Namespace, commitcheck.RepositoryName)
	tmppath := fmt.Sprintf("%s/%s/", commitcheck.Namespace, commitcheck.RepositoryName)

	for _, v := range commitcheck.ChangedFiles {
		if strings.Contains(v.Path, "arch/") {
			changes, archalts := altconfig.Arch_commit(v, repourl, installflag, c)
			iacfileinfo = append(iacfileinfo, changes...)
			altinfos = append(altinfos, archalts...)
		} else if strings.Contains(v.Path, "sql/") {
			change := Sql_commit(v, tmppath, installflag, c)
			dbfileinfo = append(dbfileinfo, change)
		}
	}

	if installflag {
		for _, val := range altinfos {
			val.Arch_info.Install(c)
			val.Org_info.OrgInstall(c)
		}
	}
	altconfig.Tmpt(iacfileinfo, constset.IacUrl, *constset.IacBranch, constset.Iacpath, installflag, c)
	altconfig.Tmpt(dbfileinfo, constset.Dburl, "master", constset.DbPath, installflag, c)

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
	})
}

func arch_install(c *gin.Context) {
	archinfo := archfig.Arch_config{}
	// json := make(map[string]string) //注意该结构接受的内容
	// mm := make(map[string]interface{})
	c.BindJSON(&archinfo)
	bs, err := json.Marshal(archinfo)
	json.Marshal(archfig.Arch_config{})
	log := logagent.InstPlatform(c)
	if err != nil {
		log.Print(err)
	}
	log.Println(string(bs))

	archinfo.Install(c)
	c.String(http.StatusOK, "install success")
}

func Commit_hook(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			c.JSON(http.StatusOK, gin.H{
				"ret": 1,
				"err": fmt.Sprint(e),
				"msg": fmt.Sprint(e),
			})
		}
	}()
	commitcheck := altconfig.CommitCheckHookInfo{}
	// json := make(map[string]string) //注意该结构接受的内容
	// mm := make(map[string]interface{})
	c.BindJSON(&commitcheck)
	bs, _ := json.Marshal(commitcheck)
	log := logagent.InstPlatform(c)
	log.Println(string(bs))

	var installflag = false
	if commitcheck.Branch == "master" {
		installflag = true
	}

	switch commitcheck.RepositoryName {
	case "af-db":
		Sql_commits(commitcheck, c)
	case "templeinfo":
		altconfig.Archdef_commit(commitcheck, installflag, c)
	case "archinfo":
		// if installflag {
		altconfig.Arch_commits(commitcheck, installflag, c)
		// } else {
		// 	arch_cc(commitcheck)
		// }
	// case "jenkins-library":

	default:
	}
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
	})
}
func MrHook(c *gin.Context) {

}

func Sql_commits(commitinfo altconfig.CommitCheckHookInfo, c context.Context) {
	log := logagent.InstPlatform(c)
	for _, c := range commitinfo.ChangedFiles {
		// for _, a := range c..Added {
		if strings.Contains(c.Status, "R") {
			log.Panic("delete and add can't be in one commit,please split and push")
		} else if strings.Contains(c.Path, "sql") {
			if c.Status == "A" || c.Status == "M" {
				lowcontent := strings.ToLower(c.Content)
				if strings.Contains("drop", lowcontent) || strings.Contains("truncate", lowcontent) {
					log.Panic("sql cant contain drop or truncate")
				}
			}
		}
	}
}

func Sql_commit(filecontentinfo archfig.FileContentInfo, repopath string, install bool, c context.Context) githelp.Writeinfo {
	var fileinfo githelp.Writeinfo

	log := logagent.InstPlatform(c)
	if strings.Contains(filecontentinfo.Path, ".sql") {
		if strings.Contains(filecontentinfo.Status, "R") {
			log.Panic("delete and add can't be in one commit,please split and push")
		} else if filecontentinfo.Status == "A" || filecontentinfo.Status == "M" {
			lowcontent := strings.ToLower(filecontentinfo.Content)
			if strings.Contains(lowcontent, "drop") || strings.Contains(lowcontent, "truncate") {
				log.Panic("sql cant contain drop or truncate")
			}

			if install {

				fileinfo = githelp.Writeinfo{Filepath: constset.DbPath + repopath + filecontentinfo.Path, Content: filecontentinfo.Content, Del: false}

			}
		}
	}
	return fileinfo
}
