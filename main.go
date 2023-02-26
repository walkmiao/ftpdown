package main

import (
	"backup/conf"
	"backup/logic"
	"backup/models"
	"backup/utils"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

var (
	cpus       = runtime.GOMAXPROCS(0)
	maxWorkers = cpus
	wg         = &sync.WaitGroup{}
)
var (
	Version string
	Build   string
)

func main() {
	logrus.Printf("version:%s build:%s\n", Version, Build)
	if err := conf.InitConfig(""); err != nil {
		logrus.Fatalf("init config error:%v\n", err)
	}
	num := conf.GlobalCfg.Fetch.Factor
	if num != 0 {
		maxWorkers = cpus * num
	} else {
		maxWorkers = cpus * 2
	}
	interval := time.Hour * time.Duration(conf.GlobalCfg.Fetch.Interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	var cs = make([]logic.QueryCollector, 0, len(conf.GlobalCfg.Mysql))
	for _, cnf := range conf.GlobalCfg.Mysql {
		cs = append(cs, &logic.MysqlQuery{MysqlConf: cnf})
	}
	//立即备份然后再定时
	backUp(wg, cs...)
	go func() {
		logrus.Infof("开启定时任务,间隔:%v", interval)
		for {
			select {
			case <-ticker.C:
				backUp(wg, cs...)
			}
		}
	}()
	sig := <-sigChan
	logrus.Infof("监听到信号%s，已结束运行", sig.String())
}

func HandleServerList(ctx context.Context, serverList []*models.Collector, park string) (err error) {
	total := len(serverList)
	var notTry int32 = 0
	sema := semaphore.NewWeighted(int64(maxWorkers))
	var num atomic.Int32
	for _, server := range serverList {
		v, ok := utils.ErrMap.Load(server.Addr)
		if ok {
			//如果已经被禁或者尝试的时间没到,则跳过下载
			if inv := v.(*utils.Invalid); inv.Forbidden || time.Now().Sub(inv.RetryTime()) < 0 {
				notTry++
				continue
			}
		}
		svrInfo := utils.NewServerInfo(server, park)
		if err = sema.Acquire(context.Background(), 1); err != nil {
			svrInfo.Errorf("semaphore accquire resource error:%v", err)
			break
		}
		server := server
		go func() {
			defer sema.Release(1)
			if err = svrInfo.Work(); err != nil {
				svrInfo.Errorf("备份KC521B【%s(%s)】出错:%v", server.EquipName, server.Addr, err)
				return
			}
			num.Add(1)
		}()
	}

	if err = sema.Acquire(context.Background(), int64(maxWorkers)); err != nil {
		return fmt.Errorf("semaphore accquire failed:%v", err)
	}
	success := num.Load()
	logrus.Infof("%s KC521B共【%d】台,成功:%d 台 失败:%d 台 跳过:%d 台!",
		park, total, success, int32(total)-success-notTry, notTry)
	return nil
}

func pause() {
	cmd := exec.Command("cmd", "/c", "pause")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func backUp(wg *sync.WaitGroup, queries ...logic.QueryCollector) {
	wg.Add(len(queries))
	errCh := make(chan error, len(queries))
	for _, query := range queries {
		q := query
		go func(errCh chan error) {
			defer wg.Done()
			cs, err := q.GetCollector()
			if err != nil {
				errCh <- fmt.Errorf("从%s获取数据信息错误:%v", q.GetTag(), err)
				return
			}
			errCh <- backgroundWork(cs, q.GetTag())
		}(errCh)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			logrus.Errorf("备份任务出错:%v", err)
		}
	}
}

func backgroundWork(cs []*models.Collector, park string) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second*time.Duration(conf.GlobalCfg.Fetch.Timeout))
	defer cancel()
	err := HandleServerList(ctx, cs, park)
	if err != nil {
		if err = utils.SendEmail("KC521B备份错误", fmt.Sprintf("备份%s错误:%v", park, err)); err != nil {
			logrus.Errorf("send email error:%v", err)
		}
		return err
	}
	return nil
}
