package utils

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func Mkdir(name string) error {
	if _, err := os.Stat(name); err != nil {
		if err := os.MkdirAll(name, os.ModePerm); err != nil {
			fmt.Printf("make dir %s fail:%v\n", name, err)
			return err
		}
	}
	fmt.Printf("make dir %s success!\n", name)
	return nil
}

func FromRespToFile(resp *ftp.Response, dstPath string) {
	defer resp.Close()
	fs, err := os.OpenFile(dstPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Errorf("open file error:%v\n", err)
		return
	}
	defer fs.Close()
	_, err = io.Copy(fs, resp)
	if err != nil {
		log.Errorf("io copy error:%v\n", err)
		return
	}
}
