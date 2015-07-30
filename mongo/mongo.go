package mongo

import (
	"gopkg.in/mgo.v2"
)

type MgoSession struct {
	Session *mgo.Session
}

func (ms MgoSession) DB(db string) IMgoDatabase {
	return MgoDatabase{ms.Session.DB(db)}
}

type MgoDatabase struct {
	DB *mgo.Database
}

//Implements MgoDatabase
func (md MgoDatabase) C(c string) IMgoCollection {
	return MgoCollection{md.DB.C(c)}
}

type MgoCollection struct {
	C *mgo.Collection
}

//Implements IMgoCollection
func (mc MgoCollection) Find(q interface{}) IMgoQuery {
	mq := mc.C.Find(q)
	mgq := MgoQuery{mq}
	return &mgq
}
func (mc MgoCollection) Update(q interface{}, u interface{}) error {
	return mc.C.Update(q, u)
}
func (mc MgoCollection) Remove(q interface{}) error {
	return mc.C.Remove(q)
}
func (mc MgoCollection) Insert(i interface{}) error {
	return mc.C.Insert(i)
}

type MgoQuery struct {
	Q *mgo.Query
}

//Implements IMgoQuery
func (mq MgoQuery) All(r interface{}) error {
	return mq.Q.All(r)
}

func (mq MgoQuery) One(r interface{}) error {
	return mq.Q.One(r)
}
