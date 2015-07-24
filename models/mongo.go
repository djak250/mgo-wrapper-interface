package models

import (
	"gopkg.in/mgo.v2"
)

type MongoDatabase struct {
	DB *mgo.Database
}

//Implements MgoDatabase
func (md MongoDatabase) C(c string) MgoCollection {
	return MongoCollection{md.DB.C(c)}
}

type MongoCollection struct {
	C *mgo.Collection
}

//Implements MgoCollection
func (mc MongoCollection) Find(q interface{}) MgoQuery {
	mq := mc.C.Find(q)
	mgq := MongoQuery{mq}
	return &mgq
}
func (mc MongoCollection) Update(q interface{}, u interface{}) error {
	return mc.C.Update(q, u)
}
func (mc MongoCollection) Remove(q interface{}) error {
	return mc.C.Remove(q)
}
func (mc MongoCollection) Insert(i interface{}) error {
	return mc.C.Insert(i)
}

type MongoQuery struct {
	Q *mgo.Query
}

//Implements MgoQuery
func (mq MongoQuery) All(r interface{}) error {
	return mq.Q.All(r)
}

func (mq MongoQuery) One(r interface{}) error {
	return mq.Q.One(r)
}
