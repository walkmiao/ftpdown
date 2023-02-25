package logic

import (
	"backup/conf"
	"backup/models"
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	JDDB  *gorm.DB
	WGQDB *gorm.DB
)

type QueryCollector interface {
	GetCollector() ([]*models.Collector, error)
	GetTag() string
}

type MysqlQuery struct {
	MysqlConf *conf.MysqlConf
}

func (q *MysqlQuery) GetTag() string {
	switch q.MysqlConf.Tag {
	case "jd":
		return "嘉定园区"
	case "wgq":
		return "外高桥园区"
	default:
		return "Mixed"
	}
}

var _ QueryCollector = (*MysqlQuery)(nil)

func (q *MysqlQuery) initDB(db *gorm.DB) error {
	if db == nil {
		dbOpen, err := gorm.Open(mysql.Open(q.MysqlConf.Dsn()), &gorm.Config{})
		if err != nil {
			return err
		}
		db = dbOpen
	}
	return nil
}

func (q *MysqlQuery) GetCollector() ([]*models.Collector, error) {
	switch q.MysqlConf.Tag {
	case "jd":
		if err := q.initDB(JDDB); err != nil {
			return nil, err
		}
		var collectos []*models.Collector
		if err := JDDB.Table("cfgequipment").Select("equipname,cfgport.portsetting").
			Joins("left join cfgport on cfgequipment.portid=cfgport.portid").Where("LIBNAME LIKE ?", "%IEC104%").
			Scan(&collectos).Error; err != nil {
			return nil, err
		}
		for _, c := range collectos {
			c.Addr = strings.Replace(c.Addr, ":2404", ":21", 1)
		}
		return collectos, nil
	case "wgq":
		if err := q.initDB(WGQDB); err != nil {
			return nil, err
		}
		var collectos []*models.Collector
		if err := WGQDB.Table("cfgequipment").Select("equipname,cfgport.portsetting").
			Joins("left join cfgport on cfgequipment.portid=cfgport.portid").Where("LIBNAME LIKE ?", "%IEC104%").
			Scan(&collectos).Error; err != nil {
			return nil, err
		}
		for _, c := range collectos {
			c.Addr = strings.Replace(c.Addr, ":2404", ":21", 1)
		}
		return collectos, nil
	default:
		return nil, fmt.Errorf("tag %s not supported", q.MysqlConf.Tag)
	}

}
