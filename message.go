package main

import (
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"time"
)

const messageFetchSize =10

type Message struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	RoomId    bson.ObjectId `bson:"room_id" json:"room_id"`
	Content   string        `bson:"content" json:"content"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	User      *User         `bson:"user" json:"user"`
}

func (m *Message) create() error {
	session := mongoSession.Copy()
	defer session.Close()

	m.ID = bson.NewObjectId()
	m.CreatedAt = time.Now()
	c := session.DB("test").C("message")

	if err := c.Insert(m); err != nil {
		return err
	}
	return nil
}

func retrieveMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mongoSession.Copy()
	defer session.Close()

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = messageFetchSize
	}

	var messages []Message
	err = session.DB("test").C("messages").
		Find(bson.M{"room_id": bson.ObjectIdHex(ps.ByName("id"))}).
		Sort("-_id").Limit(limit).All(&messages)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
	}
	renderer.JSON(w, http.StatusOK, messages)
}