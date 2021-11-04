package main

import (
	"flag"

	"github.com/max-gui/consulagent/pkg/consulhelp"
	"github.com/max-gui/logagent/pkg/confload"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/logagent/pkg/logsets"
	"github.com/max-gui/spells/internal/githelp"
	"github.com/max-gui/spells/internal/pkg/constset"
	"github.com/max-gui/spells/router"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func main() {

	// flag.Parse()
	// router.Test(os.Args[1])
	// var Argsetmap = make(map[string]interface{})

	flag.Parse()
	bytes := confload.Load()
	c := logagent.GetRootContextWithTrace()
	constset.StartupInit(bytes, c)
	// constset.StartupInit()
	// config := consulhelp.Getconfaml(*constset.ConfResPrefix, "redis", "redis-sentinel-proxy", *constset.Appenv)
	// redisops.Url = config["url"].(string)
	// redisops.Pwd = config["password"].(string)
	go consulhelp.StartWatch(*constset.ConfWatchPrefix, true, c)

	// if len(os.Args) > 2 {
	// 	port = os.Args[1]
	// }
	// if len(os.Args) > 2 {
	// 	consulhelp.Consulurl = os.Args[2]
	// }
	// if len(os.Args) > 3 {
	// 	consulhelp.AclToken = os.Args[3]
	// }

	// router.Envs
	//port := "4000"

	githelp.UpdateAll(c)
	r := router.SetupRouter()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	r.Run(":" + *logsets.Port)
}
