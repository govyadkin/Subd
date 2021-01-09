package forum

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"subd/dz/models"
	forumRep "subd/dz/server/forum/rep"
	threadRep "subd/dz/server/thread/rep"
	userRep "subd/dz/server/user/rep"
)

func ClearHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := forumRep.ClearForum()
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("null"))
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := forumRep.StatusForum()
	body, err := json.Marshal(status)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	forum := models.Forum{}
	err := json.NewDecoder(r.Body).Decode(&forum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	user, err := userRep.FindByNickname(forum.User)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write(models.MarshalErrorSt("Can't find user"))
			return
		}
		log.Println(err)
		return
	}

	forum.User = user.Nickname

	conflictForum, err := forumRep.FindForum(forum.Slug)
	if err == sql.ErrNoRows {
		err = forumRep.InsertForum(forum)
		if err != nil {
			log.Println(err)
			return
		}

		body, err := json.Marshal(forum)
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(body)

		return
	}
	if err != nil {
		log.Println(err)
		return
	}

	body, err := json.Marshal(conflictForum)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusConflict)
	w.Write(body)
}

func Details(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slug := vars["slug"]

	forum, err := forumRep.FindForum(slug)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write(models.MarshalErrorSt("Can't find forum"))
			return
		}
		log.Println(err)
		return
	}

	body, err := json.Marshal(forum)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func Users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slug := vars["slug"]

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}

	since := r.URL.Query().Get("since")

	desc, err := strconv.ParseBool(r.URL.Query().Get("desc"))
	if err != nil {
		desc = false
	}

	if !forumRep.CheckForum(slug) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.MarshalErrorSt("Can't find forum"))
		return
	}

	users, err := userRep.FindByForum(slug, since, limit, desc)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func Threads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}

	since := r.URL.Query().Get("since")

	desc, err := strconv.ParseBool(r.URL.Query().Get("desc"))
	if err != nil {
		desc = false
	}

	vars := mux.Vars(r)
	slug := vars["slug"]

	if !forumRep.CheckForum(slug) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.MarshalErrorSt("Can't find forum"))
		return
	}

	threads, err := threadRep.FindThreads(slug, since, limit, desc)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := json.Marshal(threads)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
