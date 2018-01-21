package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

// ValRow represents a row from our 'vals' table
type ValRow struct {
	Id  int
	Val string
}

// MySQLConnectionString is the connection string for your
// database. It should be set as an environment variable.
var MySQLConnectionString = os.Getenv("MYSQL_CONNECTION_STRING")

func getVals(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	db, err := sql.Open("mysql", MySQLConnectionString)

	if err != nil {
		t, _ := template.ParseFiles("views/error.html")
		t.Execute(res, err.Error())
	}

	rows, err := db.Query("SELECT * FROM vals")

	if err != nil {
		t, _ := template.ParseFiles("views/error.html")
		t.Execute(res, err.Error())
	} else {
		var vals []ValRow
		for rows.Next() {
			var id int
			var val string
			_ = rows.Scan(&id, &val)
			row := ValRow{id, val}
			vals = append(vals, row)
		}

		t, _ := template.ParseFiles("views/index.html")
		t.Execute(res, vals)
	}
}

func createVal(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	req.ParseForm()
	newvals := req.Form["val"]

	if newvals == nil || len(newvals) < 1 {
		t, _ := template.ParseFiles("views/error.html")
		t.Execute(res, "You didn't submit anything!")
		return
	}

	newval := newvals[0]

	if len(newval) < 5 || len(newval) > 100 {
		t, _ := template.ParseFiles("views/error.html")
		t.Execute(res, "New values must have 5 < length > 100.")
		return
	}

	db, err := sql.Open("mysql", MySQLConnectionString)
	if err != nil {
		t, _ := template.ParseFiles("views/error.html")
		t.Execute(res, err.Error())
		return
	}

	createStatement, err := db.Prepare("INSERT INTO vals (val) VALUES (?)")
	if err != nil {
		t, _ := template.ParseFiles("views/error.html")
		t.Execute(res, err.Error())
		return
	}

	_, err = createStatement.Exec(newval)
	if err != nil {
		t, _ := template.ParseFiles("views/error.html")
		t.Execute(res, err.Error())
		return
	}

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func main() {
	router := httprouter.New()

	router.GET("/", getVals)
	router.POST("/vals", createVal)

	log.Println("Server is running at http://localhost:8080")

	http.ListenAndServe(":8080", router)
}
