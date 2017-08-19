package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

//App owns the handals for the http request to be able to link the request with the database
type App struct {
	cardDatabase *sql.DB
	router       *mux.Router
}

//GenerateApp constructor for App Struct it connects to the database and then creates
//a new router object
func GenerateApp() App {
	var newApp App
	// connectToDB
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		panic(err)
	}
	newApp.cardDatabase = db

	newApp.router = mux.NewRouter().StrictSlash(true)
	newApp.InitalizeRoutes()
	return newApp

}

//SearchByIDHandle Http Handle to get card my database ID
func (CurrentApp *App) SearchByIDHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["card_id"]

	//fmt.Fprintln(w, "id:", id)

	idAsInt, err := strconv.Atoi(id)

	if err != nil {
		fmt.Fprintln(w, "invalid Card Id")
	} else {

		var c card
		err := c.getCardFromID(CurrentApp.cardDatabase, idAsInt)
		if err != nil || c.ID == 0 {
			fmt.Fprintln(w, "no card with that found found")
			fmt.Fprintln(w, err.Error())

		}
		json.NewEncoder(w).Encode(c)
	}

}

//SearchByNameHandle Http Handle to get card my database ID
func (CurrentApp *App) SearchByNameHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["Name"]

}

//Listen start listeing on given port number
func (CurrentApp *App) Listen(portNumber string) {
	log.Fatal(http.ListenAndServe(":"+portNumber, CurrentApp.router))
}

//Index Home (really just for testing)
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

//InitalizeRoutes Defines the routes for the router to set which get request trigger what function
func (CurrentApp *App) InitalizeRoutes() {
	CurrentApp.router.HandleFunc("/", Index)
	CurrentApp.router.HandleFunc("/card/id/{card_id:[0-9]+}", CurrentApp.SearchByIDHandle)
	CurrentApp.router.HandleFunc("/card/Name/{language_code}/{card_Name}", CurrentApp.SearchByIDHandle)

}
