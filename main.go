package main

import (
	"context"
	"fmt"
	"ftpHelper/conf"
	"ftpHelper/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/sync/semaphore"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sync"
	"time"
)

var (
	cpus = runtime.GOMAXPROCS(0)
	maxWorkers = 8
	wg = &sync.WaitGroup{}
)

func main() {
	if err := conf.InitConfig(""); err != nil {
		log.Printf("init config error:%v\n", err)
		return
	}

	num:=viper.GetInt("Fetch.Factor")
	if num!=0{
		maxWorkers = cpus*num
	}else{
		maxWorkers = cpus*2
	}
	timeout:=viper.GetInt("Fetch.Timeout")
	ctx,_  := context.WithTimeout(context.Background(),
		time.Second*time.Duration(timeout))
	log.Printf("cpu:%d factor:%d maxWorkers(cpu*factor):%d\n",cpus,num,maxWorkers)
	log.Printf("timeout is set to %d s!\n",timeout)
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		if err:=HandleServerList(ctx,utils.DefaultCommonInfo.JDServer,"jd");err!=nil{
			logrus.Errorf("handle server list jd error:%v\n",err)
		}
	}(wg)


	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		if err:=HandleServerList(ctx,utils.DefaultCommonInfo.WGQServer,"wgq");err!=nil{
			logrus.Errorf("handle server list wgq error:%v\n",err)
		}
	}(wg)
	wg.Wait()
	logrus.Info(" All work done!")
	pause()
}

func HandleServerList(ctx context.Context,serverList []string,serverDir string)(err error){
	//设置为cpu的核数*2
	if len(serverList)<=0{
		return fmt.Errorf("%s server list is empty!\n",serverDir)
	}

	sema := semaphore.NewWeighted(int64(maxWorkers))

	for _, addr := range serverList {
		svrInfo := utils.NewServerInfo(addr,serverDir)
		svrInfo.CommonInfo = utils.DefaultCommonInfo
		svrInfo.TargetPath = path.Join(svrInfo.StorePath,serverDir,addr)
		svrInfo.PreviousDir = svrInfo.TargetPath
		svrInfo.Debugf("target store path:%s",svrInfo.TargetPath)
		if err := sema.Acquire(ctx, 1); err != nil {
			svrInfo.Errorf("accquire resource error:%v\n",err)
			break
		}
		go func() {
			if err:=svrInfo.Work();err!=nil{
				svrInfo.Errorf("download from server error:%v\n",err)
			}
			sema.Release(1)
			svrInfo.Warningf("%s released!\n",svrInfo.Addr)
		}()
	}

	if err := sema.Acquire(ctx, int64(maxWorkers)); err != nil {
		return fmt.Errorf("%s work  do  failed:%v\n", serverDir,err)
	}
	logrus.Infof("[DONE] %s  work done!\n",serverDir)
	return nil
}


func pause(){
	cmd:=exec.Command("cmd","/c","pause")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}