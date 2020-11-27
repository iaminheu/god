version := $(shell /bin/date "+%Y-%m-%d %H:%M")

mac:
	GOOS=darwin go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o god god.go
	command -v upx &> /dev/null && upx god
	mv god ~/go/bin/