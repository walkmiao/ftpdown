package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/jlaffaye/ftp"
)

func (s *ServerInfo) WriteRespToFile(resp *ftp.Response, path string) error {
	fs, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer fs.Close()
	_, err = io.Copy(fs, resp)
	if err != nil {
		return err
	}

	if err = resp.Close(); err != nil {
		return fmt.Errorf("resp close error:%v\n", err)
	}
	return nil
}
