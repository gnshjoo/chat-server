package main

import (
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

var (
	renderer *render.Render
	mongoSession *mgo.Session
	upgrader = &websocket.Upgrader{
		ReadBufferSize: socketBufferSize,
		WriteBufferSize: socketBufferSize,
	}
)

const (
	sessionKey = "simple_chat_session"
	sessionSecret = "simple_chat_session_secret"
	socketBufferSize = 1024
)

func init() {
	renderer = render.New()

	s, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	mongoSession = s
}

func main() {
	// create router
	router := httprouter.New()

	// defined handler
	router.GET("/", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		// template redner
		renderer.HTML(writer, http.StatusOK, "index", map[string]string{"title":"Simple Chat!"})
	})

	router.GET("/login", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		renderer.HTML(writer, http.StatusOK, "login", nil)
	})

	router.GET("/logout", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		sessions.GetSession(request).Delete(currentUserKey)
		http.Redirect(writer, request, "/login", http.StatusFound)
	})

	router.GET("/auth/:action/:provider", loginHandler)

	router.POST("/rooms", createRoom)
	router.GET("/rooms", retrieveRooms)

	router.GET("/rooms/:id/messages", retrieveMessages)

	router.GET("/ws/:room_id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
			socket , err := upgrader.Upgrade(writer, request, nil)
			if err != nil {
				log.Fatalln("ServerHTTP : ", err)
				return
			}
			newClient(socket, params.ByName("room_id"), GetCurrentUser(request))
	})

	// negroni middleware create
	n := negroni.Classic()
	store := cookiestore.New([]byte(sessionSecret))
	n.Use(sessions.Sessions(sessionKey, store))


	// negroni add router fo handler
	n.UseHandler(router)

	// start web server
	n.Run(":8080")
}
