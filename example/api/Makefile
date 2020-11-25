version := $(shell /bin/date "+%Y-%m-%d %H:%M")

mac:
	GOOS=darwin go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o example main.go
	command -v upx &> /dev/null && upx example