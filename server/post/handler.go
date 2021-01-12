package post

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"subd/dz/models"
	forumRep "subd/dz/server/forum/rep"
	postRep "subd/dz/server/post/rep"
	threadRep "subd/dz/server/thread/rep"
	userRep "subd/dz/server/user/rep"
)

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slugOrID := vars["slug_or_id"]

	var thread *models.Thread
	var err error
	id, errInt := strconv.Atoi(slugOrID)

	if errInt != nil {
		thread, err = threadRep.FindThread(slugOrID)
	} else {
		thread, err = threadRep.FindThreadByID(id)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write(models.MarshalErrorSt("Can't find post thread"))
			return
		}
		// log.Println(err)
		return
	}

	posts := models.Posts{}
	err = json.NewDecoder(r.Body).Decode(&posts)
	if err != nil {
		// log.Println(err)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("[]"))
		return
	}

	err = postRep.InsertPosts(&posts, thread.ID, thread.Forum)

	if err != nil {
		if strings.Contains(err.Error(), "bad parent thread") {
			w.WriteHeader(http.StatusConflict)
			w.Write(models.MarshalErrorSt("Parent post was created in another thread"))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.MarshalErrorSt("Can't find post author by nickname"))
		return
	}

	body, err := json.Marshal(posts)
	if err != nil {
		// log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func ThreadPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}

	since, err := strconv.Atoi(r.URL.Query().Get("since"))
	if err != nil {
		since = 0
	}
	sort := r.URL.Query().Get("sort")

	desc, err := strconv.ParseBool(r.URL.Query().Get("desc"))
	if err != nil {
		desc = false
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

	posts, err := postRep.FindPosts(thread, limit, since, sort, desc)
	if err != nil {
		// log.Println(err)
		return
	}

	body, err := json.Marshal(posts)
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
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// log.Println(err)
		return
	}

	related := r.URL.Query().Get("related")

	post, err := postRep.FindByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write(models.MarshalErrorSt("Can't find post by id"))
			return
		}
		// log.Println(err)
		return
	}

	if r.Method == "GET" {
		postFull := models.PostFull{}
		if strings.Contains(related, "user") {
			user, err := userRep.FindByNickname(post.Author)
			if err != nil {
				// log.Println(err)
				return
			}
			postFull.Author = user
		}

		if strings.Contains(related, "forum") {
			forum, err := forumRep.FindForum(post.Forum)
			if err != nil {
				// log.Println(err)
				return
			}
			postFull.Forum = forum
		}

		if strings.Contains(related, "thread") {
			thread, err := threadRep.FindThreadByID(post.Thread)
			if err != nil {
				// log.Println(err)
				return
			}
			postFull.Thread = thread
		}

		postFull.Post = post

		body, err := json.Marshal(postFull)
		if err != nil {
			// log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
		return
	}

	postUpdate := models.PostUpdate{}
	err = json.NewDecoder(r.Body).Decode(&postUpdate)
	if err != nil {
		// log.Println(err)
		return
	}

	err = postRep.UpdatePost(post, postUpdate)
	if err != nil {
		// log.Println(err)
		return
	}

	body, err := json.Marshal(post)
	if err != nil {
		// log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
