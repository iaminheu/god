package gogen

import (
	"bytes"
	"strings"
	"text/template"

	"git.zc0901.com/go/god/tools/god/api/spec"
	"git.zc0901.com/go/god/tools/god/api/util"
	"git.zc0901.com/go/god/tools/god/config"
	"git.zc0901.com/go/god/tools/god/util/format"
)

var middlewareImplementCode = `
package middleware

import "net/http"

type {{.name}} struct {
}

func New{{.name}}() *{{.name}} {	
	return &{{.name}}{}
}

func (m *{{.name}})Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

		// Pass through to next handler if need 
		next(w, r)
	}	
}
`

func genMiddleware(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	var middlewares = getMiddleware(api)
	for _, item := range middlewares {
		middlewareFilename := strings.TrimSuffix(strings.ToLower(item), "middleware") + "_middleware"
		formatName, err := format.FileNamingFormat(cfg.NamingFormat, middlewareFilename)
		if err != nil {
			return err
		}

		filename := formatName + ".go"
		fp, created, err := util.MaybeCreateFile(dir, middlewareDir, filename)
		if err != nil {
			return err
		}
		if !created {
			return nil
		}
		defer fp.Close()

		name := strings.TrimSuffix(item, "Middleware") + "Middleware"
		t := template.Must(template.New("contextTemplate").Parse(middlewareImplementCode))
		buffer := new(bytes.Buffer)
		err = t.Execute(buffer, map[string]string{
			"name": strings.Title(name),
		})
		if err != nil {
			return err
		}

		formatCode := formatCode(buffer.String())
		_, err = fp.WriteString(formatCode)
		return err
	}
	return nil
}
