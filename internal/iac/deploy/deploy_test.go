package deploy

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/max-gui/consulagent/pkg/consulsets"
	"github.com/max-gui/logagent/pkg/confload"
	"github.com/max-gui/logagent/pkg/logsets"
	"github.com/max-gui/spells/internal/iac/archfig"
	"github.com/max-gui/spells/internal/iac/dockfig"
	"github.com/max-gui/spells/internal/iac/jenfig"
	"github.com/max-gui/spells/internal/pkg/constset"
	"gopkg.in/yaml.v2"
)

var orgconfigPth, abstestpath, pthSep string

func setup() {

	// flag.Parse()
	// plaintext = "123"
	// cryptedHexText = "1bda1896724a4521cfb7f38646824197929cd1"
	*logsets.Apppath = "/Users/jimmy/Projects/hercules/spells"
	*logsets.Port = "8080"
	*logsets.Appenv = "test"

	flag.Parse()
	bytes := confload.Load(context.Background())
	constset.StartupInit(bytes, context.Background())

	*consulsets.Consul_host = ""
	*consulsets.Acltoken = "" //"245d0a09-7139-config-prod-ff170a0562b1"
	// testpath = makeconfiglist()
	pthSep = string(os.PathSeparator)
	orgconfigPth = abstestpath + pthSep + "orgconfig" + pthSep
}

func teardown() {

}

func Test_genjenkinsmap(t *testing.T) {
	c := context.Background()
	appconf := archfig.GetArchfigSin("af-orderservice", c)
	genjenkinsmap(appconf, archfig.EnvInfo{Env: "test", Dc: "LFB"}, "test", "default", "default", c)
}

func Test_genConfAppend(t *testing.T) {
	c := context.Background()
	appconf := archfig.Arch_config{}
	appconf.Application.Name = "arch-spells"
	appconf.Environment.Tag = make(map[string]string)
	genConfAppend(appconf, "test", c)
	log.Print(appconf.Environment.Tag)
	genConfAppend(appconf, "uat", c)
	log.Print(appconf.Environment.Tag)
	genConfAppend(appconf, "prod", c)
	log.Print(appconf.Environment.Tag)
}

func Test_gen002(t *testing.T) {
	c := context.Background()

	filebytes, _ := os.ReadFile("/Users/jimmy/Downloads/git2consul.yaml")

	ddd := archfig.GenArchConfig(filebytes, "arch", "hercules", "git2consul", false, c)

	// ddd := archfig.GenArchConfigFrom("team0/project0/iactest2.yaml", "team0", "project0", "iactest2", false, c)
	ddd.Application.NoSource = true
	bytes, _ := yaml.Marshal(ddd)
	log.Println(string(bytes))

	dockerfile := dockfig.GenDocfile(ddd, c)
	log.Println(dockerfile)

	jenconfig := jenfig.GenJenfig(ddd)
	jenfile := jenfig.GenJenfile(jenconfig, c)
	log.Println(jenfile)

	jenkinspara := genjenkinsmap(ddd, archfig.EnvInfo{}, "prod", "default", "master", c)
	jenkinspara["archfig"] = ""
	jenkinspara["Dockerstring"] = ""
	jenkinspara["Valstring"] = ""
	log.Println("-------------------------")
	log.Println(jenkinspara)
	// bytes, _ = yaml.Marshal(Gen4old("fls-afch", "afch-authservice", c))
	// log.Print(string(bytes))
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
