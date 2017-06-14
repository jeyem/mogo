package mogo

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

var (
	TestURI = "127.0.0.1:27017/test"
)

type EmbededTest struct {
	Name string `bson:"name"`
}

type TestCollection struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	TestField1 string        `bson:"test_field_1"`
	TestField2 int           `bson:"test_field_2"`
	Names      []EmbededTest `bson:"names"`
}

func (TestCollection) CollectionName() string {
	return "test"
}

func TestConnection(t *testing.T) {
	db, err := Conn(TestURI)
	if err != nil {
		t.Error(err)
		return
	}
	db.Close()
}

func TestCreateCollection(t *testing.T) {
	db, err := Conn(TestURI)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	testData := new(TestCollection)
	testData.TestField1 = "test"
	testData.TestField2 = 7777777
	if err := db.Create(testData); err != nil {
		t.Error(err)
		return
	}
	getData := new(TestCollection)
	if err := db.Get(getData, testData.ID); err != nil {
		t.Error(err)
		return
	}
	if getData.ID != testData.ID ||
		testData.TestField1 != getData.TestField1 ||
		testData.TestField2 != getData.TestField2 {
		t.Error("not match data values")

	}
}

func TestUpdateCollection(t *testing.T) {
	db, err := Conn(TestURI)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	testData := new(TestCollection)
	defer db.DropCollection(testData)
	testData.TestField1 = "test2"
	testData.TestField2 = 88888
	if err := db.Create(testData); err != nil {
		t.Error(err)
		return
	}
	field1To := "test35"
	testData.TestField1 = field1To
	if err := db.Update(testData); err != nil {
		t.Error(err)
		return
	}
	newtestData := new(TestCollection)
	if err := db.Get(newtestData, testData.ID); err != nil {
		t.Error(err)
		return
	}
	if newtestData.TestField1 != field1To {
		t.Error("collection not updated")
		return
	}
}

func TestQuery(t *testing.T) {
	db, err := Conn(TestURI)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	testData := new(TestCollection)
	defer db.DropCollection(testData)
	testData.TestField1 = "test3"
	testData.TestField2 = 88888
	if err := db.Create(testData); err != nil {
		t.Error(err)
		return
	}
	testResult := new(TestCollection)
	if err := db.Where(bson.M{"test_field_1": testData.TestField1}).Find(
		testResult); err != nil {
		t.Error(err)
		return
	}
	if testResult.ID != testData.ID ||
		testData.TestField1 != testResult.TestField1 ||
		testData.TestField2 != testResult.TestField2 {
		t.Error("not match data values", testData.ID, ">>", testResult.ID)

	}
}

func TestEmbeddedQuery(t *testing.T) {
	db, err := Conn(TestURI)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	testData := new(TestCollection)
	// defer db.DropCollection(testData)
	testData.TestField1 = "test3"
	testData.TestField2 = 88888
	for i := 1; i < 3; i++ {
		name := fmt.Sprintf("test%d", i)
		testData.Names = append(testData.Names, EmbededTest{Name: name})
	}
	if err := db.Create(testData); err != nil {
		t.Error(err)
		return
	}
	result := new(TestCollection)
	if err := db.Where(bson.M{"names.name": "test1"}).
		Find(result); err != nil {
		t.Error(err)
	}
}
