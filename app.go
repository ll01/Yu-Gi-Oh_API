package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	//_ "github.com/mattn/go-sqlite3"
)

//App owns the handals for the http request to be able to link the request with the database
type App struct {
	cardDatabase *sql.DB
	router       *mux.Router
}

var currentBufferWriter http.ResponseWriter

//GenerateApp constructor for App Struct it connects to the database and then creates
//a new router object
func GenerateApp() App {
	var newApp App
	// connectToDB
	db, err := sql.Open("mysql", "x:x@/card_db")
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
	currentBufferWriter = w
	vars := mux.Vars(r)
	id := vars["card_id"]

	idAsInt, err := strconv.Atoi(id)

	if err != nil {
		fmt.Fprintln(w, "invalid Card Id")
	} else {

		var c card
		c.getCardFromID(CurrentApp.cardDatabase, idAsInt)

		json.NewEncoder(w).Encode(c)
	}

}

func (CurrentApp *App) SearchByNameHandle(w http.ResponseWriter, r *http.Request) {
	currentBufferWriter = w
	pathVariables := mux.Vars(r)

	contryCode := pathVariables["contry_code"]
	contryCode = strings.ToUpper(contryCode)
	cardName := pathVariables["card_name"]
	var cardToGet card
	cardToGet.getCardFromName(CurrentApp.cardDatabase, cardName, contryCode)
	json.NewEncoder(w).Encode(cardToGet)
}

func HandleCardSearchError(err error) {
	if err != nil {
		fmt.Fprintln(currentBufferWriter, "no card with that found found")
		fmt.Fprintln(currentBufferWriter, err.Error())
	}
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
	CurrentApp.router.HandleFunc("/card/passcode/{card_id:[0-9]+}", CurrentApp.SearchByIDHandle)
	CurrentApp.router.HandleFunc("/card/name/{contry_code:[A-Za-z][A-Za-z]}/{card_name}", CurrentApp.SearchByNameHandle)
	CurrentApp.router.HandleFunc("/card/archtype/{archtype_name}", CurrentApp.SearchByArchtypeHandle)
	//CurrentApp.router.HandleFunc("/card/effect/{effect_name}", CurrentApp.SearchByEffectHandle)
	//CurrentApp.router.HandleFunc("/card/effect/{attribute_name}", CurrentApp.SearchByAttributeHandle)
	//CurrentApp.router.HandleFunc("/card/effect/{link_arrows}", CurrentApp.SearchByLinkArrowsHandle)
	//CurrentApp.router.HandleFunc("/card/effect/{comparator}/{scale[1]+[0-3]|[0-9]}", CurrentApp.SearchByLinkscalesHandle)
	//CurrentApp.router.HandleFunc("/card/effect/{comparator}/{level:[1]+[0-2]|[1-9]}", CurrentApp.SearchByLevelHandle)

}
