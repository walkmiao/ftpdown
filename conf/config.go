package conf

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Name string
}

func InitConfig(name string) error {
	c := &Config{name}
	if err := c.initConfig(); err != nil {
		return err
	}
	if err := initLog(); err != nil {
		return err
	}
	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func initLog() error {
	formatter := &logrus.TextFormatter{
		TimestampFormat:"2006-01-02 15:04:05",
		ForceQuote:true,
	}

	if level := viper.GetString("Logger.Level"); level != "" {
		lvl, err := logrus.ParseLevel(level)
		if err != nil {
			return err
		}
		logrus.SetLevel(lvl)
	}
	if name := viper.GetString("Logger.FileName"); name != "" {
		fs, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		logrus.SetOutput(fs)
	} else {
		formatter.ForceColors =true

	}
	logrus.SetFormatter(formatter)
	return nil
}
