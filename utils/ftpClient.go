package utils

import (
	"fmt"
	"ftpHelper/conf"
	"github.com/jlaffaye/ftp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path"
	"time"
)

type ServerInfo struct {
	addr string
	*CommonInfo
	*ftp.ServerConn
	SinglePath string
}

type CommonInfo struct {
	port                string
	username, password  string
	fetchPath           string
	ReadPath, StorePath string
	timeout             time.Duration
	ServerList          []string
}

var (
	DefaultCommonInfo = &CommonInfo{}
)

func init() {
	if err := conf.InitConfig(""); err != nil {
		log.Errorf("init config error:%v\n", err)
		return
	}
	port := viper.GetString("Server.Port")
	timeout := viper.GetInt("Server.Timeout")

	username, password := viper.GetString("Account.UserName"),
		viper.GetString("Account.Password")
	rootPath := viper.GetString("Fetch.RootPath")
	serverList := viper.GetStringSlice("Server.List")
	DefaultCommonInfo.port = port
	DefaultCommonInfo.timeout = time.Second * time.Duration(timeout)
	DefaultCommonInfo.username = username
	DefaultCommonInfo.password = password
	DefaultCommonInfo.fetchPath = rootPath
	DefaultCommonInfo.ServerList = serverList
	DefaultCommonInfo.ReadPath, DefaultCommonInfo.StorePath =
		viper.GetString("Fetch.ReadPath"), viper.GetString("Fetch.StorePath")
}

func NewServerInfo(addr string) *ServerInfo {
	return &ServerInfo{addr: addr}
}

func (s *ServerInfo) ConnectServer() error {
	conn, err := ftp.Dial(fmt.Sprintf("%s:%s", s.addr, s.port),
		ftp.DialWithTimeout(s.timeout))

	if err != nil {
		return err
	}

	err = conn.Login(s.username, s.password)
	if err != nil {
		return err
	}
	s.ServerConn = conn
	log.Infof("connect to server [%s] success", s.addr)
	return nil
}

func (s *ServerInfo) WalkAndBuild() {
	log.Infof("[BEGIN]download work from server [%s]!\n", s.addr)
	walker := s.Walk(s.ReadPath)
	if err := Mkdir(s.SinglePath); err != nil {
		log.Errorf("[%s] Make Dir error:%v\n", s.addr, err)
		return
	}
	for walker.Next() {
		entry := walker.Stat()
		dstPath := path.Join(s.SinglePath, walker.Path())
		log.Debugf("dstPath:%s\n", dstPath)
		if entry.Type == ftp.EntryTypeFolder {
			if err := Mkdir(dstPath); err != nil {
				log.Errorf("Make Dir error:%v\n", err)
				return
			}
		} else {
			if entry.Type == ftp.EntryTypeFile {
				resp, err := s.Retr(walker.Path())
				if err != nil {
					log.Errorf("Retr target %s error:%v\n", entry.Target, err)
					return
				}
				FromRespToFile(resp, path.Join(s.SinglePath, entry.Name))
			}
		}
	}
	log.Infof("[END]download from server [%s] success!\n", s.addr)
}

func (s *ServerInfo) Work() {
	if err := s.ConnectServer(); err != nil {
		log.Errorf("connect to server [%s] error:%v\n", s.addr, err)
		return
	}
	s.WalkAndBuild()
}
