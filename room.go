package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type Room struct {
	ID bson.ObjectId `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
}

func (r *Room) FieldMap(req *http.Request) binding.FieldMap{
	return binding.FieldMap{&r.Name: "name"}
}

func createRoom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	r := new(Room)
	errs := binding.Bind(req, r)
	if errs != nil {
		return
	}

	session := mongoSession.Copy()
	defer session.Close()
	r.ID = bson.NewObjectId()
	c := session.DB("test").C("rooms")

	if err := c.Insert(r); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	renderer.JSON(w, http.StatusCreated, r)
}

func retrieveRooms(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	session := mongoSession.Copy()
	defer session.Close()

	var rooms []Room
	err := session.DB("test").C("rooms").Find(nil).All(&rooms)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err )
		return
	}
	renderer.JSON(w, http.StatusOK, rooms)
}