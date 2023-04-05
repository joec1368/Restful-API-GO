package main

import (
	"Dcard/controller"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/header", controller.GetAllHeader).Methods("GET")
	router.HandleFunc("/header", controller.AddHeader).Methods("POST")
	router.HandleFunc("/header/GetHead/{header}", controller.GetCertainHeader).Methods("POST")
	router.HandleFunc("/header/Clear", controller.ClearPage).Methods("DELETE")
	router.HandleFunc("/header/ClearAll", controller.ClearAll).Methods("DELETE")
	router.HandleFunc("/header/{header}", controller.DeleteHeader).Methods("DELETE")
	router.HandleFunc("/header/{header}", controller.UpdateHeader).Methods("PUT")
	router.HandleFunc("/page/setPage/{header}", controller.AddPage).Methods("POST")
	router.HandleFunc("/page/getPage", controller.GetPage).Methods("POST")
	router.HandleFunc("/page/{header}", controller.UpdatePage).Methods("PUT")
	router.HandleFunc("/page/{header}", controller.DeletePage).Methods("DELETE")

	fmt.Println("Server at 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
