version := $(shell /bin/date "+%Y-%m-%d %H:%M")

build:
	go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o loadfromyaml main.go
	$(if $(shell command -v upx), upx loadfromyaml)
mac:
	GOOS=darwin go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o loadfromyaml-darwin main.go
	$(if $(shell command -v upx), upx loadfromyaml-darwin)
win:
	GOOS=windows go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o loadfromyaml.exe main.go
	$(if $(shell command -v upx), upx loadfromyaml.exe)
linux:
	GOOS=linux go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o loadfromyaml-linux main.go
	$(if $(shell command -v upx), upx loadfromyaml-linux)
