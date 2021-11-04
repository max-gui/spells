package deploy

import (
	"context"
	"log"
	"testing"

	"github.com/max-gui/spells/internal/iac/archfig"
)

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
