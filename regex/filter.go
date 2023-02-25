package regex

import (
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
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
	if expr == "" {
		return false
	}
	checks := strings.Split(expr, "|")
	for _, check := range checks {
		if strings.Contains(file, check) {
			return true
		}
	}
	return false
}
