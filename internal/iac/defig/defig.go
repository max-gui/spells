package defig

import (
	"context"
	"encoding/json"
	"os"
	"sort"

	"github.com/gomodule/redigo/redis"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/redisagent/pkg/redisops"
	"github.com/max-gui/spells/internal/pkg/constset"
	"gopkg.in/yaml.v2"
)

type Defconf struct {
	Capacitylable map[string]struct {
		Cpu string `yaml:"cpu"`
		Mem string `yaml:"mem"`
	}
	Defualtinfo struct {
		Cmdarg struct {
			Args map[string][]string `yaml:"args,omitempty"`
			Ign  map[string][]string `yaml:"ign,omitempty"`
		} `yaml:"cmdarg,omitempty"`

		Build map[string]struct {
			Cmd       string   `yaml:"cmd,omitempty"`
			Arg       string   `yaml:"arg,omitempty"`
			Config    string   `yaml:"config,omitempty"`
			Jenkignor []string `yaml:"jenkignor,omitempty"`
		} `yaml:"build,omitempty"`

		Sidecar struct {
			Neighbour map[string][]string `yaml:"neighbour,omitempty"`
			Ign       map[string][]string `yaml:"ign,omitempty"`
		} `yaml:"sidecar,omitempty"`
		Output   string              `yaml:"output,omitempty"`
		Port     string              `yaml:"port,omitempty"`
		Resource map[string][]string `yaml:"resource,omitempty"`
		// Volumes  []string
		Capacity map[string]struct {
			Capacity string `yaml:"capacity,omitempty"`
			Replica  int    `yaml:"replica,omitempty"`
			Cpu      string `yaml:"cpu,omitempty"`
			Mem      string `yaml:"mem,omitempty"`
		} `yaml:"capacity,omitempty"`
		Deploy struct {
			Limited  []string `yaml:"limited,omitempty"`
			Strategy []struct {
				Flow string   `yaml:"flow,omitempty"`
				Env  []string `yaml:"env,omitempty"`
			} `yaml:"strategy,omitempty"`
		} `yaml:"deploy,omitempty"`
	}
}

var defconfig Defconf

func GetDefconfig(c context.Context) Defconf {
	rediscli := redisops.Pool().Get()
	defer rediscli.Close()

	if len(defconfig.Capacitylable) > 0 {

		rediscli.Do("EXPIRE", "arch-spell-default", 600)
		return defconfig
	} else {

		mm, err := redis.Bytes(rediscli.Do("GET", "arch-spell-default"))
		var defconf Defconf
		json.Unmarshal(mm, &defconf)

		if nil == err && len(defconfig.Capacitylable) > 0 {
			defconfig = defconf
			rediscli.Do("EXPIRE", "arch-spell-default", 600)
			return defconfig
		} else {
			return GenDefig(true, c)
		}
	}
}

func ClearDefconfig() {
	defconfig = Defconf{}
	rediscli := redisops.Pool().Get()
	defer rediscli.Close()
	rediscli.Do("EXPIRE", "arch-spell-default", 0)
}

// var dockerpath = constset.Apppath + constset.PthSep + "repo" + constset.PthSep + "temple" + constset.PthSep
func GenDefig(isinstall bool, c context.Context) Defconf {

	log := logagent.InstArch(c)
	dirPth := constset.Defconfpath // constset.Templepath + "defaultconfig.yaml"
	str, err := os.ReadFile(dirPth)
	if err != nil {
		log.Panic(err)
	}
	return GenDefigFrom(str, isinstall, c)

	// Defconfig = defconf
	// return defconf
}

func GenDefigFrom(content []byte, isinstall bool, c context.Context) Defconf {

	log := logagent.InstArch(c)
	log.Println(string(content))
	defconf := Defconf{}
	err := yaml.Unmarshal(content, &defconf)

	if err != nil {
		log.Panic(err)
	}

	if isinstall {

		rediscli := redisops.Pool().Get()

		defer rediscli.Close()
		jsonbs, _ := json.Marshal(defconf)

		rediscli.Do("SETEX", "arch-spell-default", 600, jsonbs)
		defconfig = defconf
	}
	// Defconfig = defconf
	return defconf
}

func (v *Defconf) GetEnvKeys() []string {
	envkeys := make([]string, 0, len(v.Capacitylable))
	for k := range v.Defualtinfo.Capacity {
		envkeys = append(envkeys, k)
	}
	sort.Strings(envkeys)
	return envkeys
}

func (v *Defconf) GetCapacity() []string {
	capkeys := make([]string, 0, len(v.Capacitylable))
	for k := range v.Capacitylable {
		capkeys = append(capkeys, k)
	}
	sort.Strings(capkeys)
	return capkeys
}
