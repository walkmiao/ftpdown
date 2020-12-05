package main

import (
	"context"
	"ftpHelper/conf"
	"ftpHelper/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
	"log"
	"path"
	"runtime"
	"time"
)

var (
	maxWorkers = runtime.GOMAXPROCS(0) * 2
	//设置为cpu的核数*2
	sema = semaphore.NewWeighted(int64(maxWorkers))
	ctx,_  = context.WithTimeout(context.Background(),time.Second*5)
)

func main() {
	if err := conf.InitConfig(""); err != nil {
		log.Printf("init config error:%v\n", err)
		return
	}
	for _, addr := range utils.DefaultCommonInfo.ServerList {
		svrInfo := utils.NewServerInfo(addr)
		svrInfo.CommonInfo = utils.DefaultCommonInfo
		svrInfo.TargetPath = path.Join(svrInfo.StorePath, addr)
		svrInfo.PreviousDir = svrInfo.TargetPath
		svrInfo.Debugf("target store path:%s",svrInfo.TargetPath)
		if err := sema.Acquire(ctx, 1); err != nil {
			break
		}
		go func() {
			if err:=svrInfo.Work();err!=nil{
				svrInfo.Errorf("download from server error:%v\n",err)
			}
			sema.Release(1)
		}()
	}

	if err := sema.Acquire(ctx, int64(maxWorkers)); err != nil {
		logrus.Errorf("work  do  failed:%v\n", err)
		return
	}
	logrus.Info(" All work done!")
}
