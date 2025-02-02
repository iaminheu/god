package tsgen

import (
	"errors"
	"git.zc0901.com/go/god/lib/fs"
	"path"
	"strings"
	"text/template"

	"git.zc0901.com/go/god/tools/god/api/spec"
	apiutil "git.zc0901.com/go/god/tools/god/api/util"
)

const (
	componentsTemplate = `// Code generated by god. DO NOT EDIT.

{{.componentTypes}}
`
)

func genComponents(dir string, api *spec.ApiSpec) error {
	types := apiutil.GetSharedTypes(api)
	if len(types) == 0 {
		return nil
	}

	val, err := buildTypes(types, func(name string) (*spec.Type, error) {
		for _, ty := range api.Types {
			if strings.ToLower(ty.Name) == strings.ToLower(name) {
				return &ty, nil
			}
		}
		return nil, errors.New("inline type " + name + " not exist, please correct api file")
	})
	if err != nil {
		return err
	}

	outputFile := apiutil.ComponentName(api) + ".ts"
	filename := path.Join(dir, outputFile)
	if err := fs.RemoveIfExist(filename); err != nil {
		return err
	}

	fp, created, err := apiutil.MaybeCreateFile(dir, ".", outputFile)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	t := template.Must(template.New("componentsTemplate").Parse(componentsTemplate))
	return t.Execute(fp, map[string]string{
		"componentTypes": val,
	})
}

func buildTypes(types []spec.Type, inlineType func(string) (*spec.Type, error)) (string, error) {
	var builder strings.Builder
	first := true
	for _, tp := range types {
		if first {
			first = false
		} else {
			builder.WriteString("\n")
		}
		if err := writeType(&builder, tp, func(name string) (*spec.Type, error) {
			return inlineType(name)
		}, func(tp string) string {
			return ""
		}); err != nil {
			return "", apiutil.WrapErr(err, "Type "+tp.Name+" generate error")
		}
	}

	return builder.String(), nil
}
