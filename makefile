GOBUILD = go build -ldflags="-w -s"
WIN_TARGET = ftpFetch.exe
TARGET = ftpFetch
build : build-windows

build-linux: main.go config.yaml
	 $(GOBUILD)  -o $(TARGET)  main.go

build-windows: main.go config.yaml
	GOOS=windows GOARCH=amd64  GO111MODULE=on $(GOBUILD)  -o $(WIN_TARGET)  main.go

clean :  cleanlog  cleanprogram
cleanlog :
	rm -f *.log
cleanprogram :
	rm -f ${WIN_TARGET} ${TARGET}


.PHONY: clean