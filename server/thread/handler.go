package thread

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"subd/dz/models"
	forumRep "subd/dz/server/forum/rep"
	threadRep "subd/dz/server/thread/rep"
)

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slug := vars["slug"]

	thread := models.Thread{}
	err := json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {
		// log.Println(err)
		return
	}

	//if !userRep.CheckByNickname(thread.Author) {
	//	w.WriteHeader(http.StatusNotFound)
	//	w.Write(models.MarshalErrorSt("Can't find user"))
	//	return
	//}

	forum, err := forumRep.FindForum(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write(models.MarshalErrorSt("Can't find thread forum"))
			return
		}
		// log.Println(err)
		return
	}

	thread.Forum = slug
	if thread.Slug != "" {
		thread2, err := threadRep.FindThread(thread.Slug)

		if err != sql.ErrNoRows {
			if err != nil {
				// log.Println(err)
				return
			}
			body, err := json.Marshal(thread2)
			if err != nil {
				// log.Println(err)
				return
			}

			w.WriteHeader(http.StatusConflict)
			w.Write(body)
			return
		}
	}

	thread.Forum = forum.Slug
	err = threadRep.InsertThread(&thread)
	if err != nil {
		// log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.MarshalErrorSt("Can't find user"))
		return
	}

	body, err := json.Marshal(thread)
	if err != nil {
		// log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func Vote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vote := models.Vote{}
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		// log.Println(err)
		return
	}

	vars := mux.Vars(r)
	slugOrID := vars["slug_or_id"]

	var thread int
	id, errInt := strconv.Atoi(slugOrID)

	if errInt != nil {
		thread, err = threadRep.FindThreadID(slugOrID)
	} else {
		thread, err = threadRep.FindThreadByIDID(id)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write(models.MarshalErrorSt("Can't find thread"))
			return
		}
		// log.Println(err)
		return
	}

	vote.Thread = thread
	err = threadRep.InsertVote(vote)
	if err != nil {
		// log.Println(err)
		if strings.Contains(err.Error(), "duplicate key") {
			err = threadRep.UpdateVote(vote)
		}
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(models.MarshalErrorSt("Can't find post author by nickname"))
			return
		}
	}

	threadUpdate, err := threadRep.FindThreadByID(thread)
	if err != nil {
		// log.Println(err)
		return
	}

	if errInt != nil && threadUpdate.Votes < 0 {
		threadUpdate.Votes *= -1
	}

	body, err := json.Marshal(threadUpdate)
	if err != nil {
		// log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func Details(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slugOrID := vars["slug_or_id"]

	var thread *models.Thread
	id, errInt := strconv.Atoi(slugOrID)
	var err error
	if errInt != nil {
		thread, err = threadRep.FindThread(slugOrID)
	} else {
		thread, err = threadRep.FindThreadByID(id)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write(models.MarshalErrorSt("Can't find thread"))
			return
		}
		// log.Println(err)
		return
	}
	if r.Method == "GET" {
		body, err := json.Marshal(thread)
		if err != nil {
			// log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
		return
	}

	threadUpdate := models.ThreadUpdate{}
	err = json.NewDecoder(r.Body).Decode(&threadUpdate)
	if err != nil {
		// log.Println(err)
		return
	}

	err = threadRep.UpdateThread(thread, threadUpdate)
	if err != nil {
		// log.Println(err)
		return
	}

	body, err := json.Marshal(thread)
	if err != nil {
		// log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
