package ctx

import (
	"git.zc0901.com/go/god/lib/fs"
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"git.zc0901.com/go/god/lib/stringx"
	"git.zc0901.com/go/god/tools/god/rpc/execx"
	"github.com/stretchr/testify/assert"
)

func TestProjectFromGoMod(t *testing.T) {
	dft := build.Default
	gp := dft.GOPATH
	if len(gp) == 0 {
		return
	}
	projectName := stringx.Rand()
	dir := filepath.Join(gp, "src", projectName)
	err := fs.MkdirIfNotExist(dir)
	if err != nil {
		return
	}

	_, err = execx.Run("go mod init "+projectName, dir)
	assert.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	ctx, err := projectFromGoMod(dir)
	assert.Nil(t, err)
	assert.Equal(t, projectName, ctx.Path)
	assert.Equal(t, dir, ctx.Dir)
}
