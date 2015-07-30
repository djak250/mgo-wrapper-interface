package mongo

type IMgoSession interface {
	DB(string) IMgoDatabase
}

type IMgoDatabase interface {
	C(string) IMgoCollection
}

type IMgoCollection interface {
	Find(interface{}) IMgoQuery
	Update(interface{}, interface{}) error
	Remove(interface{}) error
	Insert(interface{}) error
}

type IMgoQuery interface {
	All(interface{}) error
	One(interface{}) error
	Select(interface{}) IMgoQuery
}
