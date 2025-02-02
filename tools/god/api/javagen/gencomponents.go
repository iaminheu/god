package javagen

import (
	"fmt"
	"git.zc0901.com/go/god/lib/fs"
	"io"
	"path"
	"strings"
	"text/template"

	"git.zc0901.com/go/god/tools/god/api/spec"
	apiutil "git.zc0901.com/go/god/tools/god/api/util"
	"git.zc0901.com/go/god/tools/god/util"
)

const (
	componentTemplate = `// Code generated by god. DO NOT EDIT.
package com.god.logic.http.packet.{{.packet}}.model;

import com.god.logic.http.DeProguardable;

{{.componentType}}
`
)

func genComponents(dir, packetName string, api *spec.ApiSpec) error {
	types := apiutil.GetSharedTypes(api)
	if len(types) == 0 {
		return nil
	}
	for _, ty := range types {
		if err := createComponent(dir, packetName, ty); err != nil {
			return err
		}
	}

	return nil
}

func createComponent(dir, packetName string, ty spec.Type) error {
	modelFile := util.Title(ty.Name) + ".java"
	filename := path.Join(dir, modelDir, modelFile)
	if err := fs.RemoveOrQuit(filename); err != nil {
		return err
	}

	fp, created, err := apiutil.MaybeCreateFile(dir, modelDir, modelFile)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	tys, err := buildType(ty)
	if err != nil {
		return err
	}

	t := template.Must(template.New("componentType").Parse(componentTemplate))
	return t.Execute(fp, map[string]string{
		"componentType": tys,
		"packet":        packetName,
	})
}

func buildType(ty spec.Type) (string, error) {
	var builder strings.Builder
	if err := writeType(&builder, ty); err != nil {
		return "", apiutil.WrapErr(err, "Type "+ty.Name+" generate error")
	}
	return builder.String(), nil
}

func writeType(writer io.Writer, tp spec.Type) error {
	fmt.Fprintf(writer, "public class %s implements DeProguardable {\n", util.Title(tp.Name))
	for _, member := range tp.Members {
		if err := writeProperty(writer, member, 1); err != nil {
			return err
		}
	}
	genGetSet(writer, tp, 1)
	fmt.Fprintf(writer, "}\n")
	return nil
}
