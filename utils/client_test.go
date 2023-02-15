package utils

import (
	"path"
	"testing"
)

func TestFtpDownload(t *testing.T){
	svr:=NewServerInfo("192.168.01.1","test")
	svr.TargetPath=path.Join(".", svr.Addr)
	if err:=svr.Work();err!=nil{
		t.Fatal(err)
	}
}
