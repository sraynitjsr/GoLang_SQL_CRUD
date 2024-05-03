package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/juju/ratelimit"
	"github.com/sraynitjsr/controller"
	"github.com/sraynitjsr/repository"
	"github.com/sraynitjsr/service"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
}

func main() {
	config, err := loadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	db, err := repository.ConnectToDB(config.Database.Username + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + config.Database.Port + ")/" + config.Database.Name)
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

func loadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
