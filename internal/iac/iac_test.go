package iac

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/max-gui/consulagent/pkg/consulhelp"
	"github.com/max-gui/consulagent/pkg/consulsets"
	"github.com/max-gui/logagent/pkg/confload"
	"github.com/max-gui/logagent/pkg/logsets"
	"github.com/max-gui/redisagent/pkg/redisops"
	"github.com/max-gui/spells/internal/iac/altconfig"
	"github.com/max-gui/spells/internal/iac/archfig"
	"github.com/max-gui/spells/internal/iac/defig"
	"github.com/max-gui/spells/internal/iac/dockfig"
	"github.com/max-gui/spells/internal/iac/jenfig"
	"github.com/max-gui/spells/internal/iac/valfig"
	"github.com/max-gui/spells/internal/pkg/constset"
	"github.com/max-gui/spells/internal/pkg/jenkinsops"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"

	"gopkg.in/yaml.v2"
)

// var Configlist []map[string]interface{}

// var args []string
var orgconfigPth, abstestpath, pthSep string

func setup() {

	// flag.Parse()
	// plaintext = "123"
	// cryptedHexText = "1bda1896724a4521cfb7f38646824197929cd1"
	*logsets.Apppath = "/Users/jimmy/Projects/hercules/spells"
	*logsets.Port = "8080"
	*consulsets.Consul_host = "http://consul-prod.paic.com.cn"

	flag.Parse()
	bytes := confload.Load(context.Background())
	constset.StartupInit(bytes, context.Background())
	// testpath = makeconfiglist()
	pthSep = string(os.PathSeparator)
	orgconfigPth = abstestpath + pthSep + "orgconfig" + pthSep
}

func teardown() {

}

// func makeconfiglist() string { //f0 func(entitytype, entityid, env, configcontent string)) {

// 	// pathname := "yamls"
// 	// pwd, _ := os.Getwd()
// 	pathname := *constset.Apppath + string(os.PathSeparator) + "yamls"
// 	abspath, _ := filepath.Abs(pathname)
// 	// consulhelp.Consulurl = "http://localhost:32771"
// 	// consulhelp.AclToken = ""
// 	files, err := ioutil.ReadDir(abspath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	PthSep := string(os.PathSeparator)
// 	var filename, entitytype, entityid, env, configfilepath string
// 	var entityinfos []string
// 	// var configlist []map[string]interface{}
// 	// var config = make(map[string]interface{})
// 	for _, file := range files {
// 		if path.Ext(file.Name()) == ".yaml" {
// 			filename = strings.Split(file.Name(), ".")[0]
// 			// fmt.Println(filename)
// 			entityinfos = strings.Split(filename, "-")
// 			entitytype = entityinfos[0]
// 			// fmt.Println(entitytype)
// 			entityid = entityinfos[1]
// 			// fmt.Println(entityid)
// 			env = entityinfos[2]
// 			// fmt.Println(env)
// 			if len(entityinfos) > 3 {
// 				env += "-" + entityinfos[3]
// 			}

// 			configfilepath = pathname + PthSep + file.Name()
// 			// fmt.Println(configfilepath)

// 			content, _ := ioutil.ReadFile(configfilepath)

// 			// config := make(map[string]interface{})
// 			// configlist = append(configlist, config)

// 			// fmt.Println(string(content))
// 			_, err := consulhelp.Sendconfig2consul(entitytype, entityid, env, string(content))
// 			if err != nil {
// 				fmt.Println(err.Error())
// 			}
// 			// f0(entitytype, entityid, env, string(content))
// 			// resp, err := consulhelp.Sendconfig2consul(entitytype, entityid, env, string(content))
// 			// if err != nil {
// 			// 	fmt.Println(err.Error())
// 			// }
// 			// fmt.Println(resp)
// 		}
// 	}

// 	return abspath
// }

// func Test_Cases(t *testing.T) {
// 	// <setup code>
// 	setup()

// 	t.Run("Getconfig=Getconfig", Test_Getconfig)
// 	t.Run("GenerateConfig=String", Test_GenerateConfigString)
// 	t.Run("GetPostFileConfig=Encrypt", Test_GetPostFileConfigWithEncrypt)
// 	t.Run("GetPostFileConfig=Decrypt", Test_GetPostFileConfigWithDecrypt)
// 	t.Run("GenerateConfig=ContentList", Test_GenerateConfigContentList)
// 	// <tear-down code>
// 	teardown()
// }

// func prepareTestConfigs() {
// 	Makeconfiglist(func(entitytype, entityid, env, configcontent string) {
// 		resp, err := consulhelp.Sendconfig2consul(entitytype, entityid, env, configcontent)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}
// 		fmt.Println(resp)
// 	})
// }

// func Test_Getconfig(t *testing.T) {

// 	orgstr, _ := fileops.Read(testpath + string(os.PathSeparator) + "consul-consul1-uat.yaml")
// 	orgm := ConvertMap4Json(convertops.ConvertYamlToMap(orgstr), cypher.Decryptbyhex2str)
// 	configm := Getconfig("consul", "consul1", "uat")
// 	// configstr := convertops.ConvertStrMapToYaml(&m)

// 	// assert.NoError(t, err, "read is ok")
// 	assert.Equal(t, orgm, configm)
// 	t.Logf("Test_Getconfig result is:\n%s", configm)

// 	// if !convertops.CompareTwoMapInterface(orgm, configm) {
// 	// 	t.Fatalf("Test_Getconfig failed! config should be:\n%s \nget:\n%s", orgm, configm)
// 	// }
// 	// t.Logf("Test_Getconfig is ok! result is:\n%s", configm)
// }

// func Test_GenerateConfigString(t *testing.T) {
// 	envstr := "sit"
// 	filestr, _ := fileops.Read(testpath + string(os.PathSeparator) + "orgconfig" + string(os.PathSeparator) + "application-test-singalvar.yml")
// 	configstr := GenerateConfigString(filestr, envstr)
// 	configm := ConvertMap4Json(convertops.ConvertYamlToMap(configstr), cypher.Decryptbyhex2str)

// 	var d map[string]interface{}
// 	for k, v := range configm["af-arch"].(map[string]interface{})["resource"].(map[string]interface{}) {
// 		if m, ok := v.(map[string]interface{}); ok {
// 			d = Getconfig(m["entityType"].(string), k, envstr)
// 			// configstr := convertops.ConvertStrMapToYaml(&m)
// 			d["entityType"] = m["entityType"]
// 			assert.Equal(t, v, d)
// 			// if !convertops.CompareTwoMapInterface(d, v.(map[string]interface{})) {
// 			// 	t.Fatalf("Test_Getconfig failed! config should be:\n%s \nget:\n%s", d, v)
// 			// }
// 		}
// 	}

// 	t.Logf("Test_GenerateConfigString result is:\n%s", configm)
// 	// t.Logf("Test_GenerateConfigString is ok! result is:\n%s", configstr)

// }

