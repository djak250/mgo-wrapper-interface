package models

import (
	"errors"
	"fmt"
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

type TMongoQuery struct {
	OneResponse interface{}
	AllResponse interface{}
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
	tmq := TMongoQuery{oneResponse, allResponse}
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
	*(r).(*[]map[string]interface{}) = tmq.AllResponse.([]map[string]interface{})
	return nil
}

func (tmq TMongoQuery) One(r interface{}) error {
	if _, ok := tmq.OneResponse.(error); ok {
		return tmq.OneResponse.(error)
	}
	*(r).(*map[string]interface{}) = tmq.OneResponse.(map[string]interface{})
	return nil
}
