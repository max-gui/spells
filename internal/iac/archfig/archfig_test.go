package archfig

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/max-gui/logagent/pkg/confload"
	"github.com/max-gui/spells/internal/pkg/constset"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func setup() {

	flag.Parse()
	bytes := confload.Load()
	constset.StartupInit(bytes, context.Background())
}

func teardown() {

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

func Test_GenDefig(t *testing.T) {
	// RMArch()
	bytes, err := ioutil.ReadFile(constset.Archpath + "team0/project0/iactest2.yaml")
	c := context.Background()
	if err != nil {
		log.Panic(err)
	}
	fc := FileContentInfo{
		Path:    "team0/project0/iactest2.yaml",
		Content: string(bytes),
	}
	archfigget := GetArchfig(fc, c)

	appfigget := GetArchfigByGitContent(fc, false, c)

	assert.Equal(t, appfigget, archfigget)
}

func Test_AppExist(t *testing.T) {
	c := context.Background()
	flat, paths := AppExist("iactest2", c)
	assert.Equal(t, true, flat)
	assert.Equal(t, []string{"ops/iac/arch/team0/project0/iactest2"}, paths)

	flat, paths = AppExist("iactest22", c)
	assert.Equal(t, false, flat)
	assert.Equal(t, []string{}, paths)
}

func Test_RMArch(t *testing.T) {
	c := context.Background()
	app_conf := Arch_config{}
	team := "team"
	proj := "proj"
	appfname := "app"
	app_conf.Application.Name = appfname
	app_conf.Application.Appid = "SS"
	app_conf.Application.Repositry = "ssh"
	app_conf.Environment.Strategy = map[string]struct {
		Capacity string `yaml:"capacity,omitempty"`
		Cpu      string `yaml:"cpu,omitempty"`
		Mem      string `yaml:"mem,omitempty"`
		Replica  int    `yaml:"replica,omitempty"`
	}{"test": {Capacity: "mini"}}
	app_conf.Deploy.Build.Output = "dd"

	bytes, _ := yaml.Marshal(app_conf)
	app_conf = GenArchConfig(bytes, team, proj, appfname, true, c)

	flag, _ := AppExist(app_conf.Application.Name, c)
	assert.True(t, flag)

	app_conf.RMArch(c)
	flag, _ = AppExist(app_conf.Application.Name, c)
	assert.False(t, flag)
}
func Test_GenArchConfigFrom(t *testing.T) {
	c := context.Background()
	app_conf := Arch_config{}
	// arfig.Application.UnGenfig = true

	bytes, _ := yaml.Marshal(app_conf)
	team := "team"
	proj := "proj"
	appfname := "app"

	assert.PanicsWithValue(t,
		fmt.Sprintf("team or project should not be empty, get:\nteam:%s project:%s", app_conf.Application.Team, app_conf.Application.Project),
		func() { GenArchConfig(bytes, "", "", "", false, c) })

	assert.PanicsWithValue(t,
		fmt.Sprintf("application name shoud be equal with filename, get:\napplication name:%s appfilename:%s", app_conf.Application.Name, appfname),
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Application.Name = appfname
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		fmt.Sprintf("application Appid shoud not be empty,captitalism or with space, get:\napplication Appid:%s", app_conf.Application.Appid),
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Application.Appid = "SS"
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		fmt.Sprintf("application Appid shoud not be empty,captitalism or with space, get:\napplication Appid:%s", app_conf.Application.Appid),
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Application.Appid = "s s"
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		fmt.Sprintf("application Appid shoud not be empty,captitalism or with space, get:\napplication Appid:%s", app_conf.Application.Appid),
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Application.Appid = "ss"
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		fmt.Sprintf("application Type shoud not be empty,captitalism or with space, get:\n%s", app_conf.Application.Type),
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Application.Type = "Java"
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		fmt.Sprintf("application Type shoud not be empty,captitalism or with space, get:\n%s", app_conf.Application.Type),
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Application.Type = " java"
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		fmt.Sprintf("application Type shoud not be empty,captitalism or with space, get:\n%s", app_conf.Application.Type),
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Application.Type = "java"
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		fmt.Sprintf("repo shoud use ssh or git, get:%s", app_conf.Application.Repositry),
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Application.Repositry = "ssh"
	app_conf.Environment.Expose.Unsafe = true
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		fmt.Sprintf("Expose.Path can't be empty if you want to expose a service, get:%s", app_conf.Environment.Expose.PrefixPath),
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Environment.Expose.PrefixPath = "/asdf"
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		"app_conf.Environment.Expose.Path should not start with /",
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Environment.Expose.PrefixPath = "asdf"
	app_conf.Environment.Strategy = map[string]struct {
		Capacity string `yaml:"capacity,omitempty"`
		Cpu      string `yaml:"cpu,omitempty"`
		Mem      string `yaml:"mem,omitempty"`
		Replica  int    `yaml:"replica,omitempty"`
	}{"dd": {}}
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		"env should be in range of [dr prod test uat], get:dd",
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Environment.Strategy = map[string]struct {
		Capacity string `yaml:"capacity,omitempty"`
		Cpu      string `yaml:"cpu,omitempty"`
		Mem      string `yaml:"mem,omitempty"`
		Replica  int    `yaml:"replica,omitempty"`
	}{"test": {}}
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		"test cpu config is wrong, get:",
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Environment.Strategy = map[string]struct {
		Capacity string `yaml:"capacity,omitempty"`
		Cpu      string `yaml:"cpu,omitempty"`
		Mem      string `yaml:"mem,omitempty"`
		Replica  int    `yaml:"replica,omitempty"`
	}{"test": {Cpu: "-1"}}
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		"test cpu config is wrong, get:-1",
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Environment.Strategy = map[string]struct {
		Capacity string `yaml:"capacity,omitempty"`
		Cpu      string `yaml:"cpu,omitempty"`
		Mem      string `yaml:"mem,omitempty"`
		Replica  int    `yaml:"replica,omitempty"`
	}{"test": {Cpu: "0.2", Mem: "11"}}
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		"test mem config is wrong, should with suffix 'Mi' get:11",
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Environment.Strategy = map[string]struct {
		Capacity string `yaml:"capacity,omitempty"`
		Cpu      string `yaml:"cpu,omitempty"`
		Mem      string `yaml:"mem,omitempty"`
		Replica  int    `yaml:"replica,omitempty"`
	}{"test": {Cpu: "0.2", Mem: "-1Mi"}}
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		"test mem config is wrong, get:-1Mi",
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Environment.Strategy = map[string]struct {
		Capacity string `yaml:"capacity,omitempty"`
		Cpu      string `yaml:"cpu,omitempty"`
		Mem      string `yaml:"mem,omitempty"`
		Replica  int    `yaml:"replica,omitempty"`
	}{"test": {Capacity: "dd"}}
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		"capacity should be in range of [high low mid mini], get:dd",
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Environment.Strategy = map[string]struct {
		Capacity string `yaml:"capacity,omitempty"`
		Cpu      string `yaml:"cpu,omitempty"`
		Mem      string `yaml:"mem,omitempty"`
		Replica  int    `yaml:"replica,omitempty"`
	}{"test": {Capacity: "mini"}}
	// app_conf.Deploy.Build.Output = "dd"
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		"Deploy.Build.Output should not be empty",
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

	app_conf.Environment.Strategy = map[string]struct {
		Capacity string `yaml:"capacity,omitempty"`
		Cpu      string `yaml:"cpu,omitempty"`
		Mem      string `yaml:"mem,omitempty"`
		Replica  int    `yaml:"replica,omitempty"`
	}{"test": {Capacity: "mini"}}
	app_conf.Deploy.Build.Output = "/dd"
	bytes, _ = yaml.Marshal(app_conf)
	assert.PanicsWithValue(t,
		"Deploy.Build.Output should not start with /",
		func() { GenArchConfig(bytes, team, proj, appfname, false, c) })

}

func Test_getContentInfo(t *testing.T) {
	c := context.Background()
	type args struct {
		c FileContentInfo
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getContentInfo(tt.args.c, c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getContentInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
