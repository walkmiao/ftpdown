package utils

import (
	"backup/conf"
	"backup/models"
	"log"
	"os"
	"testing"
)

func preLoad() {
	if err := conf.InitConfig(""); err != nil {
		log.Printf("init config error:%v", err)
		os.Exit(1)
	}
}
func TestFtpDownload(t *testing.T) {
	preLoad()
	server := &models.Collector{
		EquipName: "TESTUC01",
		Addr:      "192.168.8.102:21",
	}
	svrInfo := NewServerInfo(server, "外高桥")
	if err := svrInfo.Work(); err != nil {
		t.Fatal(err)
	}
}