// func Test_GetPostFileConfigWithEncrypt(t *testing.T) {
// 	f, _ := os.OpenFile(testpath+string(os.PathSeparator)+"orgconfig"+string(os.PathSeparator)+"pg-plain-sit.yaml", os.O_RDONLY, 0644)
// 	orgf, _ := fileops.Read(testpath + string(os.PathSeparator) + "orgconfig" + string(os.PathSeparator) + "pg-pgcypher-sit.yaml")
// 	ortfm := ConvertMap4Json(convertops.ConvertYamlToMap(orgf), func(ciphertext string, key, nonce []byte) string { return "cypher=" + ciphertext })
// 	// ReadFrom(f)
// 	m, _ := GetPostFileConfigWithEncrypt(f)

// 	defer f.Close()
// 	assert.Equal(t, ortfm, m)
// 	t.Logf("Test_GetPostFileConfigWithEncrypt result is:\n%s", m)
// 	// if !convertops.CompareTwoMapInterface(m, ortfm) {
// 	// 	t.Fatalf("GetPostFileConfigWithEncrypt failed! config should be:\n%s \nget:\n%s", m, ortfm)
// 	// }
// 	// t.Logf("GetPostFileConfigWithEncrypt is ok! result is:\n%s", m)

// }

// func Test_GetPostFileConfigWithDecrypt(t *testing.T) {
// 	f, _ := os.OpenFile(testpath+string(os.PathSeparator)+"orgconfig"+string(os.PathSeparator)+"pg-pgcypher-sit.yaml", os.O_RDONLY, 0644)
// 	orgf, _ := fileops.Read(testpath + string(os.PathSeparator) + "orgconfig" + string(os.PathSeparator) + "pg-pg2-sit.yaml")
// 	ortfm := ConvertMap4Json(convertops.ConvertYamlToMap(orgf), func(ciphertext string, key, nonce []byte) string { return ciphertext })
// 	// ReadFrom(f)
// 	m, _ := GetPostFileConfigWithDecrypt(f)

// 	defer f.Close()
// 	assert.Equal(t, ortfm, m)
// 	t.Logf("Test_GetPostFileConfigWithDecrypt result is:\n%s", m)
// 	// if !convertops.CompareTwoMapInterface(m, ortfm) {
// 	// 	t.Fatalf("GetPostFileConfigWithDecrypt failed! config should be:\n%s \nget:\n%s", m, ortfm)
// 	// }
// 	// t.Logf("GetPostFileConfigWithDecrypt is ok! result is:\n%s", m)

// }

// func Test_GenerateConfigContentList(t *testing.T) {
// 	// envstr := "sit"
// 	filestr, _ := fileops.Read(testpath + string(os.PathSeparator) + "orgconfig" + string(os.PathSeparator) + "application-test-singalvar.yml")
// 	var f0 = func(content map[string]interface{}, env string) (map[string]interface{}, error) {
// 		return content, nil
// 	}

// 	m, _ := GenerateConfigContentList(filestr, []string{"sit", "prod"}, f0)

// 	// configstr := GenerateConfigString(filestr, envstr)
// 	// configm := convertMap4Json(convertops.ConvertYamlToMap(configstr), cypher.Decryptbyhex2str)

// 	var d map[string]interface{}
// 	for k, v := range m {
// 		if m, ok := v.(map[string]interface{}); ok {
// 			d = ConvertMap4Json(convertops.ConvertYamlToMap(GenerateConfigString(filestr, k)), cypher.Decryptbyhex2str)

// 			// d = Getconfig(k, m["entityType"].(string), envstr)
// 			// configstr := convertops.ConvertStrMapToYaml(&m)

// 			assert.Equal(t, d, m)
// 			// if !convertops.CompareTwoMapInterface(d, m) {
// 			// 	t.Fatalf("Test_GenerateConfigContentList failed! config should be:\n%s \nget:\n%s", d, m)
// 			// }
// 		}
// 	}

// 	t.Logf("Test_GenerateConfigContentList result is:\n%s", m)
// 	// t.Logf("Test_GenerateConfigContentList is ok! result is:\n%s", m)

// }

func Test_cmdargIgn(t *testing.T) {

	c := context.Background()
	appconf := archfig.GetArchfigSin("af-affe-security-gateway", c)
	appconf.Deploy.Runtime.Args = []string{" aa bb ", " -XX:+UnlockExperimentalVMOptions "}
	appconf.Deploy.Runtime.Ign = map[string][]string{"test": {" aa bb ", " -XX:+UseCGroupMemoryLimitForHeap  "}}

	realarch := archfig.GenArchConfigSinFrominst(appconf, appconf.Application.Name, false, c)
	log.Println(realarch.Deploy.Runtime.Ign["prod"])

	// log.Println(realarch)

	testvf := valfig.GenValfig(realarch, archfig.EnvInfo{Env: "test", Dc: "LFB"}, "test", c)
	testval := valfig.GenValfile(testvf, c)
	log.Println(testval)
	ClearcacheAll(c)
	appconf = archfig.GetArchfigSin("af-affe-security-gateway", c)
	appconf.Deploy.Runtime.Args = []string{" aa bb ", " -XX:+UnlockExperimentalVMOptions "}
	appconf.Deploy.Runtime.Ign = map[string][]string{"test": {" aa bb ", " -XX:+UseCGroupMemoryLimitForHeap  "}}
	realarch = archfig.GenArchConfigSinFrominst(appconf, appconf.Application.Name, false, c)
	log.Println(realarch.Deploy.Runtime.Ign["prod"])
	// log.Println(realarch)
	prodvf := valfig.GenValfig(realarch, archfig.EnvInfo{Env: "prod", Dc: "LFB"}, "prod", c)
	prodval := valfig.GenValfile(prodvf, c)
	log.Println(prodval)
	assert.Contains(t, realarch.Deploy.Runtime.Args, "-javaagent:/wls/wls81/lbagent-1.0.0.jar=http://\\$(POD_IP)/agentcall/$ServiceName/$Buildenv/$dc/")
	log.Println(realarch)
}

func Test_valuesfull(t *testing.T) {

	c := context.Background()
	appconf := archfig.Arch_config{}

	bytes, _ := os.ReadFile("/Users/jimmy/Projects/hercules/iac-tools/charon/arch/af-charon.yaml")
	yaml.Unmarshal(bytes, &appconf)

	// appconf.Application.Name = "test"
	appconf.Application.Type = "java"
	appconf.Application.Service = []string{"abc", "123"}
	appconf.Deploy.Sidecar.Neighbour = []string{"nginx-exporter", "nginx-exporter"}
	// appconf.Application.Resource = map[string]string{"a": "1", "xxl": "2"}
	appconf.Environment.Resource = map[string][]string{"pinpoint": {"detector"}}

	valconfig := valfig.GenValfig(appconf, archfig.EnvInfo{Env: "prod", Dc: "LFB"}, "test", c)

	valfile := valfig.GenValfile(valconfig, c)
	log.Println(valfile)

	appconf.Application.Resource["xxl-viechle"] = "xxl"
	valconfig = valfig.GenValfig(appconf, archfig.EnvInfo{Env: "test", Dc: "LFB"}, "test", c)

	valfile = valfig.GenValfile(valconfig, c)
	log.Println(valfile)

	appconf.Environment.EnHostportable = true
	appconf.Environment.IsHostNetwork = true
	appconf.Environment.Port = ""
	appconf2 := archfig.GenArchConfigSinFrominst(appconf, appconf.Application.Name, false, c)
	valconfig = valfig.GenValfig(appconf2, archfig.EnvInfo{Env: "prod", Dc: "LFB"}, "test", c)

	valfile = valfig.GenValfile(valconfig, c)
	log.Println(valfile)
	// dockerfile := dockfig.GenRuntimeDocfile(appconf, valconfig)
	// dockerfile := dockfig.GenDocfile(appconf, c)
	// log.Print(dockerfile)
}

