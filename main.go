package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func welcomeFunction(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome"))
}

// func simulateIDFunction(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte(`{"name": "Table", "responsible": "Rakhman", "place": "Informatics", "initial_date": "18.05.2020", "initial_cost": "1000000", "current_cost": "1000000"}`))
// }

const (
	Db_USER     = "admin"
	Db_PASSWORD = "rahman03"
	Db_NAME     = "testdb"
)

type Object struct {
	Do_exist     bool      `json:"do_exist"`
	Name         string    `json:"name"`
	Responsible  string    `json:"responsible"`
	Place        string    `json:"place"`
	Initial_date time.Time `json:"initial_date"`
	Initial_cost int64     `json:"initial_cost"`
}

func NewEmptyObject() Object {
	obj := Object{Name: "", Responsible: "", Place: "", Initial_cost: 0}
	return obj
}

func handleRequests() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/welcome", welcomeFunction)
	http.HandleFunc("/getdata/", get1Function)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func get1Function(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[9:]

	getQuery := fmt.Sprintf("SELECT name, responsible, place, initial_date, initial_cost FROM objects1 WHERE id='%s';", url)
	rows, err := Db.Query(getQuery)

	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer rows.Close()

	obj := NewEmptyObject()

	if rows.Next() {

		err := rows.Scan(&obj.Name, &obj.Responsible, &obj.Place, &obj.Initial_date, &obj.Initial_cost)

		if err != nil {
			fmt.Println(err)
			return
		}
		obj.Do_exist = true
		fmt.Println(true, obj.Name, obj.Responsible, obj.Place, obj.Initial_date, obj.Initial_cost)
	} else {
		obj.Do_exist = false
		fmt.Println(false, obj.Name, obj.Responsible, obj.Place, obj.Initial_date, obj.Initial_cost)
	}

	jsonByte, err := json.Marshal(obj)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}

	w.Write(jsonByte)
	fmt.Println(string(jsonByte))
	if rows.Next() {
		http.Error(w, "Duplicate in one code", http.StatusConflict)
	}
}

var Db *sql.DB

func main() {

	fmt.Println("It works.")

	DbInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", "localhost", 5432, Db_USER, Db_PASSWORD, Db_NAME)

	var err error
	Db, err = sql.Open("postgres", DbInfo)

	if err != nil {
		log.Fatal(err)
	}

	err = Db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	defer Db.Close()

	handleRequests()
}

/*
	XXXUUUUUUUUU - code
	XXX - Company hash
	UUUUUUUUU - Thing ID
*/
