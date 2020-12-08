package utils

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"os"
)

func(s *ServerInfo) Mkdir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		s.Infof("dir %s is not exist,make dir!\n", dir)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
		s.Infof("make dir %s success!\n",dir)
	}
	return nil
}

func (s *ServerInfo)WriteRespToFile(resp *ftp.Response, path string)error{
	fs, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)

	if err != nil {
		return err
	}
	defer fs.Close()
	_, err = io.Copy(fs, resp)
	if err != nil {
		return err
	}

	if err:=resp.Close();err!=nil{
		return fmt.Errorf("resp close error:%v\n",err)
	}
	s.Debugf("write file %s end!\n",path)
	return nil
}