func Test_setNetworkList(t *testing.T) {

	c := context.Background()
	appconf := archfig.GetArchfigSin("af-affe-security-gateway", c)
	appconf.Environment.Expose.PrefixPath = "sdfwe/dsf"
	appconf.Environment.Expose.Clusternet.Open = true
	appconf.Environment.Expose.Ptrnet.Open = false
	appconf.Environment.Expose.Internet.Visible = true
	appconf.Environment.Expose.Internet.Blacklist = []string{"sdfwe/dsf/bsa", "sdfwe/dsf/dsfa/fed"}
	appconf.Environment.Expose.Intranet.Visible = false
	appconf.Environment.Expose.Intranet.Blacklist = []string{"sdfwe/dsf/a", "sdfwe/dsf/b"}
	// appconf.Environment.Expose.Intranet.Blacklist = []string{"a", "b"}
	// appconf.FireWallRefresh(appconf, c)

	appconf = archfig.GenArchConfigSinFrominst(appconf, "af-affe-security-gateway", false, c)
	appconf.FireWallRefresh4Wthie(c)

	valconf := valfig.GenValfig(appconf, archfig.EnvInfo{Dc: "LFB", Env: "test"}, "test", c)
	values := valfig.GenValfile(valconf, c)
	log.Print(values)

	app2conf := archfig.GetArchfigSin("af-affe-security-gateway", c)
	app2conf.Environment.Expose.Clusternet.Open = false
	app2conf.Environment.Expose.Ptrnet.Open = true
	app2conf.Environment.Expose.Internet.Visible = false
	app2conf.Environment.Expose.Internet.Blacklist = []string{"a", "b"}
	app2conf.Environment.Expose.Intranet.Visible = true
	// app2conf.FireWallRefresh(appconf, c)
	appconf.FireWallRefresh4Wthie(c)

	valconf = valfig.GenValfig(app2conf, archfig.EnvInfo{Dc: "LFB", Env: "test"}, "test", c)
	values = valfig.GenValfile(valconf, c)
	log.Print(values)
}

func GenValfig4envs(appconf archfig.Arch_config, envdcs []archfig.EnvInfo) map[string]valfig.ValuesInfo {
	var res = make(map[string]valfig.ValuesInfo)
	c := context.Background()
	for _, envdc := range envdcs {
		x0 := envdc.Env + envdc.Dc
		res[x0] = valfig.GenValfig(appconf, envdc, x0, c)
	}

	return res
}

func Test_004(t *testing.T) {
	file, err := os.Open(constset.Iacpath + "volumn.csv") //fi.Name() + string(os.PathSeparator) + subfi.Name() + string(os.PathSeparator) + iacfi.Name())
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	c := context.Background()
	scanner := bufio.NewScanner(file)
	mmm := make(map[string]map[string]map[string]string)
	for scanner.Scan() {
		cc := strings.Split(scanner.Text(), ",")
		confenv := make(map[string]map[string]string)
		confm := make(map[string]string)
		if v, ok := mmm[cc[0]]; ok {
			confm["mountPath"] = cc[3]
			confm["hostPath"] = cc[2]
			confm["name"] = cc[0]
			v[cc[1]] = confm
		} else {
			confm["mountPath"] = cc[3]
			confm["hostPath"] = cc[2]
			confm["name"] = cc[0]
			confenv[cc[1]] = confm

			mmm[cc[0]] = confenv
		}
	}

	// mountPath: /wls/wls81/logs
	// hostPath: /nfsc/cnas_csp_stg_fls_aflm_id9192_vol1003_stg/logs/test
	// name: logs
	for _, env := range []string{"prod", "uat", "dr", "test"} {
		dd := make(map[string]string)
		dd["mountPath"] = "/wls/wls81/fake"
		dd["hostPath"] = "/data/fake"
		dd["name"] = "fake.volumn"
		tmp, _ := yaml.Marshal(dd)
		strtmp := string(tmp)
		consulhelp.Sendconfig2consul("volumn", "fakevolumn", env, strtmp, c)
	}

	for k, v := range mmm {
		for _, env := range []string{"prod", "uat", "dr", "test"} {
			if vv, ok := v[env]; ok {
				tmp, _ := yaml.Marshal(vv)
				strtmp := string(tmp)
				consulhelp.Sendconfig2consul("volumn", k, env, strtmp, c)
			} else {
				dd := make(map[string]string)
				for _, vvv := range v {
					dd["mountPath"] = "/wls/wls81/fake"
					dd["hostPath"] = "/data/fake"
					dd["name"] = vvv["name"]
					break
				}

				tmp, _ := yaml.Marshal(dd)
				strtmp := string(tmp)
				consulhelp.Sendconfig2consul("volumn", k, env, strtmp, c)
			}
		}
	}
}

func Test_003(t *testing.T) {
	file, err := os.Open(constset.Iacpath + "hosts.CSV") //fi.Name() + string(os.PathSeparator) + subfi.Name() + string(os.PathSeparator) + iacfi.Name())
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	c := context.Background()
	scanner := bufio.NewScanner(file)
	mmm := make(map[string]map[string]map[string]string)
	mmmttt := make(map[string]string)
	for scanner.Scan() {
		cc := strings.Split(scanner.Text(), ",")
		confenv := make(map[string]map[string]string)
		confm := make(map[string]string)
		if v, ok := mmm[cc[2]]; ok {
			confm["host"] = cc[1]
			confm["ip"] = cc[3]
			v[cc[0]] = confm
		} else {
			confm["host"] = cc[1]
			confm["ip"] = cc[3]
			confenv[cc[0]] = confm

			mmm[cc[2]] = confenv
		}
		mmmttt[cc[0]+cc[1]] = cc[2]
	}

	for _, env := range []string{"prod", "uat", "dr", "test"} {
		dd := make(map[string]string)
		dd["ip"] = "127.0.0.1"
		dd["host"] = "fake.host"
		tmp, _ := yaml.Marshal(dd)
		strtmp := string(tmp)
		consulhelp.Sendconfig2consul("hostAlias", "fakehost", env, strtmp, c)
	}

	for k, v := range mmm {
		for _, env := range []string{"prod", "uat", "dr", "test"} {
			if vv, ok := v[env]; ok {
				tmp, _ := yaml.Marshal(vv)
				strtmp := string(tmp)
				consulhelp.Sendconfig2consul("hostAlias", k, env, strtmp, c)
			} else {
				dd := make(map[string]string)
				dd["real-id"] = "fakehost"
				tmp, _ := yaml.Marshal(dd)
				strtmp := string(tmp)
				consulhelp.Sendconfig2consul("hostAlias", k, env, strtmp, c)
			}
		}
	}
}

