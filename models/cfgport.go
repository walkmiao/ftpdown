package models

import "database/sql"

type Cfgport struct {
	PORTID        int64          `gorm:"primaryKey"`
	MONITORUNITID int64          `gorm:"MONITORUNITID"`
	PORTTYPE      int64          `gorm:"PORTTYPE"`
	PORTNO        sql.NullString `gorm:"PORTNO"`
	PORTSETTING   sql.NullString `gorm:"PORTSETTING"`
	PORTLIBNAME   sql.NullString `gorm:"PORTLIBNAME"`
	DESCRIPTION   sql.NullString `gorm:"DESCRIPTION"`
	EXTENDFIELD1  sql.NullString `gorm:"EXTENDFIELD1"`
	EXTENDFIELD2  sql.NullString `gorm:"EXTENDFIELD2"`
	EXTENDFIELD3  sql.NullString `gorm:"EXTENDFIELD3"`
	EXTENDFIELD4  sql.NullString `gorm:"EXTENDFIELD4"`
	EXTENDFIELD5  sql.NullString `gorm:"EXTENDFIELD5"`
}
