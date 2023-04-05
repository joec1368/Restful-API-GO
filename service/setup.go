package service

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"unicode"
)

const (
	DB_USER     = ""
	DB_PASSWORD = ""
	DB_NAME     = ""
)

// DB set up
func SetupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	DB, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	return DB
}

func RemoveSpace(s string) string {
	rr := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			rr = append(rr, r)
		}
	}
	return string(rr)
}
