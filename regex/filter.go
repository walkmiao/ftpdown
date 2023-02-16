package regex

import (
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

func Filter(expr string, file string) bool {
	reg, err := regexp.Compile(expr)
	if err != nil {
		log.Errorf("regexp compile error:%v\n", err)
		return false
	}
	return reg.MatchString(file)
}

func FilterString(expr string, file string) bool {
	checks := strings.Split(expr, "|")
	for _, check := range checks {
		if strings.Contains(file, check) {
			return true
		}
	}
	return false
}
