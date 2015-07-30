package models

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"testing"
)

type TMongoDatabase struct {
	Collections []*TMongoCollection
	T           *testing.T
}

type TMongoCollection struct {
	Name            string
	FindResponses   []TMongoQuery
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
type TMongoQuery struct {
	OneResponse interface{}
	AllResponse interface{}
	T           *testing.T
}

//Implements MgoDatabase
func (tmd TMongoDatabase) C(colName string) MgoCollection {
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

func (tmd *TMongoDatabase) ExpectFind(colName string, oneResponse interface{}, allResponse interface{}) {
	tmq := TMongoQuery{oneResponse, allResponse, tmd.T}
	for _, c := range tmd.Collections {
		if c.Name == colName {
			c.FindResponses = append(c.FindResponses, tmq)
			return
		}
	}

	c := TMongoCollection{Name: colName}
	c.FindResponses = []TMongoQuery{tmq}
	tmd.Collections = append(tmd.Collections, &c)
}
func (tmd *TMongoDatabase) ExpectUpdate(colName string, response error) {
	for _, c := range tmd.Collections {
		if c.Name == colName {
			c.UpdateResponses = append(c.UpdateResponses, response)
			return
		}
	}
	c := TMongoCollection{Name: colName}
	c.UpdateResponses = []error{response}
	tmd.Collections = append(tmd.Collections, &c)
}
func (tmd *TMongoDatabase) ExpectRemove(colName string, response error) {
	for _, c := range tmd.Collections {
		if c.Name == colName {
			c.RemoveResponses = append(c.RemoveResponses, response)
			return
		}
	}
	c := TMongoCollection{Name: colName}
	c.RemoveResponses = []error{response}
	tmd.Collections = append(tmd.Collections, &c)
}
func (tmd *TMongoDatabase) ExpectInsert(colName string, response error) {
	for _, c := range tmd.Collections {
		if c.Name == colName {
			c.InsertResponses = append(c.InsertResponses, response)
			return
		}
	}
	c := TMongoCollection{Name: colName}
	c.InsertResponses = []error{response}
	tmd.Collections = append(tmd.Collections, &c)
}

//Implements MgoCollection
func (tmc *TMongoCollection) Find(q interface{}) MgoQuery {
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
func (tmc *TMongoCollection) Update(q interface{}, u interface{}) error {
	//Unshift first response off and return it
	if len(tmc.UpdateResponses) > 0 {
		ur, url := tmc.UpdateResponses[0], tmc.UpdateResponses[1:]
		tmc.UpdateResponses = url
		return ur
	} else {
		return errors.New("No Response Defined Update")
	}
}
func (tmc *TMongoCollection) Remove(q interface{}) error {
	//Unshift first response off and return it
	if len(tmc.RemoveResponses) > 0 {
		rr, rrl := tmc.RemoveResponses[0], tmc.RemoveResponses[1:]
		tmc.RemoveResponses = rrl
		return rr
	} else {
		return errors.New("No Response Defined Remove")
	}
}
func (tmc *TMongoCollection) Insert(i interface{}) error {
	//Unshift first response off and return it
	if len(tmc.InsertResponses) > 0 {
		ir, irl := tmc.InsertResponses[0], tmc.InsertResponses[1:]
		tmc.InsertResponses = irl
		return ir
	} else {
		return errors.New("No Response Defined Insert")
	}
}

//Implements MgoQuery
func (tmq TMongoQuery) All(r interface{}) error {
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

func (tmq TMongoQuery) One(r interface{}) error {
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
