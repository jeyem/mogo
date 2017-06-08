package mogo

import (
	"errors"
	"strings"

	"gopkg.in/mgo.v2"
)

var (
	// ErrorURI raise for parsing uri error
	ErrorURI = errors.New("mongo: could not parse URI")
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

// --------------------- inpackage ---------------
func parseURI(uri string) (string, string, error) {
	var url, db string
	splited := strings.Split(uri, "/")
	if len(splited) < 2 {
		return url, db, ErrorURI
	}
	url, db = splited[0], splited[1]
	return url, db, nil
}
