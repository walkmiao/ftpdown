package conf

import (
	"fmt"
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var GlobalCfg Cfg

type Config struct {
	Name string
}

type Cfg struct {
	Logger   Logger
	Server   Server
	Mysql    Mysql
	Fetch    Fetch
	Retry    *Retry
	Accounts []Account
	Email    Email
}
type Retry struct {
	RetryInterval   int     `yaml:"retryInterval"`
	ThresholdFactor float64 `yaml:"thresholdFactor"`
	MaxFailed       int64   `yaml:"maxFailed"`
}
type Logger struct {
	Level    string `yaml:"level"`
	FileName string `yaml:"fileName"`
}
type Server struct {
	Port    int `yaml:"port"`
	Timeout int `yaml:"timeout"`
}
type MysqlConf struct {
	Addr     string `yaml:"addr"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
	Tag      string `yaml:"tag"`
}

type Mysql struct {
	Jd  *MysqlConf `yaml:"jd"`
	Wgq *MysqlConf `yaml:"wgq"`
}

func (m MysqlConf) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.User, m.Password, m.Addr, m.Db)
}

type Fetch struct {
	Timeout    int    `yaml:"timeout"`
	ServerPath string `yaml:"serverPath"`
	StorePath  string `yaml:"storePath"`
	Filters    string `yaml:"filters"`
	Factor     int    `yaml:"factor"`
	Interval   int    `yaml:"interval"`
	Retain     int    `yaml:"retain"`
}
type Email struct {
	Endpoint  string   `yaml:"endpoint"`
	Receivers []string `yaml:"receivers"`
}
type Account struct {
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
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
		viper.AddConfigPath("..")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&GlobalCfg); err != nil {
		return err
	}
	if GlobalCfg.Fetch.Interval == 0 {
		GlobalCfg.Fetch.Interval = 24
	}
	if GlobalCfg.Fetch.Retain == 0 {
		GlobalCfg.Fetch.Retain = 10
	}
	if GlobalCfg.Retry == nil {
		GlobalCfg.Retry = &Retry{
			RetryInterval:   8,
			ThresholdFactor: 0.25,
			MaxFailed:       10,
		}
	}

	return nil
}

func initLog() error {
	formatter := &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		ForceQuote:      true,
		ForceColors:     true,
	}

	if level := GlobalCfg.Logger.Level; level != "" {
		lvl, err := logrus.ParseLevel(level)
		if err != nil {
			return err
		}
		logrus.SetLevel(lvl)
	}
	if logrus.GetLevel() == logrus.DebugLevel {
		logrus.SetReportCaller(true)
	}
	writer, err := rotatelogs.New(
		GlobalCfg.Logger.FileName+".%Y%m%d%H%M",
		rotatelogs.WithMaxAge(24*7*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		return err
	}
	logrus.SetOutput(io.MultiWriter(os.Stderr, writer))
	logrus.SetFormatter(formatter)
	return nil
}
