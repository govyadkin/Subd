package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"subd/dz/models"
	"subd/dz/server/forum"
	"subd/dz/server/post"
	"subd/dz/server/thread"
	"subd/dz/server/user"
)

func DBConnection() *sql.DB {
	connString := "host=localhost user=misha password=password dbname=subdproject sslmode=disable"
	//connString := "host=localhost user=password password=password dbname=password sslmode=disable"

	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	models.DB = DBConnection()

	router := mux.NewRouter()

	router.HandleFunc("/api/forum/create", forum.Create).Methods(http.MethodPost)
	router.HandleFunc("/api/forum/{slug}/details", forum.Details).Methods(http.MethodGet)
	router.HandleFunc("/api/forum/{slug}/create", thread.Create).Methods(http.MethodPost)
	router.HandleFunc("/api/forum/{slug}/users", forum.Users).Methods(http.MethodGet)
	router.HandleFunc("/api/forum/{slug}/threads", forum.Threads).Methods(http.MethodGet)

	router.HandleFunc("/api/post/{id:[0-9]+}/details", post.Details).Methods(http.MethodGet, http.MethodPost)

	router.HandleFunc("/api/service/clear", forum.ClearHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/service/status", forum.StatusHandler).Methods(http.MethodGet)

	router.HandleFunc("/api/thread/{slug_or_id}/create", post.Create).Methods(http.MethodPost)
	router.HandleFunc("/api/thread/{slug_or_id}/details", thread.Details).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/api/thread/{slug_or_id}/posts", post.ThreadPosts).Methods(http.MethodGet)
	router.HandleFunc("/api/thread/{slug_or_id}/vote", thread.Vote).Methods(http.MethodPost)

	router.HandleFunc("/api/user/{nickname}/create", user.Create).Methods(http.MethodPost)
	router.HandleFunc("/api/user/{nickname}/profile", user.Profile).Methods(http.MethodGet, http.MethodPost)

	http.Handle("/", router)

	fmt.Println("Starting server at: 5000")
	http.ListenAndServe(":5000", nil)
}
