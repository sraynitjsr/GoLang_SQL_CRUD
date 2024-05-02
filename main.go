package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/database_name")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var users []User

	result, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var user User
		err := result.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	result, err := db.Query("SELECT id, name, age FROM users WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	var user User
	for result.Next() {
		err := result.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	result, err := db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", user.Name, user.Age)
	if err != nil {
		panic(err.Error())
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	user.ID = int(lastInsertID)

	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.ID, _ = strconv.Atoi(params["id"])

	_, err := db.Exec("UPDATE users SET name = ?, age = ? WHERE id = ?", user.Name, user.Age, user.ID)
	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "User with ID = %s was deleted", params["id"])
}
