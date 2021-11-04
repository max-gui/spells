package templ

import (
	"bytes"
	"context"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/max-gui/fileconvagt/pkg/fileops"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/redisagent/pkg/redisops"
	"github.com/max-gui/spells/internal/pkg/constset"
)

func GemplFrom(templname string, a interface{}, c context.Context) string {
	// rediscli := redisops.Pool().Get()

	// defer rediscli.Close()

	if v, ok := templMap[templname]; ok {
		// rediscli.Do("EXPIRE", "arch-spell-temple", 600)
		return Gempl(v, a, c)
	} else {
		dirPth := constset.Templepath + templname
		res := GetemplFrom(dirPth, templname, true, c)
		return Gempl(res, a, c)
		// var templtmp *template.Template
		// mm, err := redis.Bytes(rediscli.Do("HGET", "arch-spell-temple", templname))
		// json.Unmarshal(mm, &templtmp)

		// if err != nil {
		// 	dirPth := constset.Templepath + templname
		// 	res := GetemplFrom(dirPth, templname, true)
		// 	return Gempl(res, a)
		// } else {
		// 	templMap[templname] = templtmp
		// 	rediscli.Do("EXPIRE", "arch-spell-temple", 600)
		// 	return Gempl(templtmp, a)
		// }
	}
}

func ClsGempl() {
	templMap = make(map[string]*template.Template)

	// rediscli := redisops.Pool().Get()

	// defer rediscli.Close()
	// rediscli.Do("DEL", "arch-spell-temple")
}

func Gempl(templ *template.Template, a interface{}, c context.Context) string {
	// dirPth = abstestpath + PthSep + "orgconfig" + PthSep + "Dockerfile." + a.Application.Type
	// templ := fn1(dirPth)
	var buffer bytes.Buffer
	log := logagent.Inst(c)

	err := templ.Execute(&buffer, a)
	if err != nil {
		log.Panic(err)
	}

	bs := buffer.String()
	return bs
}

var templMap = make(map[string]*template.Template)

func Getempl(tempfile string, templname string, isinstall bool, c context.Context) *template.Template {

	// var buffer bytes.Buffer
	log := logagent.Inst(c)
	templ, err := template.New(templname).Funcs(sprig.TxtFuncMap()).Parse(tempfile) //.Funcs(sprig.FuncMap()).Parse
	// templ, err := template.New(templname).Parse(tempfile)
	if err != nil {
		log.Panic(err)
	}
	if isinstall {
		templMap[templname] = templ

		// rediscli := redisops.Pool().Get()

		// defer rediscli.Close()
		// jsonbs, _ := json.Marshal(templ)
		// rediscli.Do("HSET", "arch-spell-templ", templname, jsonbs)

		// rediscli.Do("EXPIRE", "arch-spell-temple", 600)
	}
	return templ
}

func RMtempl(templname string) {
	delete(templMap, templname)
	rediscli := redisops.Pool().Get()

	defer rediscli.Close()
	rediscli.Do("HDEL", "arch-spell-templ", templname)
}

func GetemplFrom(dirPth, templname string, isinstall bool, c context.Context) *template.Template {
	var err error
	var f string
	log := logagent.Inst(c)
	f, err = fileops.Read(dirPth)
	if err != nil {
		log.Panic(err)
	}
	// log.Print(f)
	return Getempl(f, templname, isinstall, c)
}
