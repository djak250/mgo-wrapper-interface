package main

import (
	"fmt"
	"github.com/djak250/mgo-wrapper-interface/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// func wrapMongoRequest(handler func(w http.ResponseWriter, r *http.Request, db *mgo.Database)) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		db := mongo.Connection(r)
// 		handler(w, r, db)
// 	}
// }

func main() {
	session, merr := mgo.Dial("127.0.0.1:27017")
	defer session.Close()
	if merr != nil {
		fmt.Println(merr.Error())
		return
	}
	db := session.DB("test")
	mdb := models.MongoDatabase{db}

	testObj := map[string]interface{}{
		"_id":   bson.NewObjectId(),
		"test1": "1",
		"test2": "2",
		"test3": "3",
	}

	runDbSuite(mdb, testObj)
}

func runDbSuite(rDb models.MgoDatabase, testObj map[string]interface{}) {
	err := rDb.C("testCol").Insert(testObj)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Inserted:", false)
		return
	} else {
		fmt.Println("Inserted:", true)
	}

	testFindOneObj := map[string]interface{}{}
	testFindAllArray := make([]map[string]interface{}, 0)
	testUpdateObj := map[string]interface{}{}

	q := rDb.C("testCol").Find(bson.M{"_id": testObj["_id"]})

	err = q.One(&testFindOneObj)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Found: ", false)
		return
	}

	if _, ok := testFindOneObj["_id"]; ok && testFindOneObj["_id"].(bson.ObjectId) == testObj["_id"] {
		fmt.Println("Found: ", true)

	} else {
		fmt.Println("Found: ", false)
	}

	err = q.All(&testFindAllArray)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("FoundAll: ", false)
		return
	}
	if testFindAllArray[0] != nil {
		if _, ok2 := testFindAllArray[0]["_id"]; ok2 && testFindOneObj["_id"].(bson.ObjectId) == testObj["_id"] {
			fmt.Println("FoundAll: ", true)
		}
	} else {
		fmt.Println("FoundAll: ", false)
	}

	err = rDb.C("testCol").Update(bson.M{"_id": testObj["_id"]}, bson.M{"$set": map[string]interface{}{"updated": true}})
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Updated: ", false)
		return
	}
	rDb.C("testCol").Find(bson.M{"_id": testObj["_id"]}).One(&testUpdateObj)
	if _, ok := testUpdateObj["updated"]; ok && testUpdateObj["updated"].(bool) == true {
		fmt.Println("Updated: ", true)
	}

	err = rDb.C("testCol").Remove(bson.M{"_id": testObj["_id"]})
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Removed: ", false)
		return
	}

	q = rDb.C("testCol").Find(bson.M{"_id": testObj["_id"]})
	err = q.One(nil)
	if err != nil {
		if err.Error() == "not found" {
			fmt.Println("Removed: ", true)
			return
		} else {
			fmt.Println(err.Error())
		}
	}
	fmt.Println("Removed: ", false)
	return
}
