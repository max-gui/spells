package defig

import (
	"context"
	"os"
	"testing"

	"github.com/max-gui/spells/internal/pkg/constset"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func setup() {
	constset.StartupInit(nil, context.Background())
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
	// bytes, err := ioutil.ReadFile(constset.Defconfpath)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// filemap := make(map[interface{}]interface{})
	// err = yaml.Unmarshal(bytes, &filemap)
	// if err != nil {
	// 	log.Panic(err)
	// }
	c := context.Background()
	defigval := GenDefig(false, c)
	defigbytes, _ := yaml.Marshal(defigval)
	defigval4byte := GenDefigFrom(defigbytes, false, c)

	assert.Equal(t, defigval, defigval4byte)
}

func Test_GetDefconfig(t *testing.T) {
	defconfig = Defconf{}
	c := context.Background()
	// *constset.RedisServer = ""
	getdefig := GetDefconfig(c)
	defconfig = Defconf{}

	gendefig := GenDefig(false, c)
	assert.Equal(t, gendefig, getdefig)
}
