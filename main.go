package main

import (
	"context"
	"ftpHelper/conf"
	"ftpHelper/utils"
	"golang.org/x/sync/semaphore"
	"log"
	"path"
	"runtime"
)

var (
	maxWorkers = runtime.GOMAXPROCS(0) * 2
	//设置为cpu的核数*2
	sema = semaphore.NewWeighted(int64(maxWorkers))
	ctx  = context.Background()
)

func main() {
	if err := conf.InitConfig(""); err != nil {
		log.Printf("init config error:%v\n", err)
		return
	}
	log.Printf("max workers:%d\n", maxWorkers)

	for _, addr := range utils.DefaultCommonInfo.ServerList {
		svrInfo := utils.NewServerInfo(addr)
		svrInfo.CommonInfo = utils.DefaultCommonInfo
		svrInfo.SinglePath = path.Join(svrInfo.StorePath, addr)
		log.Println(svrInfo.SinglePath)
		if err := sema.Acquire(ctx, 1); err != nil {
			break
		}
		go func() {
			defer sema.Release(1)
			svrInfo.Work()
		}()
	}

	if err := sema.Acquire(ctx, int64(maxWorkers)); err != nil {
		log.Printf("all work exec failed:%v\n", err)
		return
	}
	log.Println("work all success")
}
