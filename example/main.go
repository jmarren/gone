package main

import (
	"fmt"
	"net/http"

	"github.com/jmarren/gone/example/middleware"
	"github.com/jmarren/gone/gone"
)

type appData struct {
	username string
}

func main() {

	// Creates a new gone application
	gone := gone.New[*appData]()

	data := new(appData)

	data.username = "john"

	gone.SetData(data)

	// the returned Route handles '/' by default

	// you can register a get handler with route.Get
	gone.Get(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("'/' route")
		w.Write([]byte("/"))
	}))

	// Append a route onto another with gone.Then(<pattern>)
	// this returns a new route that you can register additional routes with
	songs := gone.Then("songs/")

	songs.Get(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("songs route")
		fmt.Printf("data: %s\n", songs.GetData().username)
		w.Write([]byte("songs"))
	}))

	// register/append another subroute to the '/' route (named 'gone' here)
	users := gone.Then("users/")

	users.Get(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("users route")
		w.Write([]byte("/users"))
	}))

	// append a {id} route to the users/ route
	id := users.Then("{id}")

	// use some middlewares for the /users/{id} route
	id.Use(middleware.LogHiMiddleware, middleware.LogBaconMiddleware)

	id.Get(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("id route")
		idVal := r.PathValue("id")
		fmt.Printf("id: %s\n", idVal)
		w.Write([]byte(idVal))
	}))

	// serve the application on the provided port
	gone.Serve(":8080")
}
