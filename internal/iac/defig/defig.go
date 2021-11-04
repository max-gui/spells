package defig

import (
	"context"
	"encoding/json"
	"io/ioutil"
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
			Args map[string]string
			Ign  map[string][]string
		}

		Build map[string]struct {
			Cmd       string
			Arg       string
			Config    string
			Jenkignor []string
		}
		Output string

		Resource map[string][]string
		// Volumes  []string
		Capacity map[string]struct {
			Capacity string
			Replica  int
			Cpu      string
			Mem      string
		}
		Deploy struct {
			Limited  []string
			Strategy []struct {
				Flow string
				Env  []string
			}
		}
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

// var dockerpath = constset.Apppath + constset.PthSep + "repo" + constset.PthSep + "temple" + constset.PthSep
func GenDefig(isinstall bool, c context.Context) Defconf {

	log := logagent.Inst(c)
	dirPth := constset.Defconfpath // constset.Templepath + "defaultconfig.yaml"
	str, err := ioutil.ReadFile(dirPth)
	if err != nil {
		log.Panic(err)
	}
	return GenDefigFrom(str, isinstall, c)

	// Defconfig = defconf
	// return defconf
}

func GenDefigFrom(content []byte, isinstall bool, c context.Context) Defconf {

	log := logagent.Inst(c)
	log.Println(string(content))
	defconf := Defconf{}
	yaml.Unmarshal(content, &defconf)

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
