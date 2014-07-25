package templates

import (
	"html/template"
	"log"
)

var sources map[string]string
var T *template.Template

func MyEq(a, b interface{}) bool {
	return a == b
}

func init() {
	T = template.New("_top_").Funcs(template.FuncMap{"myeq": MyEq})
}

func registerTemplate(name string, source string) {
	if sources == nil {
		sources = make(map[string]string)
	}
	_, ok := sources[name]
	if ok {
		log.Fatalf("ERROR: duplicate template %s", name)
	}
	sources[name] = source
}

func Parse() {
	for name, source := range sources {
		_, err := T.New(name).Parse(source)
		if err != nil {
			log.Fatalf("ERROR: failed to parse template %s - %s", name, err)
		}
	}
}