func Test_00002(t *testing.T) {
	c := context.Background()

	ddd := archfig.GenArchConfigFrom("team0/project0/iactest2.yaml", "team0", "project0", "iactest2", false, c)
	bytes, _ := yaml.Marshal(ddd)
	log.Print(string(bytes))
	bytes, _ = yaml.Marshal(Gen4old("fls-afch", "afch-authservice", c))
	log.Print(string(bytes))
}

func Test_genvalfile(t *testing.T) {

	c := context.Background()
	appconf := archfig.GetArchfigSin("aflm-cmnsrv-customer-gateway", c)
	valcofn := valfig.GenValfig(appconf, archfig.EnvInfo{Dc: "LFB", Env: "test"}, "test", c)
	str := valfig.GenValfile(valcofn, c)
	log.Print(str)

}

func Test_002(t *testing.T) {
	// dirPth := constset.iacpath + "app" + consts et.PthSep
	// host_config.CSV
	file, err := os.Open(constset.Iacpath + "hosts.CSV") //fi.Name() + string(os.PathSeparator) + subfi.Name() + string(os.PathSeparator) + iacfi.Name())
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	mmm := make(map[string]map[string]map[string]string)
	mmmttt := make(map[string]string)
	for scanner.Scan() {
		cc := strings.Split(scanner.Text(), ",")
		confenv := make(map[string]map[string]string)
		confm := make(map[string]string)
		if v, ok := mmm[cc[2]]; ok {
			confm["host"] = cc[1]
			confm["ip"] = cc[3]
			v[cc[0]] = confm
		} else {
			confm["host"] = cc[1]
			confm["ip"] = cc[3]
			confenv[cc[0]] = confm

			mmm[cc[2]] = confenv
		}
		mmmttt[cc[0]+cc[1]] = cc[2]
	}

	bytess, _ := yaml.Marshal(mmm)
	log.Print(string(bytess))
	bytess, _ = yaml.Marshal(mmmttt)
	log.Print(string(bytess))
	dirPth := constset.Iacpath + "app" + constset.PthSep
	dir, _ := os.ReadDir(dirPth)
	var subdir, iacdir []os.DirEntry
	var nonhost []struct {
		Host string
		Ip   string
		Env  string
		Path string
	}
	var nonvol = []struct {
		MountPath string
		HostPath  string
		Name      string
		Env       string
		Path      string
	}{}

	// var webmap := map[string]string{
	// 	"aflm-apollo-assetfinance-web":  "pazl-web.war",
	// 	"aflm-apollo-pazl-web":  "pazl-web.war",
	// 	"fls-aflm-gdreport-web":  "pazl-web.war",
	// 	"fls-aflm-orderservice-web":  "pazl-web.war",
	// 	"fls-zk-assetfinance-web":  "pazl-web.war",
	// 	"fls-zk-orderservice-web":  "pazl-web.war",
	// 	"fls-zk-report-web":  "pazl-web.war",
	// 	"zk-apollo-web":  "pazl-web.war",
	// }
	// filemap := map[string]struct{}{}

	var nonvalues []string
	var archfigs []archfig.Arch_config
	rediscli := redisops.Pool().Get()

	defer rediscli.Close()

	// bytes, err := redis.String(rediscli.Do("GET"))
	// if err == nil && len(bytes) > 0 {
	// 	rediscli.Do("SETEX", key, 600, bytes)
	// 	return bytes
	// }

	for _, fi := range dir {
		if fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") && !strings.HasPrefix(fi.Name(), "iactest") && fi.Name() != "iac" { // 目录, 递归遍历team
			subdir, _ = os.ReadDir(dirPth + fi.Name())
			for _, subfi := range subdir { //proj

				if str, _ := redis.String(rediscli.Do("GET", dirPth+fi.Name()+string(os.PathSeparator)+subfi.Name())); str != "" {
					log.Println("skip!!!!!!!" + dirPth + fi.Name() + string(os.PathSeparator) + subfi.Name() + string(os.PathSeparator))
					continue
				} else {
					log.Println("in!!!!!!!" + dirPth + fi.Name() + string(os.PathSeparator) + subfi.Name() + string(os.PathSeparator))

				}
				if subfi.IsDir() && !strings.HasPrefix(subfi.Name(), ".") {
					iacdir, _ = os.ReadDir(dirPth + fi.Name() + string(os.PathSeparator) + subfi.Name())
					var AppId, giturl, Cmd, CmdArgs, Pom, Output, AppName, Env, PrePackage, Cpu, Mem, Rtargs string
					var Replica int
					var ExpoviceOk bool
					AppId = fi.Name()

					valuesm := make(map[string]interface{})
					stratgym := make(map[string]struct {
						Capacity string `yaml:"capacity,omitempty"`
						Cpu      string `yaml:"cpu,omitempty"`
						Mem      string `yaml:"mem,omitempty"`
						Replica  int    `yaml:"replica,omitempty"`
					})
					resources := make(map[string]map[string]struct{})
					for _, iacfi := range iacdir { //app

						log.Print("===================================filename===================================================")
						vfilepath := dirPth + fi.Name() + string(os.PathSeparator) + subfi.Name() + string(os.PathSeparator) + iacfi.Name()
						log.Print(dirPth + fi.Name() + string(os.PathSeparator) + subfi.Name() + string(os.PathSeparator) + iacfi.Name())

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
							bytes, err := os.ReadFile(dirPth + fi.Name() + string(os.PathSeparator) + subfi.Name() + string(os.PathSeparator) + iacfi.Name())
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
								nonvalues = append(nonvalues, vfilepath)
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
									ip := vv["ip"].(string)
									hostnames := vv["hostnames"].([]interface{})
									for _, hs := range hostnames {
										if hostval, hostvalok := mmmttt[Env+hs.(string)]; hostvalok {
											// resid := mmmttt[hs+ip]
											if resources["hostAlias"] == nil {
												resources["hostAlias"] = map[string]struct{}{}
											}
											resources["hostAlias"][hostval] = struct{}{} //append(resources["hostAlias"], hostval)
										} else {
											nonhost = append(nonhost, struct {
												Host string
												Ip   string
												Env  string
												Path string
											}{
												Host: hs.(string),
												Ip:   ip,
												Env:  Env,
												Path: vfilepath,
											})
										}
										// nonhost = append(nonhost, nonhost hs.(string)+"--"+ip)
										// id,host,ip,env

									}
								}
							}
							// volumeMounts:
							// - mountPath: /wls/wls81/logs
							//   name: logs
							// volumes:
							// - hostPath:
							//     path: /nfsc/cnas_csp_stg_fls_aflm_id9192_vol1003_stg/logs/test
							//   name: logs
							if val, ok := valuesm["volumeMounts"]; ok {
								for _, v := range val.([]interface{}) {
									vv := v.(map[interface{}]interface{})
									volname := vv["name"].(string)
									mountPath := vv["mountPath"].(string)
									var hostPath string
									for _, vvv := range valuesm["volumes"].([]interface{}) {
										vvvv := vvv.(map[interface{}]interface{})
										if vvvv["name"].(string) == volname {
											// log.Print(vvvv)
											log.Print(vvvv)

											vd := vvvv["hostPath"].(map[interface{}]interface{})
											hostPath = vd["path"].(string)
											break
										}
									}
									// volumes:[map[hostPath:map[path:/data/logs] name:logs] map[hostPath:map[path:/data/classpath] name:classpath]]]
									nonvol = append(nonvol, struct {
										MountPath string
										HostPath  string
										Name      string
										Env       string
										Path      string
									}{
										Name:      volname,
										MountPath: mountPath,
										HostPath:  hostPath,
										Env:       Env,
										Path:      vfilepath,
									})
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
							bytes, err := os.ReadFile(dirPth + fi.Name() + string(os.PathSeparator) + subfi.Name() + string(os.PathSeparator) + iacfi.Name())
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
							file, err := os.Open(dirPth + fi.Name() + string(os.PathSeparator) + subfi.Name() + string(os.PathSeparator) + iacfi.Name())
							if err != nil {
								log.Panic(err)
							}
							defer file.Close()

							scanner = bufio.NewScanner(file)
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
					archfig := archfig.Arch_config{}
					archfig.Application.Appid = AppId
					archfig.Application.Name = AppName
					archfig.Application.Project = ""
					archfig.Application.Repositry = giturl
					archfig.Application.Team = ""
					if PrePackage == "" {
						archfig.Application.Type = "java"
					} else {
						archfig.Application.Type = "h5"

						if ExpoviceOk {
							archfig.Environment.Expose.Expovice = PrePackage
							archfig.Application.Name = PrePackage
						}
					}
					for rk, rv := range resources {
						var rvks = []string{}
						for rvk := range rv {
							rvks = append(rvks, rvk)
						}
						if archfig.Environment.Resource == nil {
							archfig.Environment.Resource = map[string][]string{}
						}
						archfig.Environment.Resource[rk] = rvks
					}
					// archfig.Environment.Resource = resources
					archfig.Environment.Strategy = stratgym
					archfig.Deploy.Build.Args = strings.Trim(CmdArgs, "'")
					archfig.Deploy.Build.Cmd = strings.Trim(Cmd, "'")
					archfig.Deploy.Build.Output = strings.Trim(Output, "'")
					archfig.Deploy.Build.Pkgconf = Pom
					archfig.Deploy.Runtime.Args = strings.Split(Rtargs, " ")
					archfig.Application.Ungenfig = true

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
					rediscli.Do("SET", dirPth+fi.Name()+string(os.PathSeparator)+subfi.Name(), "skip")
				}
			}
			// if len(resmap[fi.Name()]) > 0 {
			// 	rediscli.Do("hmset", redis.Args{}.Add("arch-spell-projteam-"+fi.Name()).AddFlat(resmap[fi.Name()])...)
			// }
		}
	}
	// nonhostm := make(map[string]struct {
	// 	Host string
	// 	Ip   string
	// 	Env  string
	// 	Path string
	// })
	// for _, v := range nonhost {
	// 	nonhostm[v.Host+v.Ip] = v

	// }
	for _, v := range nonhost {
		log.Println(v.Env + "," + v.Host + "," + v.Ip + "," + strings.TrimPrefix(v.Path, dirPth))
	}
	for _, v := range archfigs {
		log.Println(v.Application.Name + "," + v.Application.Repositry + ",")
	}
	for _, v := range nonvol {
		log.Println(v.Name + "," + v.Env + "," + v.HostPath + "," + v.MountPath + "," + strings.TrimPrefix(v.Path, dirPth))
	}

	// // bbbb, _ := json.Marshal(nonhost)
	// // log.Print(string(bbbb))
	bbbb, _ := json.Marshal(nonvalues)
	log.Print(string(bbbb))
}

