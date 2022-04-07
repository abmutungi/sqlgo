package main

import (
	"database/sql"
	//"fmt"
	"net/http"
	//"html/template"
	"sqlgo/crud"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	db, _ = sql.Open("sqlite3", "database/godb.db")
	db.Exec("create table if not exists testTable (id integer primary key,username text, surname text,age Integer,university text)")

	// crud.AddUser(db, "Natasha", "Mutungi", 29, "University of Sheffield") // added data to database
	// crud.AddUser(db, "Alexandra", "Mutungi", 22, "")
	// crud.UpdateUser(db, 27, "Arnold", "Mutungi", 29, "01 Founders") // update data to database

	//crud.DeleteUser(db, 1) // delete data to database

	// fmt.Println(crud.GetUsers(db, 2)) // printing the user

	// stmnt, _ := db.Prepare("INSERT INTO testTable (id, username, surname, age, university) VALUES (?, ?, ?, ?, ?)")
	// stmnt.Exec(25, "Arnold", "Mutungi", 29, "University of Nottingham")

	crud.UpdateUser(db, 1, "test", "test", 20, "university")


	http.HandleFunc("/", crud.Handler)
	http.HandleFunc("/insert", crud.InsertHandler)
	http.HandleFunc("/browse", crud.BrowseHandler)
	http.HandleFunc("/delete/", crud.DeleteHandler)
	http.HandleFunc("/update/", crud.UpdateHandler)
	http.HandleFunc("/updateresult/", crud.UpdateResultHandler)

	http.ListenAndServe(":8080", nil)
}
