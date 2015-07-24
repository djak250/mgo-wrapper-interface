package models

type MgoDatabase interface {
	C(string) MgoCollection
}

type MgoCollection interface {
	Find(interface{}) MgoQuery
	Update(interface{}, interface{}) error
	Remove(interface{}) error
	Insert(interface{}) error
}

type MgoQuery interface {
	All(interface{}) error
	One(interface{}) error
}
