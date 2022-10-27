package valfig

import (
	"context"
	"log"

	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/spells/internal/iac/archfig"
	"github.com/max-gui/spells/internal/iac/templ"
	"gopkg.in/yaml.v3"
)

type Sidefig struct {
	Appconf    archfig.Arch_config
	Values     map[string]interface{}
	Neighbours []string
}

func GenSideContent(sides Sidefig, c context.Context) []interface{} {
	// dirPth := orgconfigPth + pthSep + "arch.yaml"
	// str, _ := ioutil.ReadFile(dirPth)
	// a := Arch_config{}
	// yaml.Unmarshal(str, &a)
	// t.Log(a)
	logger := logagent.InstPlatform(c)

	neighbourmap := []interface{}{}
	for _, neighbour := range sides.Neighbours {

		result := templ.GemplFromType(neighbour+".yml", "sidecar", sides.Values, c)
		log.Println(result)
		neighbourinfo := map[string]interface{}{}
		err := yaml.Unmarshal([]byte(result), &neighbourinfo)

		if err != nil {
			logger.Panic(err)
		}
		neighbourmap = append(neighbourmap, neighbourinfo)
	}

	return neighbourmap
}

func GenSidefig(appconf archfig.Arch_config, valuesinfo map[string]interface{}, env, dc string, c context.Context) Sidefig {

	ignmap := func(slice []string) map[string]struct{} {
		res := map[string]struct{}{}
		for _, v := range slice {
			res[v] = struct{}{}
		}

		return res
	}(appconf.Deploy.Sidecar.Ign[env])

	neighbours := []string{}
	for _, neighbour := range appconf.Deploy.Sidecar.Neighbour {
		if _, ok := ignmap[neighbour]; !ok {
			neighbours = append(neighbours, neighbour)
		}
	}

	sidefiginfo := Sidefig{
		Appconf:    appconf,
		Values:     valuesinfo,
		Neighbours: neighbours,
	}

	return sidefiginfo

}