func Test_put(t *testing.T) {

	c := context.Background()
	mm := consulhelp.Getconfibytes(*constset.ConfResPrefix, "ops", "s", "d", c)
	log.Print(mm)
}

func Test_startWatch(t *testing.T) {
	go startWatch()
	// util.CheckError(err)
	// util.CheckError(err)
	// Get a new client
	// conf.Address
	// startWatch()
	done := make(chan int)
	<-done
}

func startWatch() {
	watchConfig := make(map[string]interface{})
	var first = true
	watchConfig["type"] = "keyprefix"
	watchConfig["prefix"] = "ops/"
	// watchConfig["handler_type"] = "script"
	watchPlan, err := watch.Parse(watchConfig)
	watchPlan.Token = *consulsets.Acltoken
	if err != nil {
		log.Panic(err)
	}

	var kvmap = make(map[string][]byte)
	watchPlan.Handler = func(lastIndex uint64, result interface{}) {
		keys := result.(api.KVPairs)

		for _, v := range keys {
			if first || v.ModifyIndex == lastIndex {
				log.Print(string(v.Value))
				kvmap[v.Key] = v.Value
			}
		}
		first = false

	}

	conf := api.DefaultConfig()
	conf.Address = *consulsets.Consul_host

	conf.Token = *consulsets.Acltoken

	err = watchPlan.Run(*consulsets.Consul_host)
	if err != nil {
		log.Fatalf("start watch error, error message: %s", err.Error())
	}
}

func Test_0000(t *testing.T) {
	ctx := context.Background()
	jenkins, _, _ := jenkinsops.GetJenkins("test", ctx)
	// jenkins.GetJob(ctx, "iac-fls-usedcar-trans-ms")
	//m :=map[interface{}]interface{}{}"pub_zlqrdevops:Devops#0419""

	flag := jenkinsops.IsSameJob(jenkins, "fls-usedcar-trans-ms", "fls-cgj", ctx)
	log.Print(flag)
	flag = jenkinsops.IsSameJob(jenkins, "fls-usedcar-trans-ms", "fls-ccgj", ctx)
	log.Print(flag)

	job, _ := jenkins.GetJob(ctx, "iac-fls-usedcar-trans-ms")

	log.Print(job.Raw.Scm)

	b, _ := jenkins.GetBuild(ctx, "fls-aflm-risk-management", 68)
	log.Print(b.GetResult())

	b, _ = jenkins.GetBuild(ctx, "fls-cgj-common-coupon.bak", 220)
	log.Print(b.GetResult())
	//find text
	// log.Print(path)
	// log.Print(filepath.Dir(d.Name()))
	name := map[string]struct{}{"secgwiac.yaml": {}, "iactest.yaml": {}, "bb": {}}
	fullpath := altconfig.Findarchfile(name, ctx)
	delete(fullpath, "teamarch/projsecgw/secgwiac.yaml")
	delete(fullpath, "team0/project0/iactest.yaml")
	delete(fullpath, "bb")
	if len(fullpath) > 0 {
		log.Panicf("these apps need to be removed:%v", fullpath)
	}
}

