package main

import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	auth "github.com/ondro2208/dokkuapi/authentication"
	"github.com/ondro2208/dokkuapi/controller"
	log "github.com/ondro2208/dokkuapi/logger"
	"net/http"
	"os"
)

func main() {
	file, err := os.OpenFile("dokkuapi_webserver.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.ErrorLogger.Fatal(err)
	}
	defer file.Close()

	router := mux.NewRouter()
	router.Handle("/info", auth.IsAuthenticated(getInfo)).Methods("GET")
	router.HandleFunc("/register", controller.RegisterUserEndpoint).Methods("POST")
	loggedRouter := handlers.LoggingHandler(file, router)
	http.ListenAndServe(":3000", loggedRouter)
}

func getInfo(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	jsonData, err := json.Marshal(map[string]string{"message": "temporary info endpoint"})
	if err != nil {
		log.ErrorLogger.Fatal(err)
	}
	response.Write(jsonData)
}
