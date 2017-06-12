package mogo

import (
	"reflect"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Query struct {
	db                *DB
	q                 []bson.M
	selectCulomn      bson.M
	selectedOneColumn bool
	sort              string
	querySet          *mgo.Query
	paginated         bool
	skip              int
	limit             int
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

func (q *Query) Paginate(limit, page int) *Query {
	q.paginated = true
	q.limit = limit
	q.skip = page * limit
	return q
}

func (q *Query) Limit(limit int) *Query {
	q.paginated = true
	q.limit = limit
	q.skip = 0
	return q
}

func (q *Query) Find(model interface{}) error {
	query := q.parseQuery()
	q.loadQuerySet(model, query)
	return q.result(model)
}

func (q *Query) Count(model interface{}) (int, error) {
	query := q.parseQuery()
	q.loadQuerySet(model, query)
	return q.querySet.Count()
}

// ----------------------- in package -------------------

func (q *Query) parseQuery() bson.M {
	if len(q.q) < 2 {
		return q.q[0]
	}
	return bson.M{"$or": q.q}
}

func (q *Query) loadQuerySet(model interface{}, query bson.M) {
	col := q.db.Collection(model)
	mgoQuery := col.Find(query)
	if q.selectedOneColumn {
		mgoQuery = mgoQuery.Select(q.selectCulomn)
	}
	if q.paginated {
		mgoQuery.Skip(q.skip).Limit(q.limit)
	}
	if q.sort != "" {
		mgoQuery = mgoQuery.Sort(q.sort)
	}
	q.querySet = mgoQuery
}

func (q *Query) result(model interface{}) error {
	if isSlice(model) {
		return q.querySet.All(model)
	}
	return q.querySet.One(model)
}

func isSlice(model interface{}) bool {
	s := reflect.ValueOf(model)
	return s.Kind() == reflect.Slice
}
