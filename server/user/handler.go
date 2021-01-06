package user

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"subd/dz/models"
	"subd/dz/server/user/rep"
)

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	user.Nickname = vars["nickname"]

	users, err := rep.ConflictUsers(user.Email, user.Nickname)
	if err != nil {
		log.Println(err)
		return
	}

	if len(users) > 0 {
		body, err := json.Marshal(users)
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusConflict)
		w.Write(body)
		return
	}

	err = rep.Create(user)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func Profile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	nickname := vars["nickname"]

	if r.Method == "GET" {
		user, err := rep.FindByNickname(nickname)

		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write(models.MarshalErrorSt("Not fund user"))
			return
		}

		if err != nil {
			log.Println(err)
			return
		}

		body, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
		return
	}

	if !rep.CheckByNickname(nickname) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.MarshalErrorSt("Can't find user"))
		return
	}

	userUpdate := models.UserUpdate{}
	err := json.NewDecoder(r.Body).Decode(&userUpdate)
	if err != nil {
		log.Println(err)
		return
	}

	if rep.CheckByEmail(userUpdate.Email) {
		w.WriteHeader(http.StatusConflict)
		w.Write(models.MarshalErrorSt("This email is already registered"))
		return
	}

	res, err := rep.UpdateUser(nickname, userUpdate)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
	return
}
