package mogo

import (
	"reflect"

	"gopkg.in/mgo.v2/bson"
)

type Query struct {
	db                *DB
	q                 []bson.M
	selectCulomn      bson.M
	selectedOneColumn bool
	sort              string
}

func (q *Query) Select(s bson.M) *Query {
	q.selectCulomn = s
	q.selectedOneColumn = true
	return q
}

func (q *Query) Or(s bson.M) *Query {
	q.q = append(q.q, s)
	return q
}

func (q *Query) Sort(s string) *Query {
	q.sort = s
	return q
}

func (q *Query) Find(model interface{}) error {
	if len(q.q) < 2 {
		err := q.result(model, q.q[0])
		return err
	}
	query := bson.M{"$or": q.q}
	return q.result(model, query)
}

// ----------------------- in package -------------------

func (q *Query) result(model interface{}, query bson.M) error {
	col := q.db.Collection(model)
	mgoQuery := col.Find(query)
	if q.selectedOneColumn {
		mgoQuery = mgoQuery.Select(q.selectCulomn)
	}
	if q.sort != "" {
		mgoQuery = mgoQuery.Sort(q.sort)
	}
	if isSlice(model) {
		return mgoQuery.All(model)
	}
	return mgoQuery.One(model)
}

func isSlice(model interface{}) bool {
	s := reflect.ValueOf(model)
	return s.Kind() == reflect.Slice
}
