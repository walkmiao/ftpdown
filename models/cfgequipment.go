package models

import (
	"database/sql"
)

type Cfgequipment struct {
	EQUIPID         int64          `gorm:"primaryKey"`
	PORTID          int64          `gorm:"PORTID"`
	MONITORUNITID   int64          `gorm:"MONITORUNITID"`
	EQUIPTEMPLATEID int64          `gorm:"EQUIPTEMPLATEID"`
	EQUIPNAME       sql.NullString `gorm:"EQUIPNAME"`
	EXTPORTSETTING  sql.NullString `gorm:"EXTPORTSETTING"`
	EVENTLOCKED     int64          `gorm:"EVENTLOCKED"`
	CONTROLLOCKED   int64          `gorm:"CONTROLLOCKED"`
	EXTENDFIELD1    sql.NullString `gorm:"EXTENDFIELD1"`
	EXTENDFIELD2    sql.NullString `gorm:"EXTENDFIELD2"`
	EXTENDFIELD3    sql.NullString `gorm:"EXTENDFIELD3"`
	EXTENDFIELD4    sql.NullString `gorm:"EXTENDFIELD4"`
	EXTENDFIELD5    sql.NullString `gorm:"EXTENDFIELD5"`
	ROOMID          int64          `gorm:"ROOMID"`
	EQUIPTYPE       int64          `gorm:"EQUIPTYPE"`
	SAMPLEINTERVAL  int64          `gorm:"SAMPLEINTERVAL"`
	LIBNAME         sql.NullString `gorm:"LIBNAME"`
	DESCRIPTION     sql.NullString `gorm:"DESCRIPTION"`
	EQUIPADDRESS    int64          `gorm:"EQUIPADDRESS"`
}

func (Cfgequipment) TableName() string {
	return "cfgequipment"
}

type Collector struct {
	EquipName string `gorm:"column:equipname"`
	Addr      string `gorm:"column:portsetting"`
}
