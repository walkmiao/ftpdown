package logic

import (
	"backup/conf"
	"backup/models"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type QueryCollector interface {
	GetCollector() ([]*models.Collector, error)
	GetTag() string
}

type MysqlQuery struct {
	MysqlConf *conf.MysqlConf
	*gorm.DB
}

func (q *MysqlQuery) GetTag() string {
	if q.MysqlConf.Tag == "" {
		return "默认标签"
	}
	return q.MysqlConf.Tag
}

var _ QueryCollector = (*MysqlQuery)(nil)

func (q *MysqlQuery) initDB() error {
	if q.DB == nil {
		dbOpen, err := gorm.Open(mysql.Open(q.MysqlConf.Dsn()), &gorm.Config{})
		if err != nil {
			return err
		}
		q.DB = dbOpen
	}
	return nil
}

func (q *MysqlQuery) GetCollector() ([]*models.Collector, error) {
	if err := q.initDB(); err != nil {
		return nil, err
	}
	var collectos []*models.Collector
	if err := q.Table("cfgequipment").Select("equipname,cfgport.portsetting").
		Joins("left join cfgport on cfgequipment.portid=cfgport.portid").Where("LIBNAME LIKE ?", "%IEC104%").
		Scan(&collectos).Error; err != nil {
		return nil, err
	}
	for _, c := range collectos {
		c.Addr = strings.Replace(c.Addr, ":2404", ":21", 1)
	}
	return collectos, nil

}
