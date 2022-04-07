package crud

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id         int
	Username   string
	Surname    string
	Age        int
	University string

	// I created a struct with a struct to select the rows in the table and add data.
}

var (
	tpl, _ = template.ParseGlob("templates/*.htm")
	db     *sql.DB
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}

	// catch to error.
}

func AddUser(db *sql.DB, username string, surname string, age int, university string) {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into testTable (username, surname, age, university) values (?,?,?,?)")

	// check rows affected
	res, err := stmt.Exec(username, surname, age, university)

	rowsAffec, _ := res.RowsAffected()
	if err != nil || rowsAffec != 1 {
		fmt.Println("Error inserting row:", err)
		// tpl.ExecuteTemplate(w, "insert.htm", "Error inserting data, please check all fields.")
		// return
	}
	lastInserted, _ := res.LastInsertId()
	rowsAffected, _ := res.RowsAffected()
	fmt.Println("ID of last row inserted:", lastInserted)
	fmt.Println("number of rows affected:", rowsAffected)

	tx.Commit()
}

func GetUsers(db *sql.DB, id2 int) User {
	rows, err := db.Query("select * from testTable")
	checkError(err)
	for rows.Next() {
		var tempUser User
		err = rows.Scan(&tempUser.Id, &tempUser.Username, &tempUser.Surname, &tempUser.Age, &tempUser.University)
		checkError(err)
		if tempUser.Id == id2 {
			return tempUser
		}

	}
	return User{}
}

func UpdateUser(db *sql.DB, id int, username string, surname string, age int, university string) {
	//age2 := strconv.Itoa(age)
	tx, _ := db.Begin()

	
	stmt, _ := tx.Prepare("update testTable set username=?, surname=?, age=?, university=? where id=?")
	res, err := stmt.Exec(username, surname, age, university, id)
	checkError(err)
	tx.Commit()

	rowsAff, _ := res.RowsAffected()
	if err != nil || rowsAff != 1 {
		fmt.Println(err)
		// tpl.ExecuteTemplate(w, "result.htm", "Error updating data, please check all fields.")
		// return
	}

}

func DeleteUser(db *sql.DB, id2 int) {
	sid := strconv.Itoa(id2) // int to string
	tx, _ := db.Begin()

	stmt, _ := tx.Prepare("delete from testTable where id=?")

	// check rows affected
	res, err := stmt.Exec(sid)

	rowsAffec, _ := res.RowsAffected()
	if err != nil || rowsAffec != 1 {
		fmt.Println("Error deleting user:", err)
		// tpl.ExecuteTemplate(w, "insert.htm", "Error deleting data, please check all fields.")
		// return
	}
	lastInserted, _ := res.LastInsertId()
	rowsAffected, _ := res.RowsAffected()
	fmt.Println("ID of last row inserted:", lastInserted)
	fmt.Println("number of rows affected:", rowsAffected)

	tx.Commit()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "forum testing")
}

func InsertHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****insertHandler running*****")
	db, _ = sql.Open("sqlite3", "database/godb.db")

	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "insert.htm", nil)
		return
	}
	r.ParseForm()

	username := r.FormValue("firstName")
	surname := r.FormValue("lastName")
	university := r.FormValue("uniName")
	age := r.FormValue("age")

	var err error
	if username == "" || surname == "" || university == "" || age == "" {
		fmt.Println("Error inserting row:", err)
		tpl.ExecuteTemplate(w, "insert.htm", "Error inserting data, please check all fields.")
		return
	}

	ageint, _ := strconv.Atoi(age)

	AddUser(db, username, surname, ageint, university)

	if err != nil {
		fmt.Println("Error inserting row:", err)
		tpl.ExecuteTemplate(w, "insert.htm", "Error inserting data, please check all fields.")
		return
	}

	tpl.ExecuteTemplate(w, "insert.htm", "User Successfully Added")
}

func BrowseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("************browserHandler running********************")

	stmt := "select * from testTable"

	rows, err := db.Query(stmt)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User

		err = rows.Scan(&u.Id, &u.Username, &u.Surname, &u.Age, &u.University)
		if err != nil {
			panic(err)
		}
		users = append(users, u)
	}
	fmt.Println(users)
	tpl.ExecuteTemplate(w, "select.htm", users)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****deleteHandler running*****")
	db, _ = sql.Open("sqlite3", "database/godb.db")

	r.ParseForm()
	id := r.FormValue("idusers")
	id2, _ := strconv.Atoi(id)

	DeleteUser(db, id2)

	tpl.ExecuteTemplate(w, "result.htm", "User Successfully Deleted")
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****updateHandler running*****")
	r.ParseForm()

	id := r.FormValue("idusers")
	id2, _ := strconv.Atoi(id)

	row := db.QueryRow("select * from testTable where id = ?;", id2)

	var u User

	err := row.Scan(&u.Id, &u.Username, &u.Surname, &u.Age, &u.University)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/browse", 307)
		return
	}
	tpl.ExecuteTemplate(w, "update.htm", u)
}

func UpdateResultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****updateResultHandler running*****")

	r.ParseForm()

	id := r.FormValue("idusers")
	username := r.FormValue("firstName")
	surname := r.FormValue("lastName")
	age := r.FormValue("age")
	university := r.FormValue("uniName")

	ageint, _ := strconv.Atoi(age)	
	id2, _ := strconv.Atoi(id)


	UpdateUser(db, id2, username, surname, ageint, university)

	tpl.ExecuteTemplate(w, "result.htm", "User was successfully updated")
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/browse", 307)
}
