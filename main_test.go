package main

import (
	"errors"
	"github.com/djak250/mgo-wrapper-interface/models"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestRunDbSuite(t *testing.T) {
	tmdb := models.TMongoDatabase{
		[]*models.TMongoCollection{},
		t,
	}

	testObj := map[string]interface{}{
		"_id":   bson.NewObjectId(),
		"test1": "1",
		"test2": "2",
		"test3": "3",
	}

	tmdb.ExpectInsert("testCol", nil)
	testObjSlice := []map[string]interface{}{testObj}
	tmdb.ExpectFind("testCol", testObj, testObjSlice)
	tmdb.ExpectUpdate("testCol", nil)

	testObj2 := map[string]interface{}{
		"_id":     bson.NewObjectId(),
		"test1":   "1",
		"test2":   "2",
		"test3":   "3",
		"updated": true,
	}
	testObjSlice2 := []map[string]interface{}{testObj}

	tmdb.ExpectFind("testCol", testObj2, testObjSlice2)
	tmdb.ExpectRemove("testCol", nil)
	tmdb.ExpectFind("testCol", errors.New("not found"), errors.New("not found"))

	runDbSuite(tmdb, testObj)
}
