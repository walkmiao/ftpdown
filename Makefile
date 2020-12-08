GOBUILD = go build -ldflags="-w -s"
WIN_TARGET64 = fetch64.exe
WIN_TARGET32 = fetch32.exe
TARGET = fetch

build : build-windows32

build-linux : main.go config.yaml
	GOOS=linux GOARCH=amd64 GO111MODULE=on $(GOBUILD)  -o $(TARGET)  main.go

build-windows64 : main.go config.yaml
	GOOS=windows GOARCH=amd64  GO111MODULE=on $(GOBUILD)  -o $(WIN_TARGET64)  main.go

build-windows32 : main.go config.yaml
	GOOS=windows GOARCH=386  GO111MODULE=on $(GOBUILD)  -o $(WIN_TARGET32)  main.go

clean : cleanlog cleanprogram

cleanlog :
	rm -f *.log

cleanprogram :
	rm -f ${WIN_TARGET64} ${WIN_TARGET32} ${TARGET}

.PHONY : clean