package mongo

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"testing"
)

type TMgoSession struct {
	Databases []*TMgoDatabase
	T         *testing.T
}

type TMgoDatabase struct {
	Name        string
	Collections []*TMgoCollection
	T           *testing.T
}

type TMgoCollection struct {
	Name            string
	FindResponses   []TMgoQuery
	UpdateResponses []error
	RemoveResponses []error
	InsertResponses []error
	T               *testing.T
}

/*
  OneResponse: error / bson.M
  AllResponse: error / []bson.M

  The usage of bson.M allows us to utilize bson.Unmarshaling log to handle
  weird pointer/interface reflection for us.
*/
type TMgoQuery struct {
	OneResponse interface{}
	AllResponse interface{}
	T           *testing.T
}

func (tms TMgoSession) DB(dbName string) IMgoDatabase {
	for _, db := range tms.Databases {
		if db.Name == dbName {
			db.T = tms.T
			return db
		}
	}
	fmt.Println("Database '" + dbName + "' not defined")
	tms.T.FailNow()
	return nil
}

//Implements IMgoDatabase
func (tmd TMgoDatabase) C(colName string) IMgoCollection {
	for _, c := range tmd.Collections {
		if c.Name == colName {
			c.T = tmd.T
			return c
		}
	}
	fmt.Println("Collection '" + colName + "' not defined")
	tmd.T.FailNow()
	return nil
}

func (tmd *TMgoDatabase) ExpectFind(colName string, oneResponse interface{}, allResponse interface{}) {
	tmq := TMgoQuery{oneResponse, allResponse, tmd.T}
	for _, c := range tmd.Collections {
		if c.Name == colName {
			c.FindResponses = append(c.FindResponses, tmq)
			return
		}
	}

	c := TMgoCollection{Name: colName}
	c.FindResponses = []TMgoQuery{tmq}
	tmd.Collections = append(tmd.Collections, &c)
}
func (tmd *TMgoDatabase) ExpectUpdate(colName string, response error) {
	for _, c := range tmd.Collections {
		if c.Name == colName {
			c.UpdateResponses = append(c.UpdateResponses, response)
			return
		}
	}
	c := TMgoCollection{Name: colName}
	c.UpdateResponses = []error{response}
	tmd.Collections = append(tmd.Collections, &c)
}
func (tmd *TMgoDatabase) ExpectRemove(colName string, response error) {
	for _, c := range tmd.Collections {
		if c.Name == colName {
			c.RemoveResponses = append(c.RemoveResponses, response)
			return
		}
	}
	c := TMgoCollection{Name: colName}
	c.RemoveResponses = []error{response}
	tmd.Collections = append(tmd.Collections, &c)
}
func (tmd *TMgoDatabase) ExpectInsert(colName string, response error) {
	for _, c := range tmd.Collections {
		if c.Name == colName {
			c.InsertResponses = append(c.InsertResponses, response)
			return
		}
	}
	c := TMgoCollection{Name: colName}
	c.InsertResponses = []error{response}
	tmd.Collections = append(tmd.Collections, &c)
}

//Implements IMgoCollection
func (tmc *TMgoCollection) Find(q interface{}) IMgoQuery {
	//Unshift first response off and return it
	if len(tmc.FindResponses) > 0 {
		fr, frl := tmc.FindResponses[0], tmc.FindResponses[1:]
		tmc.FindResponses = frl
		return &fr
	} else {
		fmt.Println("No Response Defined Find")
		return nil
	}
}
func (tmc *TMgoCollection) Update(q interface{}, u interface{}) error {
	//Unshift first response off and return it
	if len(tmc.UpdateResponses) > 0 {
		ur, url := tmc.UpdateResponses[0], tmc.UpdateResponses[1:]
		tmc.UpdateResponses = url
		return ur
	} else {
		return errors.New("No Response Defined Update")
	}
}
func (tmc *TMgoCollection) Remove(q interface{}) error {
	//Unshift first response off and return it
	if len(tmc.RemoveResponses) > 0 {
		rr, rrl := tmc.RemoveResponses[0], tmc.RemoveResponses[1:]
		tmc.RemoveResponses = rrl
		return rr
	} else {
		return errors.New("No Response Defined Remove")
	}
}
func (tmc *TMgoCollection) Insert(i interface{}) error {
	//Unshift first response off and return it
	if len(tmc.InsertResponses) > 0 {
		ir, irl := tmc.InsertResponses[0], tmc.InsertResponses[1:]
		tmc.InsertResponses = irl
		return ir
	} else {
		return errors.New("No Response Defined Insert")
	}
}

//Implements IMgoQuery
func (tmq TMgoQuery) All(r interface{}) error {
	if _, ok := tmq.AllResponse.(error); ok {
		return tmq.AllResponse.(error)
	}
	rv := reflect.ValueOf(r)
	slicev := rv.Elem()
	slicev = slicev.Slice(0, slicev.Cap())
	elemt := slicev.Type().Elem()
	i := 0
	for {
		if slicev.Len() == i {
			if i > len(tmq.AllResponse.([]bson.M))-1 {
				break
			}
			elemp := reflect.New(elemt)
			_bytes, err := bson.Marshal(tmq.AllResponse.([]bson.M)[i])
			err = bson.Unmarshal(_bytes, elemp.Interface())
			if err != nil {
				tmq.T.Fatal("Invalid response type / expectation type")
			}
			slicev = reflect.Append(slicev, elemp.Elem())
			slicev = slicev.Slice(0, slicev.Cap())
		} else {
			if i > len(tmq.AllResponse.([]bson.M))-1 {
				break
			}
			_bytes, err := bson.Marshal(tmq.AllResponse.([]bson.M)[i])
			err = bson.Unmarshal(_bytes, slicev.Index(i).Addr().Interface())
			if err != nil {
				tmq.T.Fatal("Invalid response type / expectation type")
			}

		}
		i++
	}
	rv.Elem().Set(slicev.Slice(0, i))
	return nil
}

func (tmq TMgoQuery) One(r interface{}) error {
	if _, ok := tmq.OneResponse.(error); ok {
		return tmq.OneResponse.(error)
	}
	_bytes, err := bson.Marshal(tmq.OneResponse.(bson.M))
	err = bson.Unmarshal(_bytes, r)
	if err != nil {
		tmq.T.Fatal("Invalid response type / expectation type")
	}
	return nil
}

func (tmq TMgoQuery) Select(s interface{}) IMgoQuery {
	//TODO:
	//THIS IS BAD BUT I HAVEN'T DONE ANYTHING TO MOCK UP SELECT YET!!!
	return tmq
}
