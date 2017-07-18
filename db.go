// Package mogo provides a faster usage of mongo with mgo behind
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
	ErrorURI = errors.New("mogo: could not parse URI")
	// ErrorModelID passing a model for update
	// if could not find Id or ID with bson.ObjectID will raise
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

// Close db session
func (db *DB) Close() {
	db.session.Close()
}

// Collection return mgo collection from model
func (db *DB) Collection(model interface{}) *mgo.Collection {
	col := db.database.C(colName(model))
	indexes := loadIndex(model)
	for i := range indexes {
		col.EnsureIndex(indexes[i])
	}
	return col
}

// DropCollection drop a collection
func (db *DB) DropCollection(model interface{}) error {
	col := db.Collection(model)
	return col.DropCollection()
}

// func (db *DB) SetIndex(model interface{}, index mgo.Index) error {
// 	col := db.Collection(model)
// 	return col.EnsureIndex(index)
// }

// Where start generating query
func (db *DB) Where(q bson.M) *Query {
	query := new(Query)
	query.db = db
	query.q = append(query.q, q)
	return query
}

// Get a model with id
func (db *DB) Get(model, id interface{}) error {
	if val, ok := id.(string); ok {
		id = bson.ObjectIdHex(val)
	}
	col := db.Collection(model)
	return col.FindId(id).One(model)
}

// Create insert a Document to DB
func (db *DB) Create(model interface{}) error {
	col := db.Collection(model)
	setID(model)
	return col.Insert(model)
}

// Update a Document
func (db *DB) Update(model interface{}) error {
	col := db.Collection(model)
	id, err := getID(model)
	if err != nil {
		return err
	}
	query := bson.M{"_id": id}
	fieldsUpdate := parseBson(model)
	if err := col.Update(query, fieldsUpdate); err != nil {
		return err
	}
	return db.Get(model, id)
}

// --------------------- in package ---------------

func parseBson(model interface{}) bson.M {
	b, _ := bson.Marshal(model)
	var body bson.M
	bson.Unmarshal(b, &body)
	return body
}

func parseURI(uri string) (string, string, error) {
	var url, db string
	splited := strings.Split(uri, "/")
	if len(splited) < 2 {
		return url, db, ErrorURI
	}
	url, db = splited[0], splited[1]
	return url, db, nil
}

func setID(model interface{}) {
	m := structs.Map(model)
	var keyID string
	if _, ok := m["Id"]; ok {
		m["Id"] = bson.NewObjectId()
		keyID = "Id"
	}
	if _, ok := m["ID"]; ok {
		m["ID"] = bson.NewObjectId()
		keyID = "ID"
	}
	s := structs.New(model)
	field := s.Field(keyID)
	field.Set(m[keyID])

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
