package controller

import (
	"Dcard/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type Header struct {
	ID       string `json:"id"`
	Page_key string `json:"page_Key"`
	Final    string `json:"final"`
}

func GetAllHeader(w http.ResponseWriter, r *http.Request) {
	db := service.SetupDB()
	rows, err := db.Query("SELECT * FROM header")
	if err != nil {
		panic(err)
	}

	var header []Header
	for rows.Next() {
		var id string
		var page_Key string
		var final string
		err = rows.Scan(&id, &page_Key, &final)
		header = append(header, Header{ID: service.RemoveSpace(id), Page_key: service.RemoveSpace(page_Key), Final: service.RemoveSpace(final)})
	}

	json.NewEncoder(w).Encode(header)
}

func AddHeader(w http.ResponseWriter, r *http.Request) {
	db := service.SetupDB()
	respBody, _ := ioutil.ReadAll(r.Body)
	var body []map[string]string
	json.Unmarshal(respBody, &body)
	for _, value := range body {
		sqlStatement := `INSERT INTO "header" ("header") VALUES ($1)`
		_, err := db.Exec(sqlStatement, value["name"])
		if err != nil {
			var response string = "insert " + value["name"] + " fail"
			json.NewEncoder(w).Encode(response)
			panic(err)
		} else {
			var response string = "insert " + value["name"] + " succeed"
			json.NewEncoder(w).Encode(response)
		}

	}
}

func GetCertainHeader(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	header := vars["header"]
	if header == "" {
		header = r.URL.Query().Get("header")
	}
	json.NewEncoder(w).Encode(service.GetCertainHeader(header))
}

func DeleteHeader(w http.ResponseWriter, r *http.Request) {
	db := service.SetupDB()
	vars := mux.Vars(r)
	header := vars["header"]
	sqlStatement := `DELETE FROM "header" WHERE "header" = $1`
	_, err := db.Exec(sqlStatement, header)
	if err != nil {
		json.NewEncoder(w).Encode("error in DeleteHeader")
	} else {
		//response := "delete " + header + " succeed"
		//json.NewEncoder(w).Encode(response)
	}
}

func ClearPage(w http.ResponseWriter, r *http.Request) {
	db := service.SetupDB()
	sqlStatement := `TRUNCATE TABLE "page"`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		json.NewEncoder(w).Encode("error in ClearPage")
	} else {
		response := "delete all data succeed"
		json.NewEncoder(w).Encode(response)
	}
}

func ClearAll(w http.ResponseWriter, r *http.Request) {
	db := service.SetupDB()
	sqlStatement := `TRUNCATE TABLE "page"`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		json.NewEncoder(w).Encode("error in ClearAll")
	} else {
		response := "delete all data succeed"
		json.NewEncoder(w).Encode(response)
		sqlStatement := `TRUNCATE TABLE "header"`
		_, err := db.Exec(sqlStatement)
		if err != nil {
			json.NewEncoder(w).Encode("error in ClearAll")
		} else {
			//response = "delete all header succeed"
			//json.NewEncoder(w).Encode(response)
		}
	}
}

func UpdateHeader(w http.ResponseWriter, r *http.Request) {
	db := service.SetupDB()
	vars := mux.Vars(r)
	header := vars["header"]
	respBody, _ := ioutil.ReadAll(r.Body)
	var body []map[string]string
	json.Unmarshal(respBody, &body)
	sqlStatement := `UPDATE "header" SET "header" = $2 WHERE "header" = $1`
	_, err := db.Exec(sqlStatement, header, body[0]["newHeader"])
	if err != nil {
		json.NewEncoder(w).Encode("error in UpdateHeader")
	} else {
		//json.NewEncoder(w).Encode("error in UpdateHeader")
	}
}