func Test_f00(t *testing.T) {
	// pool := redisops.NewPool(":6379")

	rediscli := redisops.Pool().Get()
	c := context.Background()

	defer rediscli.Close()

	defig.GenDefig(true, c)
	// jsonbs, _ := json.Marshal(defig)
	// rediscli.Do("SET", "arch-spell-default", jsonbs)
	// Defconfig = defconf
	// var appname = "test"
	var team, proj = "teamarch", "projsecgw"

	var archconfig = archfig.GenArchConfigFrom(team+constset.PthSep+proj+constset.PthSep+"electronicsign.yaml", team, proj, "electronicsign", true, c) //constset.Apppath+constset.PthSep+"yamls"+constset.PthSep+"orgconfig"+pthSep+"arch.yaml", true)
	bytes, _ := yaml.Marshal(archconfig)
	log.Println(string(bytes))

	// jsonbs, _ = json.Marshal(archfig)
	// rediscli.Do("HSET", "arch-spell-appconfig", "demo", jsonbs)
	// mm, _ := redis.Bytes(rediscli.Do("HGET", "arch-spell-appconfig", "demo"))
	// mm, err := redis.Bytes(rediscli.Do("HGET", "arch-spell-appconfig", "dd"))
	// if err != nil {
	// 	log.Print(err)
	// }
	// log.Print(string(mm))
	dockfig.MakeDockemple(archconfig.Application.Type, true, c)
	str := dockfig.GenDocfile(archconfig, c)
	// str = DockerGen(archfig, true)
	log.Print(str)

	var jenconfig = jenfig.GenJenfig(archconfig)
	bytes, _ = yaml.Marshal(jenconfig)
	log.Println(string(bytes))

	jenfig.MakeJenkimple(archconfig.Application.Type, true, c)
	// jenfig.Output = "aaa"
	str = jenfig.GenJenfile(jenconfig, c)
	// name := "jenkins." + archfig.Application.Type
	// templtmp := templ.GetemplFrom(Templepath()+name, name)
	// str = templ.Gempl(archfig, jenfig)
	log.Print(str)

	// ctx := context.Background()
	// jenkins, _, _ := jenkinsops.GetJenkins(ctx, "test")
	// job, err := jenkins.CreateJob(ctx, jenkinsops.GetJobXML(archfig.Application.Description, archfig.Application.Name, archfig.Application.Appid), archfig.Application.Name)
	// if err != nil {
	// 	log.Print(err)
	// }

	// job, err = jenkins.GetJob(ctx, archfig.Application.Name)
	// _, err = jenkins.BuildJob(ctx, archfig.Application.Name, nil)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Print(job)

	var envinfos []archfig.EnvInfo
	// var flowtmp []string
	var ff [][]string
	// var envstr string
	log.Print(archconfig.Deploy.Stratail)

	for _, v := range archconfig.Deploy.Stratail {
		for _, deployvs := range v {
			for _, deployv := range deployvs {
				// envstr := deployv.Env + "-" + deployv.Dc
				// envstr := deployv.Env + deployv.Dc
				// envstr = strings.TrimSuffix(envstr, "LFB")
				envinfos = append(envinfos, archfig.EnvInfo{Env: deployv.Env, Dc: deployv.Dc})
			}
		}
	}

	log.Print(envinfos)
	log.Print(ff)
	// valfig.MakeValuemple(archfig.Application.Type, true)
	// templtmp = templ.GetemplFrom(Templepath()+"values.yaml", "values")
	var valfigs = GenValfig4envs(archconfig, envinfos) //[]string{"test", "prod"})
	for k, v := range valfigs {
		// bytes, _ = yaml.Marshal(v)
		// log.Printf("\nresult for env:%s\n%s", string(k), string(bytes))
		result := valfig.GenValfile(v, c) //templ.Gempl(templtmp, v)

		log.Printf("\nresult values for env:%s\n%s", string(k), result)
	}

	valconfig := valfig.GenValfig(archconfig, envinfos[0], envinfos[0].Env+envinfos[0].Dc, c)
	result := valfig.GenValfile(valconfig, c) //templ.Gempl(templtmp, v)

	log.Printf("\nresult values for env:%s\n%s", envinfos[0], result)
}

func Test_dd(t *testing.T) {
	// app_conf := archfig.Arch_config{}
	// team := "team"
	// proj := "proj"
	// appfname := "app"
	// app_conf.Application.Name = appfname
	// app_conf.Application.Appid = "SS"
	// app_conf.Application.Repositry = "ssh"
	// app_conf.Application.Type = "h5"
	// app_conf.Environment.Strategy = map[string]struct {
	// 	Capacity string
	// 	Cpu      string
	// 	Mem      string
	// 	Replica  int
	// }{"test": {Capacity: "mini"}}
	// app_conf.Deploy.Build.Output = "dd"
	// bytes, _ := yaml.Marshal(app_conf)

	// bytes, _ := ioutil.ReadFile("/Users/max/Downloads/spells/pv-core-secgw.yml")
	c := context.Background()

	// str, _ := ioutil.ReadFile("/Users/max/Downloads/dd.yaml")
	// archfig.GenArchConfig(str, "front-end", "fls-cgj-vue-pdf", "fls-cgj-vue-pdf", false, c)
	app_conf := archfig.GetAppconfig("fls-aflm-fund-loan", "", "", c)
	app_conf.Deploy.Build.Jenkignor = append(app_conf.Deploy.Build.Jenkignor, "ee")
	app_conf.Deploy.Build.Jenkexec = append(app_conf.Deploy.Build.Jenkexec, "bb")

	app_conf = archfig.GenArchConfigSinFrominst(app_conf, "fls-aflm-fund-loan", false, c)
	// app_conf := archfig.GenArchConfigSin(bytes, "pv-core-secgw", false)
	// app_conf := archfig.GenArchConfig(bytes, "arch", "hercules", "arch-spells", false)
	log.Print(app_conf.Deploy.Build.Jenkignor)
	// app_conf.Environment.NodeSelector = map[string]string{"test": "test"}
	// log.Print(app_conf.Deploy.Runtime.Args)
	valconf := valfig.GenValfig(app_conf, archfig.EnvInfo{Env: "test", Dc: "LFB"}, "test", c)

	bytes, _ := yaml.Marshal(valconf)
	log.Print(string(bytes))
	log.Print(valfig.GenValfile(valconf, c))
	// log.Print(dockfig.GenDocfile(app_conf))
	// log.Print(dockfig.GenRuntimeDocfile(app_conf, valconf))
	// log.Print(jenfig.GenJenfile(jenfig.GenJenfig(app_conf)))
}

