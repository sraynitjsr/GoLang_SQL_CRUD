package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/juju/ratelimit"
	"github.com/sraynitjsr/controller"
	"github.com/sraynitjsr/repository"
	"github.com/sraynitjsr/service"
)

func main() {
	db, err := repository.ConnectToDB("username:password@tcp(127.0.0.1:3306)/database_name")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	limiter := ratelimit.NewBucketWithRate(100, 100)

	router := mux.NewRouter()
	router.HandleFunc("/users", limitMiddleware(userController.GetUsers, limiter)).Methods("GET")
	router.HandleFunc("/users/{id}", limitMiddleware(userController.GetUser, limiter)).Methods("GET")
	router.HandleFunc("/users", userController.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", userController.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", userController.DeleteUser).Methods("DELETE")

	http.ListenAndServe(":8000", router)
}

func limitMiddleware(next http.HandlerFunc, limiter *ratelimit.Bucket) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if limiter.TakeAvailable(1) == 0 {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}
