package utils

import (
	"path"
	"testing"
)

func TestFtpDownload(t *testing.T) {
	svr := NewServerInfo("192.168.1.127", "test")
	svr.FetchConf = FetchConf{
		port:       "21",
		username:   "root",
		password:   "Sbbdlyx123",
		ServerPath: ".",
	}
	svr.Filters = ".so|CSMain|yklog"

	svr.TargetPath = path.Join(".", svr.Addr)
	if err := svr.Work(); err != nil {
		t.Fatal(err)
	}
}
