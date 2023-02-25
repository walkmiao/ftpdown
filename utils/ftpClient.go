package utils

import (
	"backup/conf"
	"backup/models"
	"backup/regex"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
	log "github.com/sirupsen/logrus"
)

const (
	LATEST = "最新下载"
)

var ErrMap sync.Map

type Invalid struct {
	ErrCount      int64
	RetryInterval time.Duration
	Forbidden     bool
}

func (inv Invalid) RetryTime() time.Time {
	return time.Now().Add(inv.RetryInterval)
}

type ServerInfo struct {
	Addr string
	Park string
	Name string
	*FetchConf
	*ftp.ServerConn
	*log.Entry
}

func (info *ServerInfo) ZipName() string {
	zipName := fmt.Sprintf("%s(%s)-%s.zip",
		info.Name, info.Addr, time.Now().Format("2006-01-02 15:04:05"))
	return path.Join(filepath.Dir(info.TargetPath()), zipName)
}

// 存储目录
func (info *ServerInfo) TargetPath() string {
	return path.Join(conf.GlobalCfg.Fetch.StorePath, info.Park,
		fmt.Sprintf("%s(%s)", info.Name, info.Addr), LATEST)
}

type FetchConf struct {
	ServerPath, StorePath string
	Filters               string
	Timeout               time.Duration
}

var DefaultFetchConfig *FetchConf

func NewServerInfo(server *models.Collector, park string) *ServerInfo {
	if DefaultFetchConfig == nil {
		DefaultFetchConfig = &FetchConf{
			ServerPath: conf.GlobalCfg.Fetch.ServerPath,
			StorePath:  conf.GlobalCfg.Fetch.StorePath,
			Filters:    conf.GlobalCfg.Fetch.Filters,
			Timeout:    time.Duration(conf.GlobalCfg.Fetch.Timeout),
		}
	}
	return &ServerInfo{
		Addr:      server.Addr,
		Name:      server.EquipName,
		FetchConf: DefaultFetchConfig,
		Park:      park,
		Entry: log.WithFields(log.Fields{
			"server": server.Addr,
			"name":   server.EquipName,
			"park":   park,
		})}
}

func (s *ServerInfo) ConnectServer() error {
	conn, err := ftp.Dial(s.Addr,
		ftp.DialWithTimeout(s.Timeout*time.Second))
	if err == nil {
		ErrMap.Delete(s.Addr)
	} else {
		v, loaded := ErrMap.LoadOrStore(s.Addr, &Invalid{
			ErrCount:      1,
			RetryInterval: time.Hour * time.Duration(conf.GlobalCfg.Retry.RetryInterval),
		})
		//出错后，8小时内不允许下载
		if loaded {
			inv := v.(*Invalid)
			if cnt := conf.GlobalCfg.Retry.MaxFailed; inv.ErrCount >= cnt {
				inv.Forbidden = true
				s.Debugf("(%s)%s 已经超过最大尝试上限 %d次,已不再下载", s.Name, s.Addr, cnt)
				return nil
			}
			inv.ErrCount += 1
			inv.RetryInterval = time.Duration(float64(inv.RetryInterval) * (1 + conf.GlobalCfg.Retry.ThresholdFactor))
		}
		return err
	}

	var errLogin error
	for _, account := range conf.GlobalCfg.Accounts {
		if errLogin = conn.Login(account.UserName, account.Password); err != nil {
			continue
		}
		break
	}
	if errLogin != nil {
		return errLogin
	}
	s.ServerConn = conn
	s.Debugf("connect to server %s  success", s.Addr)
	return nil
}

func (s *ServerInfo) WalkAndBuild() error {
	walker := s.Walk(s.ServerPath)
	if err := s.HandleWalker(walker); err != nil {
		return fmt.Errorf("handle walker error:%v\n", err)
	}
	if err := ZipFile(s.ZipName(), s.TargetPath()); err != nil {
		return fmt.Errorf("zip file error:%v", err)
	}
	return s.ServerConn.Quit()
}

func (s *ServerInfo) HandleEntry(entry *ftp.Entry, curPath string) error {
	//s.Warningf("handle entry %s!\n",entry.Name)
	switch t := entry.Type; t {
	//这个case还有待处理
	case ftp.EntryTypeFolder:
		willMkdir := path.Join(s.TargetPath(), curPath)
		if err := os.MkdirAll(willMkdir, os.ModePerm); err != nil {
			return err
		}
		walker := s.Walk(curPath)
		if err := s.HandleWalker(walker); err != nil {
			return err
		}
	case ftp.EntryTypeFile:
		if regex.FilterString(s.FetchConf.Filters, entry.Name) {
			s.Warningf("file %s has been filtered!", entry.Name)
			break
		}
		filePath := path.Join(s.TargetPath(), curPath)
		//当大小为0时 response close的时候有bug需要跳过
		if entry.Size == 0 {
			fs, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				return err
			}
			return fs.Close()
		}

		resp, err := s.Retr(curPath)
		if err != nil {
			return err
		}
		if err = s.WriteRespToFile(resp, filePath); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServerInfo) HandleWalker(walker *ftp.Walker) error {
	if walker == nil {
		return errors.New("walker is nil")
	}
	if err := os.MkdirAll(s.TargetPath(), os.ModePerm); err != nil {
		return err
	}
	for walker.Next() {
		entry := walker.Stat()
		curPath := walker.Path()
		err := s.HandleEntry(entry, curPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ServerInfo) Work() error {
	if err := s.ConnectServer(); err != nil {
		return err
	}
	return s.WalkAndBuild()
}
