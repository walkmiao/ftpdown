package utils

import (
	"backup/conf"
	"backup/regex"
	"errors"
	"fmt"
	"github.com/jlaffaye/ftp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
	"time"
)

type ServerInfo struct {
	Addr        string
	TargetPath  string //存储目标目录
	PreviousDir string
	FileCount   int
	FetchConf
	*ftp.ServerConn
	*log.Entry
}

type FetchConf struct {
	port                  string
	username, password    string
	ServerPath, StorePath string
	Filters               string
	timeout               time.Duration
	JDServer, WGQServer   []string
}

var (
	DefaultFetchConfig FetchConf
)

func init() {
	if err := conf.InitConfig(""); err != nil {
		log.Errorf("init config error:%v\n", err)
		return
	}

	DefaultFetchConfig.port = viper.GetString("Server.Port")
	DefaultFetchConfig.timeout = time.Second * time.Duration(viper.GetInt("Server.Timeout"))
	DefaultFetchConfig.username, DefaultFetchConfig.password = viper.GetString("Account.UserName"),
		viper.GetString("Account.Password")
	DefaultFetchConfig.JDServer, DefaultFetchConfig.WGQServer =
		viper.GetStringSlice("Server.List.JD"), viper.GetStringSlice("Server.List.WGQ")

	DefaultFetchConfig.ServerPath, DefaultFetchConfig.StorePath =
		viper.GetString("Fetch.ServerPath"), viper.GetString("Fetch.StorePath")
	DefaultFetchConfig.Filters = viper.GetString("Fetch.Filters")
}

func NewServerInfo(addr string, region string) *ServerInfo {
	return &ServerInfo{Addr: addr, Entry: log.WithFields(log.Fields{
		"server": addr,
		"region": region,
		"regexp": DefaultFetchConfig.Filters,
	})}
}

func (s *ServerInfo) ConnectServer() error {
	conn, err := ftp.Dial(fmt.Sprintf("%s:%s", s.Addr, s.port),
		ftp.DialWithTimeout(s.timeout))

	if err != nil {
		return err
	}

	err = conn.Login(s.username, s.password)
	if err != nil {
		return fmt.Errorf("login error:%v\n", err)
	}
	s.ServerConn = conn
	s.Infof("connect to server  success")
	return nil
}

func (s *ServerInfo) WalkAndBuild() error {
	s.Info("Start download  from server")
	walker := s.Walk(s.ServerPath)
	if err := s.HandleWalker(walker); err != nil {
		return fmt.Errorf("handle walker error:%v\n", err)
	}
	s.Infof("success download from server! file count:%d\n", s.FileCount)
	return s.ServerConn.Quit()
}

func (s *ServerInfo) HandleEntry(entry *ftp.Entry, curPath string) error {
	//s.Warningf("handle entry %s!\n",entry.Name)
	switch t := entry.Type; t {
	//这个case还有待处理
	case ftp.EntryTypeFolder:
		willMkdir := path.Join(s.TargetPath, curPath)
		if err := os.MkdirAll(willMkdir, os.ModePerm); err != nil {
			return err
		}
		walker := s.Walk(curPath)
		if err := s.HandleWalker(walker); err != nil {
			return err
		}
	case ftp.EntryTypeFile:
		s.FileCount++
		if regex.FilterString(s.FetchConf.Filters, entry.Name) {
			s.Warningf("file %s has been filtered!", entry.Name)
			break
		}
		filePath := path.Join(s.TargetPath, curPath)
		//当大小为0时 response close的时候有bug需要跳过
		if entry.Size == 0 {
			fs, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				return err
			}
			s.Debugf("create file %s end,size 0\n", filePath)
			return fs.Close()
		}

		resp, err := s.Retr(curPath)
		if err != nil {
			return err
		}
		if err = s.WriteRespToFile(resp, filePath); err != nil {
			return err
		}
		fmt.Printf("下载文件%s success\n", entry.Name)
	default:
		s.Errorf("Type %s is not supported,name:%s\n", t.String(), entry.Name)
	}

	return nil
}

func (s *ServerInfo) HandleWalker(walker *ftp.Walker) error {
	if walker == nil {
		return errors.New("walker is nil")
	}
	s.Info("start walk path from server")
	if err := s.Mkdir(s.TargetPath); err != nil {
		return err
	}
	for walker.Next() {
		entry := walker.Stat()
		//if strings.Contains(entry.Name, "test11") {
		//	fmt.Println(entry.Name)
		//}
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
		return fmt.Errorf("connect to server error:%v\n", err)
	}
	return s.WalkAndBuild()
}
