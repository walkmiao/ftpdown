package utils

import (
	"archive/zip"
	"backup/conf"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/sirupsen/logrus"
)

// Compress 压缩文件
func Compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	// 如果是目录调用CompressedDir
	if info.IsDir() {
		return CompressedDir(file, prefix, zw)
	}
	// 如果是文件调用CompressedFile
	return CompressedFile(file, prefix, zw)
}

func CompressedFile(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if filepath.Ext(file.Name()) == ".zip" {
		return nil
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = prefix + "/" + header.Name
	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	if _, err = io.Copy(writer, file); err != nil {
		return err
	}
	return nil
}

// CompressedDir
func CompressedDir(file *os.File, prefix string, zw *zip.Writer) error {
	info, _ := file.Stat()
	if info.Name() != LATEST {
		prefix = prefix + "/" + info.Name()
	}
	dirInfo, err := file.Readdir(-1)
	if err != nil {
		return err
	}
	for _, f := range dirInfo {
		f, err := os.Open(file.Name() + "/" + f.Name())
		if err != nil {
			return err
		}
		err = Compress(f, prefix, zw)
		if err != nil {
			return err
		}
	}
	return nil
}

// ZipFile 压缩directory中的文件到zipName文件
func ZipFile(zipName, directory string) error {
	f, err := os.Create(zipName)
	if err != nil {
		return err
	}
	defer func() {
		upperPath := filepath.Dir(zipName)
		if err = DeleteExpireZip(upperPath); err != nil {
			logrus.Errorf("删除备份出错:%v 路径:%s", err, upperPath)
		}
	}()
	zw := zip.NewWriter(f)
	defer zw.Close()
	sf, err := os.Open(directory)
	defer sf.Close()
	if err != nil {
		return err
	}
	return Compress(sf, "", zw)
}

// DeleteExpireZip 删除多余的zip备份
func DeleteExpireZip(dir string) error {
	newestPath := path.Join(dir, LATEST)
	if err := os.RemoveAll(newestPath); err != nil {
		logrus.Errorf("删除最新下载文件夹[%s]错误:%v", newestPath, err)
	}
	zips := make([]string, 0)
	if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && d.Name() != filepath.Base(dir) {
			return filepath.SkipDir
		}
		if filepath.Ext(d.Name()) == ".zip" {
			zips = append(zips, d.Name())
		}
		return nil
	}); err != nil {
		return err
	}
	retain := conf.GlobalCfg.Fetch.Retain
	if len(zips) <= retain {
		return nil
	}
	sort.Strings(zips)
	needDel := len(zips) - retain
	for _, zipName := range zips[:needDel] {
		if err := os.Remove(path.Join(dir, zipName)); err != nil {
			logrus.Debugf("删除备份文件%s失败", zipName)
		} else {
			logrus.Debugf("删除备份文件%s成功", zipName)
		}
	}
	return nil
}
