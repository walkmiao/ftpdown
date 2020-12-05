package regex

import (
	log "github.com/sirupsen/logrus"
	"regexp"
)

func Filter(expr string,file string)bool{
	reg,err:=regexp.Compile(expr)
	if err!=nil{
		log.Errorf("regexp compile error:%v\n",err)
		return false
	}
	return reg.MatchString(file)
}

