package mogo

import (
	"fmt"
	"strings"

	"gopkg.in/mgo.v2"
)

type coller interface {
	CollectionName() string
}

type indexer interface {
	Meta() []mgo.Index
}

func colName(model interface{}) string {
	if c, ok := model.(coller); ok {
		return c.CollectionName()
	}
	tmp := fmt.Sprintf("%T", model)
	tmp = strings.Replace(tmp, "*", "", -1)
	tmp = strings.Replace(tmp, "]", "", -1)
	tmp = strings.Replace(tmp, "[", "", -1)
	ts := strings.Split(tmp, ".")
	if len(ts) < 2 {
		return toSnake(tmp)
	}
	return toSnake(ts[1])
}

func loadIndex(model interface{}) []mgo.Index {
	if c, ok := model.(indexer); ok {
		return c.Meta()
	}
	return []mgo.Index{}
}
