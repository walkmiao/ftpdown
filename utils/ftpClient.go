package utils

import (
	"errors"
	"fmt"
	"ftpHelper/conf"
	"ftpHelper/regex"
	"github.com/jlaffaye/ftp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
	"time"
)

type ServerInfo struct {
	Addr,TargetPath string
	PreviousDir string
	FileCount int
	*CommonInfo
	*ftp.ServerConn
	*log.Entry
}

type CommonInfo struct {
	port                  string
	username, password    string
	ServerPath, StorePath string
	Regexp                string
	IsFilter			  bool
	timeout               time.Duration
	JDServer,WGQServer           []string

}

var (
	DefaultCommonInfo = &CommonInfo{}
)

func init() {
	if err := conf.InitConfig(""); err != nil {
		log.Errorf("init config error:%v\n", err)
		return
	}

	DefaultCommonInfo.port = viper.GetString("Server.Port")
	DefaultCommonInfo.timeout = time.Second * time.Duration(viper.GetInt("Server.Timeout"))
	DefaultCommonInfo.username,DefaultCommonInfo.password=viper.GetString("Account.UserName"),
		viper.GetString("Account.Password")
	DefaultCommonInfo.JDServer,DefaultCommonInfo.WGQServer =
		viper.GetStringSlice("Server.List.JD"),viper.GetStringSlice("Server.List.WGQ")

	DefaultCommonInfo.ServerPath, DefaultCommonInfo.StorePath =
		viper.GetString("Fetch.ServerPath"), viper.GetString("Fetch.StorePath")
	DefaultCommonInfo.Regexp = viper.GetString("Fetch.Regexp")
	DefaultCommonInfo.IsFilter = DefaultCommonInfo.Regexp!=""
}

func NewServerInfo(addr string,region string) *ServerInfo {
	return &ServerInfo{Addr: addr,Entry:log.WithFields(log.Fields{
		"server":addr,
		"region":region,
		"isFilterOpen":DefaultCommonInfo.IsFilter,
		"regexp":DefaultCommonInfo.Regexp,

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
		return fmt.Errorf("login error:%v\n",err)
	}
	s.ServerConn = conn
	s.Infof("connect to server  success")
	return nil
}

func (s *ServerInfo) WalkAndBuild()error {
	s.Info("Start download  from server")
	walker:=s.Walk(s.ServerPath)
	if err:=s.HandleWalker(walker);err!=nil{
		return fmt.Errorf("handle walker error:%v\n",err)
	}
	s.Infof("success download from server! file count:%d\n",s.FileCount)
	return s.ServerConn.Quit()
}


func (s *ServerInfo)HandleEntry(entry *ftp.Entry,curPath string)error{
	//s.Warningf("handle entry %s!\n",entry.Name)
	switch t:=entry.Type;t {
	//case ftp.EntryTypeFolder:
	//	s.PreviousDir = path.Join(s.PreviousDir,entry.Name)
	//	if err:=s.Mkdir(s.PreviousDir);err!=nil{
	//		return err
	//	}

	case ftp.EntryTypeFile:
		s.FileCount++
		if s.IsFilter{
			if regex.Filter(s.CommonInfo.Regexp,entry.Name){
				s.Warningf("file %s has been filtered!",entry.Name)
				break
			}
		}

		//当大小为0时 response close的时候有bug需要跳过
		if entry.Size==0{
			filePath:=path.Join(s.TargetPath,entry.Name)
			fs,err:=os.OpenFile(filePath,os.O_CREATE|os.O_WRONLY,os.ModePerm)
			if err!=nil{
				return err
			}
			s.Debugf("create file %s end,size 0\n",filePath)
			return fs.Close()
		}

		resp,err:=s.Retr(curPath)
		if err!=nil{
			return err
		}
		if err=s.WriteRespToFile(resp,path.Join(s.TargetPath,entry.Name));err!=nil{
			return err
		}
	default:
		s.Errorf("Type %s is not supported,name:%s\n",t.String(),entry.Name)
	}

	return nil
}

func(s *ServerInfo)HandleWalker(walker *ftp.Walker)error{
		if walker==nil{
			return errors.New("walker is nil")
		}
		s.Info("start walk path from server")
		if err:=s.Mkdir(s.TargetPath);err!=nil{
			return err
		}
		for walker.Next(){
			entry:=walker.Stat()
			curPath:=walker.Path()
			err:=s.HandleEntry(entry,curPath)
			if err!=nil{
				return err
			}
		}
		return nil
}

func (s *ServerInfo) Work()error {
	if err := s.ConnectServer(); err != nil {
		return fmt.Errorf("connect to server error:%v\n", err)
	}
	return s.WalkAndBuild()
}


