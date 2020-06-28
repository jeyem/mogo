// Package mogo provides a faster usage of mongo with mgo behind
package mogo

import (
	"errors"

	"github.com/fatih/structs"
	"github.com/globalsign/mgo"

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
	Session  *mgo.Session
	Database *mgo.Database
}

// Conn init db struct with uri -> host:port/db
func Conn(info *mgo.DialInfo) (*DB, error) {
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		return nil, err
	}
	database := new(DB)
	database.Database = session.DB(info.Database)
	session.SetSafe(&mgo.Safe{})
	session.SetMode(mgo.Monotonic, true)
	database.Session = session
	return database, nil
}

func ConnByURI(url string) (*DB, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	database := new(DB)
	database.Database = session.DB(info.Database)
	session.SetSafe(&mgo.Safe{})
	session.SetMode(mgo.Monotonic, true)
	database.Session = session
	return database, nil
}

// Close db session
func (db *DB) Close() {
	db.Session.Close()
}

// Collection return mgo collection from model
func (db *DB) Collection(model interface{}) *mgo.Collection {
	return db.Database.C(colName(model))
}

func (db *DB) Stream(model, query interface{}) *mgo.Iter {
	return db.Collection(model).Find(query).Iter()
}

// LoadIndexes reinitialize models indexes
func (db *DB) LoadIndexes(models ...interface{}) {
	for _, model := range models {
		col := db.Collection(model)
		indexes := loadIndex(model)
		for i := range indexes {
			col.EnsureIndex(indexes[i])
		}
	}
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

// Find start generating query
func (db *DB) Find(q bson.M) *Query {
	d := new(DB)
	d.Session = db.Session.Clone()
	d.Database = db.Database
	query := new(Query)
	query.db = d
	query.q = append(query.q, q)
	// fmt.Println(query)
	return query
}

// Get a model with id
func (db *DB) FindByID(model, id interface{}) error {
	if val, ok := id.(string); ok {
		id = bson.ObjectIdHex(val)
	}
	col := db.Collection(model)
	return col.FindId(id).One(model)
}

// Create a document in DB
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
	fieldsUpdate, err := parseBson(model)
	if err != nil {
		return err
	}
	if err := col.Update(query, bson.M{"$set": fieldsUpdate}); err != nil {
		return err
	}
	return db.FindByID(model, id)
}

// --------------------- in package ---------------

func parseBson(model interface{}) (bson.M, error) {
	b, err := bson.Marshal(model)
	if err != nil {
		return bson.M{}, err
	}
	var body bson.M
	bson.Unmarshal(b, &body)
	return body, nil
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
	if !id.Valid() {
		return id, ErrorModelID
	}
	return id, nil
}
