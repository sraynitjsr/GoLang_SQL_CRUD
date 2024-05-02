package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sraynitjsr/controller"
	"github.com/sraynitjsr/repository"
	"github.com/sraynitjsr/service"
)

func main() {
	db, err := repository.ConnectToDB("username:password@tcp(127.0.0.1:3306)/database_name")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	router := mux.NewRouter()
	router.HandleFunc("/users", userController.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", userController.GetUser).Methods("GET")
	router.HandleFunc("/users", userController.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", userController.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", userController.DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
