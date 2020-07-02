/**
* @Author:zhoutao
* @Date:2020/7/2 下午2:29
 */

package config

import (
	"regexp"
	"strings"
)

var AuthPermitConfig AuthPermitAll

type AuthPermitAll struct {
	PermitAll []interface{}
}

func Match(str string) bool {
	if len(AuthPermitConfig.PermitAll) > 0 {
		targetValue := AuthPermitConfig.PermitAll
		for i := 0; i < len(targetValue); i++ {
			s := targetValue[i].(string)
			res, _ := regexp.MatchString(strings.ReplaceAll(s, "**", "(.*?)"), str)
			if res {
				return true
			}
		}
	}
	return false
}
