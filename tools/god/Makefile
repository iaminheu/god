version := $(shell /bin/date "+%Y-%m-%d %H:%M")

mac:
	GOOS=darwin go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o god god.go
	$(if $(shell command -v upx), upx god)
	mv god ~/go/bin/

linux:
	GOOS=linux go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o god-linux god.go
    $(if $(shell command -v upx), upx god-linux)
