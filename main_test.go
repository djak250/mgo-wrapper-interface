package main

import (
	"errors"
	"github.com/djak250/mgo-wrapper-interface/models"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

// Run with test -v to show the output of right values

func TestRunDbSuite(t *testing.T) {
	tmdb := models.TMongoDatabase{
		[]*models.TMongoCollection{},
		t,
	}
	objId := bson.NewObjectId()
	testObj := map[string]interface{}{
		"_id":   objId,
		"test1": "1",
		"test2": "2",
		"test3": "3",
	}

	tmdb.ExpectInsert("testCol", nil)
	testObjBson := bson.M{
		"_id":   objId,
		"test1": "1",
		"test2": "2",
		"test3": "3",
	}
	testObjBsonSlice := []bson.M{testObjBson}

	tmdb.ExpectFind("testCol", testObjBson, testObjBsonSlice)
	tmdb.ExpectUpdate("testCol", nil)

	testObj2Bson := bson.M{
		"_id":     objId,
		"test1":   "1",
		"test2":   "2",
		"test3":   "3",
		"updated": true,
	}
	testObjBsonSlice2 := []bson.M{testObj2Bson}

	tmdb.ExpectFind("testCol", testObj2Bson, testObjBsonSlice2)
	tmdb.ExpectRemove("testCol", nil)
	tmdb.ExpectFind("testCol", errors.New("not found"), errors.New("not found"))

	runDbSuite(tmdb, testObj)
}
