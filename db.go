package mogo

import (
	"errors"
	"strings"

	"github.com/fatih/structs"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	// ErrorURI raise for parsing uri error
	ErrorURI     = errors.New("mogo: could not parse URI")
	ErrorModelID = errors.New("mogo: model id type error")
)

// DB main connection struct
type DB struct {
	session  *mgo.Session
	database *mgo.Database
}

// Conn init db struct with uri -> host:port/db
func Conn(uri string) (*DB, error) {
	url, db, err := parseURI(uri)
	if err != nil {
		return nil, err
	}
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	database := new(DB)
	database.database = session.DB(db)
	session.SetSafe(&mgo.Safe{})
	session.SetMode(mgo.Monotonic, true)
	database.session = session
	return database, nil
}

// Collection return mgo collection from model
func (db *DB) Collection(model interface{}) *mgo.Collection {
	return db.database.C(colName(model))
}

// DropCollection drop a collection
func (db *DB) DropCollection(model interface{}) error {
	col := db.Collection(model)
	return col.DropCollection()
}

func (db *DB) SetIndex(model interface{}, index mgo.Index) error {
	col := db.Collection(model)
	return col.EnsureIndex(index)
}

func (db *DB) Where(q bson.M) *Query {
	query := new(Query)
	query.db = db
	query.q = append(query.q, q)
	return query
}

func (db *DB) Get(model, id interface{}) error {
	col := db.Collection(model)
	return col.FindId(id).One(model)
}

func (db *DB) Create(model interface{}) error {
	col := db.Collection(model)
	return col.Insert(model)
}

func (db *DB) Update(model interface{}, fieldsUpdate bson.M) error {
	col := db.Collection(model)
	id, err := getID(model)
	if err != nil {
		return err
	}
	query := bson.M{"_id": id}
	return col.Update(query, fieldsUpdate)
}

// --------------------- in package ---------------
func parseURI(uri string) (string, string, error) {
	var url, db string
	splited := strings.Split(uri, "/")
	if len(splited) < 2 {
		return url, db, ErrorURI
	}
	url, db = splited[0], splited[1]
	return url, db, nil
}

func getID(model interface{}) (bson.ObjectId, error) {
	m := structs.Map(model)
	var (
		idInterface interface{}
		id          bson.ObjectId
		ok          bool
	)
	if val, ok := m["Id"]; ok {
		idInterface = val
	}
	if val, ok := m["ID"]; ok {
		idInterface = val
	}
	id, ok = idInterface.(bson.ObjectId)
	if !ok {
		return id, ErrorModelID
	}
	return id, nil
}
