package main

import (
	"fmt"
	"github.com/djak250/mgo-wrapper-interface/mongo"
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
	rS := mongo.MgoSession{session}
	defer rS.Session.Close()
	if merr != nil {
		fmt.Println(merr.Error())
		return
	}

	testObj := map[string]interface{}{
		"_id":   bson.NewObjectId(),
		"test1": "1",
		"test2": "2",
		"test3": "3",
	}

	runDbSuite(rS, testObj)
}

func runDbSuite(rs mongo.IMgoSession, testObj map[string]interface{}) {
	db := rs.DB("test")

	err := db.C("testCol").Insert(testObj)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Inserted:", false)
		return
	} else {
		fmt.Println("Inserted:", true)
	}

	testFindOneObj := map[string]interface{}{}
	testSelectObj := map[string]interface{}{}
	testFindAllArray := make([]map[string]interface{}, 0)
	testUpdateObj := map[string]interface{}{}

	q := db.C("testCol").Find(bson.M{"_id": testObj["_id"]})

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

	q.Select(bson.M{"test1": false}).One(&testSelectObj)
	if _, ok := testSelectObj["test1"]; !ok && testSelectObj["_id"].(bson.ObjectId) == testObj["_id"] {
		fmt.Println("Selected: ", true)

	} else {
		fmt.Println("Selected: ", false)
	}

	err = q.All(&testFindAllArray)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("FoundAll: ", false)
		return
	}
	if testFindAllArray[0] != nil {
		if _, ok2 := testFindAllArray[0]["_id"]; ok2 && testFindAllArray[0]["_id"].(bson.ObjectId) == testObj["_id"] {
			fmt.Println("FoundAll: ", true)
		}
	} else {
		fmt.Println("FoundAll: ", false)
	}

	err = db.C("testCol").Update(bson.M{"_id": testObj["_id"]}, bson.M{"$set": map[string]interface{}{"updated": true}})
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Updated: ", false)
		return
	}
	db.C("testCol").Find(bson.M{"_id": testObj["_id"]}).One(&testUpdateObj)
	if _, ok := testUpdateObj["updated"]; ok && testUpdateObj["updated"].(bool) == true {
		fmt.Println("Updated: ", true)
	}

	err = db.C("testCol").Remove(bson.M{"_id": testObj["_id"]})
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Removed: ", false)
		return
	}

	q = db.C("testCol").Find(bson.M{"_id": testObj["_id"]})
	err = q.One(nil)
	if err != nil {
		if err.Error() == "not found" {
			fmt.Println("Removed: ", true)
		} else {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("Removed: ", false)
	}
	fmt.Printf("\n'testFindOneObj' %+v\n", testFindOneObj)
	fmt.Printf("\n'testSelectObj' %+v\n", testSelectObj)
	fmt.Printf("\n'testFindAllArray' %+v\n", testFindAllArray)
	fmt.Printf("\n'testUpdateObj' %+v\n", testUpdateObj)

	return
}
