package tpl

var (
	Imports = `import (
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

	"god/lib/store/cache"
	"god/lib/store/sqlx"
	"god/lib/stringx"
	"god/tools/god/mysql/builder"
)
`

	ImportsNoCache = `import (
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

	"god/lib/store/sqlx"
	"god/lib/stringx"
	"god/tools/god/mysql/builder"
)`
)
