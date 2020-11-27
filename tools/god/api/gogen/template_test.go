package gogen

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"git.zc0901.com/go/god/tools/god/util"
	"github.com/stretchr/testify/assert"
)

func TestGenTemplates(t *testing.T) {
	err := util.InitTemplates(category, templates)
	assert.Nil(t, err)
	dir, err := util.GetTemplateDir(category)
	assert.Nil(t, err)
	file := filepath.Join(dir, "main.tpl")
	data, err := ioutil.ReadFile(file)
	assert.Nil(t, err)
	assert.Equal(t, string(data), mainTemplate)
}

func TestRevertTemplate(t *testing.T) {
	name := "main.tpl"
	err := util.InitTemplates(category, templates)
	assert.Nil(t, err)

	dir, err := util.GetTemplateDir(category)
	assert.Nil(t, err)

	file := filepath.Join(dir, name)
	data, err := ioutil.ReadFile(file)
	assert.Nil(t, err)

	modifyData := string(data) + "modify"
	err = util.CreateTemplate(category, name, modifyData)
	assert.Nil(t, err)

	data, err = ioutil.ReadFile(file)
	assert.Nil(t, err)

	assert.Equal(t, string(data), modifyData)

	assert.Nil(t, RevertTemplate(name))

	data, err = ioutil.ReadFile(file)
	assert.Nil(t, err)
	assert.Equal(t, mainTemplate, string(data))
}

func TestClean(t *testing.T) {
	name := "main.tpl"
	err := util.InitTemplates(category, templates)
	assert.Nil(t, err)

	assert.Nil(t, Clean())

	dir, err := util.GetTemplateDir(category)
	assert.Nil(t, err)

	file := filepath.Join(dir, name)
	_, err = ioutil.ReadFile(file)
	assert.NotNil(t, err)
}

func TestUpdate(t *testing.T) {
	name := "main.tpl"
	err := util.InitTemplates(category, templates)
	assert.Nil(t, err)

	dir, err := util.GetTemplateDir(category)
	assert.Nil(t, err)

	file := filepath.Join(dir, name)
	data, err := ioutil.ReadFile(file)
	assert.Nil(t, err)

	modifyData := string(data) + "modify"
	err = util.CreateTemplate(category, name, modifyData)
	assert.Nil(t, err)

	data, err = ioutil.ReadFile(file)
	assert.Nil(t, err)

	assert.Equal(t, string(data), modifyData)

	assert.Nil(t, Update(category))

	data, err = ioutil.ReadFile(file)
	assert.Nil(t, err)
	assert.Equal(t, mainTemplate, string(data))
}