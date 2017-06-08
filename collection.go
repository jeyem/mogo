package mogo

import (
	"fmt"
	"strings"
)

type loader interface {
	CollectionName() string
}

func colName(model interface{}) string {
	if c, ok := model.(loader); ok {
		return c.CollectionName()
	}
	tmp := fmt.Sprintf("%T", model)
	tmp = strings.Replace(tmp, "*", "", -1)
	tmp = strings.Replace(tmp, "]", "", -1)
	tmp = strings.Replace(tmp, "[", "", -1)
	ts := strings.Split(tmp, ".")
	if len(ts) < 2 {
		return tmp
	}
	return ts[1]
}
