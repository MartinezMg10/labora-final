package main

import (
	"log"
	"net/http"
	"time"

	"github.com/MartinezMg10/rest_golang/internal/course"
	"github.com/MartinezMg10/rest_golang/internal/enrollment"
	"github.com/MartinezMg10/rest_golang/internal/user"
	"github.com/MartinezMg10/rest_golang/pkg/boostrap"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	router := mux.NewRouter()
	_ = godotenv.Load()
	l := boostrap.InitLogger()
	db, err := boostrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	userRepo := user.NewRepo(l, db)
	userServ := user.NewService(l, userRepo)
	userEnd := user.MakeEndpoints(userServ)

	courseRepo := course.NewRepo(db, l)
	courseServ := course.NewService(l, courseRepo)
	courseEnd := course.MakeEndpoints(courseServ)

	enrollRepo := enrollment.NewRepo(db, l)
	enrollServ := enrollment.NewService(l, enrollRepo, courseServ, userServ)
	enrollEnd := enrollment.MakeEndpoints(enrollServ)

	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEnd.Delete).Methods("DELETE")

	router.HandleFunc("/courses", courseEnd.Create).Methods("POST")
	router.HandleFunc("/courses/{id}", courseEnd.Get).Methods("GET")
	router.HandleFunc("/courses", courseEnd.GetAll).Methods("GET")
	router.HandleFunc("/courses/{id}", courseEnd.Update).Methods("PATCH")
	router.HandleFunc("/courses/{id}", courseEnd.Delete).Methods("DELETE")

	router.HandleFunc("/enrollments", enrollEnd.Create).Methods("POST")

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
