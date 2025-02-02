package generator

import (
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/stringx"
	conf "git.zc0901.com/go/god/tools/god/config"
	"git.zc0901.com/go/god/tools/god/rpc/execx"
	"github.com/stretchr/testify/assert"
)

var cfg = &conf.Config{
	NamingFormat: "gozero",
}

func TestRpcGenerate(t *testing.T) {
	_ = Clean()
	dispatcher := NewDefaultGenerator()
	err := dispatcher.Prepare()
	if err != nil {
		logx.Error(err)
		return
	}
	projectName := stringx.Rand()
	g := NewRpcGenerator(dispatcher, cfg)

	// case go path
	src := filepath.Join(build.Default.GOPATH, "src")
	_, err = os.Stat(src)
	if err != nil {
		return
	}

	projectDir := filepath.Join(src, projectName)
	srcDir := projectDir
	defer func() {
		_ = os.RemoveAll(srcDir)
	}()
	err = g.Generate("./test.proto", projectDir, []string{src})
	assert.Nil(t, err)
	_, err = execx.Run("go test "+projectName, projectDir)
	if err != nil {
		assert.Contains(t, err.Error(), "not in GOROOT")
	}

	// case go mod
	workDir := t.TempDir()
	name := filepath.Base(workDir)
	_, err = execx.Run("go mod init "+name, workDir)
	if err != nil {
		logx.Error(err)
		return
	}

	projectDir = filepath.Join(workDir, projectName)
	err = g.Generate("./test.proto", projectDir, []string{src})
	assert.Nil(t, err)
	_, err = execx.Run("go test "+projectName, projectDir)
	if err != nil {
		assert.Contains(t, err.Error(), "not in GOROOT")
	}

	// case not in go mod and go path
	err = g.Generate("./test.proto", projectDir, []string{src})
	assert.Nil(t, err)
	_, err = execx.Run("go test "+projectName, projectDir)
	if err != nil {
		assert.Contains(t, err.Error(), "not in GOROOT")
	}

	// invalid directory
	projectDir = filepath.Join(t.TempDir(), ".....")
	err = g.Generate("./test.proto", projectDir, nil)
	assert.NotNil(t, err)
}