func Test_unsecapps(t *testing.T) {
	// app_conf := archfig.Arch_config{}
	// team := "team"
	// proj := "proj"
	// appfname := "app"
	// app_conf.Application.Name = appfname
	// app_conf.Application.Appid = "SS"
	// app_conf.Application.Repositry = "ssh"
	// app_conf.Application.Type = "h5"
	// app_conf.Environment.Strategy = map[string]struct {
	// 	Capacity string
	// 	Cpu      string
	// 	Mem      string
	// 	Replica  int
	// }{"test": {Capacity: "mini"}}
	// app_conf.Deploy.Build.Output = "dd"
	// bytes, _ := yaml.Marshal(app_conf)

	// bytes, _ := ioutil.ReadFile("/Users/max/Downloads/spells/pv-core-secgw.yml")
	kvs := consulhelp.GetConfigs("/ops/iac/arch", "", context.Background())

	var archconf archfig.Arch_config

	teammap := map[string]string{"front-end": "大前端", "bmp": "大中台", "pv-core": "乘用车", "cgj": "车服务", "arch": "架构", "autoheavytruck": "商用车", "autoqa": "测试", "data": "大数据", "css": "pmo", "pmo": "pmo"}
	envarr := map[string]string{"prod": "512Mi", "uat": "256Mi", "test": "256Mi"}
	// log.Print("团队,项目,应用,环境,内存,建议内存,代码库地址")

	// f := func(path string, processor func(*bufio.Writer)) {
	// 	os.Remove(path)
	// 	f, _ := os.Create(path)
	// 	defer f.Close()

	// 	buff := bufio.NewWriter(f)
	// 	processor(buff)
	// 	buff.Flush()
	// }

	writeExecel := func(path string, processor func(*excelize.File)) {
		os.Remove(path)

		excel := excelize.NewFile()
		defer excel.Close()

		processor(excel)
		excel.SaveAs(path)
	}

	// excel := excelize.NewFile()
	// excel.NewSheet("a")
	// excel.SetSheetRow("Sheet1", "a1", &[]interface{}{"1", "a", 2})
	// excel.SetSheetRow("a", "a2", &[]interface{}{"2", "b", 3})
	// excel.SaveAs("/Users/max/Projects/test.xlsx")
	writeExecel("/Users/jimmy/Projects/资源不合理.xlsx", func(excel *excelize.File) {
		excel.SetSheetRow("Sheet1", "a1", &[]interface{}{"团队", "项目", "应用", "环境", "内存", "建议内存", "代码库地址"})
		var row = 2
		rowstr := "a2"
		for _, v := range kvs {
			json.Unmarshal(v.Value, &archconf)
			if archconf.Application.Type == "h5" || archconf.Application.Type == "apollofront" {
				for env, memsuggi := range envarr {
					realmem, _ := strconv.Atoi(strings.TrimSuffix(archconf.Environment.Strategy[env].Mem, "Mi"))
					suggimem, _ := strconv.Atoi(strings.TrimSuffix(memsuggi, "Mi"))
					if realmem > suggimem {
						rowstr = "a" + strconv.Itoa(row)
						row++

						excel.SetSheetRow("Sheet1", rowstr, &[]interface{}{
							teammap[archconf.Application.Team],
							archconf.Application.Project,
							archconf.Application.Name,
							env,
							archconf.Environment.Strategy[env].Mem,
							memsuggi,
							archconf.Application.Repositry})
						// 	  "建议内存", "代码库地址"})
						// msg := teammap[archconf.Application.Team] + "," + archconf.Application.Project + "," + archconf.Application.Name + "," + env + "," + archconf.Environment.Strategy[env].Mem + "," + memsuggi + "," + archconf.Application.Repositry
						// buff.WriteString(msg + "\n")
					}
				}
			}
		}
	})

	// f("/Users/max/Projects/资源不合理.csv", func(buff *bufio.Writer) {
	// 	buff.WriteString("团队,项目,应用,环境,内存,建议内存,代码库地址\n")

	// 	for _, v := range kvs {
	// 		json.Unmarshal(v.Value, &archconf)
	// 		if archconf.Application.Type == "h5" || archconf.Application.Type == "apollofront" {
	// 			for env, memsuggi := range envarr {
	// 				realmem, _ := strconv.Atoi(strings.TrimSuffix(archconf.Environment.Strategy[env].Mem, "Mi"))
	// 				suggimem, _ := strconv.Atoi(strings.TrimSuffix(memsuggi, "Mi"))
	// 				if realmem > suggimem {
	// 					msg := teammap[archconf.Application.Team] + "," + archconf.Application.Project + "," + archconf.Application.Name + "," + env + "," + archconf.Environment.Strategy[env].Mem + "," + memsuggi + "," + archconf.Application.Repositry
	// 					buff.WriteString(msg + "\n")
	// 				}
	// 			}
	// 		}
	// 	}
	// })

	writeExecel("/Users/jimmy/Projects/对外暴露不合规.xlsx", func(excel *excelize.File) {
		excel.SetSheetRow("Sheet1", "a1", &[]interface{}{"团队", "项目", "应用", "expovice", "违规形式", "应用类型", "代码库地址"})
		var row = 2
		rowstr := "a2"
		for _, v := range kvs {
			json.Unmarshal(v.Value, &archconf)
			if archconf.Environment.Expose.Internet.Visible == false && archconf.Environment.Expose.Intranet.Visible == false && archconf.Environment.Expose.Expovice != "" {
				rowstr = "a" + strconv.Itoa(row)
				row++
				excel.SetSheetRow("Sheet1", rowstr, &[]interface{}{
					teammap[archconf.Application.Team],
					archconf.Application.Project,
					archconf.Application.Name,
					archconf.Environment.Expose.Expovice,
					"none area",
					archconf.Application.Type,
					archconf.Application.Repositry})
				// buff.WriteString(teammap[archconf.Application.Team] + "," + archconf.Application.Project + "," + archconf.Application.Name + "," + archconf.Environment.Expose.Expovice + "," + "none area" + "," + archconf.Application.Type + "," + archconf.Application.Repositry + "\n")
			}
			if archconf.Environment.Expose.Expovice != "" &&
				archconf.Application.Type != "h5" &&
				archconf.Application.Type != "fls-aflm-orderservice-web" &&
				archconf.Application.Type != "apollofront" {

				if !strings.Contains(archconf.Environment.Expose.Expovice, "security-gateway") {
					rowstr = "a" + strconv.Itoa(row)
					row++
					excel.SetSheetRow("Sheet1", rowstr, &[]interface{}{
						teammap[archconf.Application.Team],
						archconf.Application.Project,
						archconf.Application.Name,
						archconf.Environment.Expose.Expovice,
						"app exposed",
						archconf.Application.Type,
						archconf.Application.Repositry})
					// buff.WriteString(teammap[archconf.Application.Team] + "," + archconf.Application.Project + "," + archconf.Application.Name + "," + archconf.Environment.Expose.Expovice + "," + "app exposed" + "," + archconf.Application.Type + "," + archconf.Application.Repositry + "\n")
				} else if strings.Contains(archconf.Application.Name, "security-gateway") {
					rowstr = "a" + strconv.Itoa(row)
					row++
					excel.SetSheetRow("Sheet1", rowstr, &[]interface{}{
						teammap[archconf.Application.Team],
						archconf.Application.Project,
						archconf.Application.Name,
						archconf.Environment.Expose.Expovice,
						"secgw exposed",
						archconf.Application.Type,
						archconf.Application.Repositry})
					// buff.WriteString(teammap[archconf.Application.Team] + "," + archconf.Application.Project + "," + archconf.Application.Name + "," + archconf.Environment.Expose.Expovice + "," + "sec exposed" + "," + archconf.Application.Type + "," + archconf.Application.Repositry + "\n")
				} else {
					rowstr = "a" + strconv.Itoa(row)
					row++
					excel.SetSheetRow("Sheet1", rowstr, &[]interface{}{
						teammap[archconf.Application.Team],
						archconf.Application.Project,
						archconf.Application.Name,
						archconf.Environment.Expose.Expovice,
						"securited",
						archconf.Application.Type,
						archconf.Application.Repositry})
					// buff.WriteString(teammap[archconf.Application.Team] + "," + archconf.Application.Project + "," + archconf.Application.Name + "," + archconf.Environment.Expose.Expovice + "," + "securited" + "," + archconf.Application.Type + "," + archconf.Application.Repositry + "\n")
				}
			}
		}
	})

	writeExecel("/Users/jimmy/Projects/服务依赖完成情况.xlsx", func(excel *excelize.File) {
		excel.SetSheetRow("Sheet1", "a1", &[]interface{}{"团队", "项目", "应用", "违规形式", "应用类型", "代码库地址"})
		var row = 2
		rowstr := "a2"
		for _, v := range kvs {

			json.Unmarshal(v.Value, &archconf)
			if archconf.Application.Service == nil {
				rowstr = "a" + strconv.Itoa(row)
				row++
				excel.SetSheetRow("Sheet1", rowstr, &[]interface{}{
					teammap[archconf.Application.Team],
					archconf.Application.Project,
					archconf.Application.Name,
					"unfinished",
					archconf.Application.Type,
					archconf.Application.Repositry})
				// buff.WriteString(teammap[archconf.Application.Team] + "," + archconf.Application.Project + "," + archconf.Application.Name + "," + "unfinished" + "," + archconf.Application.Type + "," + archconf.Application.Repositry + "\n")
			} else if len(archconf.Application.Service) == 0 {
				rowstr = "a" + strconv.Itoa(row)
				row++
				excel.SetSheetRow("Sheet1", rowstr, &[]interface{}{
					teammap[archconf.Application.Team],
					archconf.Application.Project,
					archconf.Application.Name,
					"empty",
					archconf.Application.Type,
					archconf.Application.Repositry})
				// buff.WriteString(teammap[archconf.Application.Team] + "," + archconf.Application.Project + "," + archconf.Application.Name + "," + "empty" + "," + archconf.Application.Type + "," + archconf.Application.Repositry + "\n")
			} else {
				rowstr = "a" + strconv.Itoa(row)
				row++
				excel.SetSheetRow("Sheet1", rowstr, &[]interface{}{
					teammap[archconf.Application.Team],
					archconf.Application.Project,
					archconf.Application.Name,
					"finished",
					archconf.Application.Type,
					archconf.Application.Repositry})
				// buff.WriteString(teammap[archconf.Application.Team] + "," + archconf.Application.Project + "," + archconf.Application.Name + "," + "finished" + "," + archconf.Application.Type + "," + archconf.Application.Repositry + "\n")
			}
		}
	})
	writeExecel("/Users/jimmy/Projects/应用信息.xlsx", func(excel *excelize.File) {
		excel.SetSheetRow("Sheet1", "a1", &[]interface{}{"systemName", "moduleName", "gitUrl", "codeLanguage", "repositoryType", "scanCodePath"})
		var row = 2
		rowstr := "a2"
		for _, v := range kvs {

			json.Unmarshal(v.Value, &archconf)
			if archconf.Application.Type != "go" && archconf.Application.Type != "python" {
				apptype := archconf.Application.Type
				switch apptype {
				case "apollofront":
					apptype = "h5"
				case "apollobackend":
					apptype = "java"
				case "tomcat-common":
					apptype = "java"
				case "tomcat":
					apptype = "java"
				case "tomcat-job":
					apptype = "java"
				case "tomcat-wx-admin":
					apptype = "java"
				case "java-heracles":
					apptype = "java"
				case "nodejs":
					apptype = "h5"
				}
				rowstr = "a" + strconv.Itoa(row)
				row++
				excel.SetSheetRow("Sheet1", rowstr, &[]interface{}{
					archconf.Application.Appid,
					archconf.Application.Name,
					archconf.Application.Repositry,
					apptype,
					"GIT",
					"/"})
			}
		}
	})

	writeExecel("/Users/jimmy/Projects/应用资源.xlsx", func(excel *excelize.File) {
		excel.SetSheetRow("Sheet1", "a1", &[]interface{}{"团队", "项目", "应用", "环境", "内存", "cpu核心数(弹性)", "代码库地址"})
		var row = 2
		rowstr := "a2"
		for _, v := range kvs {
			json.Unmarshal(v.Value, &archconf)
			// if archconf.Application.Type == "h5" || archconf.Application.Type == "apollofront" {
			for env := range envarr {
				// realmem, _ := strconv.Atoi(strings.TrimSuffix(archconf.Environment.Strategy[env].Mem, "Mi"))
				// suggimem, _ := strconv.Atoi(strings.TrimSuffix(memsuggi, "Mi"))
				// if realmem > suggimem {
				rowstr = "a" + strconv.Itoa(row)
				row++

				excel.SetSheetRow("Sheet1", rowstr, &[]interface{}{
					teammap[archconf.Application.Team],
					archconf.Application.Project,
					archconf.Application.Name,
					env,
					archconf.Environment.Strategy[env].Mem,
					archconf.Environment.Strategy[env].Cpu,
					archconf.Application.Repositry})
				// 	  "建议内存", "代码库地址"})
				// msg := teammap[archconf.Application.Team] + "," + archconf.Application.Project + "," + archconf.Application.Name + "," + env + "," + archconf.Environment.Strategy[env].Mem + "," + memsuggi + "," + archconf.Application.Repositry
				// buff.WriteString(msg + "\n")
				// }
			}
			// }
		}
	})

	// log.Print(kvs)
}

func Test_fasfds(t *testing.T) {
	mm := make(map[string][]string)

	nn := make(map[string][]string)
	nn["dsf"] = append(nn["dsf"], mm["sdf"]...)
	log.Print(nn)

}

func TestMain(m *testing.M) {

	// flag.Parse()

	setup()
	// constset.StartupInit()
	// sendconfig2consul()
	// configgen.Getconfig = getTestConfig

	exitCode := m.Run()
	teardown()
	// // 退出
	os.Exit(exitCode)
}
