package rep

import (
	"fmt"
	"subd/dz/models"
)

func InsertThread(thread models.Thread) (models.Thread, error) {
	th := models.Thread{}
	err := models.DB.QueryRow("INSERT INTO threads(author, created, forum, message, slug, title) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;",
		thread.Author, thread.Created, thread.Forum, thread.Message, thread.Slug, thread.Title).
		Scan(&th.Author, &th.Created, &th.Forum, &th.ID, &th.Message, &th.Slug, &th.Title, &th.Votes)

	return th, err
}

func CheckThreadByID(id int) bool {
	var count int
	models.DB.QueryRow("SELECT COUNT(id) FROM threads WHERE id = $1;", id).Scan(&count)
	return count > 0
}

func FindThread(slug string) (models.Thread, error) {
	th := models.Thread{}
	err := models.DB.QueryRow("SELECT author, created, forum, id, message, slug, title, votes FROM threads WHERE slug ILIKE $1;", slug).
		Scan(&th.Author, &th.Created, &th.Forum, &th.ID, &th.Message, &th.Slug, &th.Title, &th.Votes)

	return th, err
}

func FindThreadByID(id int) (models.Thread, error) {
	th := models.Thread{}
	err := models.DB.QueryRow("SELECT author, created, forum, id, message, slug, title, votes FROM threads WHERE id = $1;", id).
		Scan(&th.Author, &th.Created, &th.Forum, &th.ID, &th.Message, &th.Slug, &th.Title, &th.Votes)

	return th, err
}

func FindThreads(forum, since string, limit int, desc bool) (models.Threads, error) {
	threads := models.Threads{}
	values := make([]interface{}, 0, 3)

	sqlRow := "SELECT author, created, forum, id, message, slug, title, votes FROM threads WHERE forum ILIKE $1 "
	values = append(values, forum)
	i := 2
	if since != "" {
		if desc {
			sqlRow += "AND created <= $2 "
		} else {
			sqlRow += "AND created >= $2 "
		}
		i++
		values = append(values, since)
	}

	if desc {
		sqlRow += "ORDER BY created DESC LIMIT $" + fmt.Sprint(i) + ";"
	} else {
		sqlRow += "ORDER BY created LIMIT $" + fmt.Sprint(i) + ";"
	}
	values = append(values, limit)
	rows, err := models.DB.Query(sqlRow, values...)

	if err != nil {
		return threads, err
	}
	defer rows.Close()

	th := models.Thread{}

	for rows.Next() {
		err = rows.Scan(&th.Author, &th.Created, &th.Forum, &th.ID, &th.Message, &th.Slug, &th.Title, &th.Votes)
		if err != nil {
			return threads, err
		}
		threads = append(threads, th)
	}
	return threads, nil
}

func FindByPost(id int) (models.Thread, error) {
	thread := models.Thread{}
	err := models.DB.QueryRow("SELECT thread FROM posts WHERE id = $1;", id).
		Scan(&thread.ID)
	if err != nil {
		return thread, err
	}

	thread, err = FindThreadByID(thread.ID)
	if err != nil {
		return thread, err
	}

	return thread, nil
}

func UpdateThread(thread models.Thread, threadUpdate models.ThreadUpdate) (models.Thread, error) {
	values := make([]interface{}, 0, 3)
	i := 1
	var s string
	if threadUpdate.Message != "" {
		values = append(values, threadUpdate.Message)
		s += " message=$" + fmt.Sprint(i)
		i++
	}
	if threadUpdate.Title != "" {
		values = append(values, threadUpdate.Title)

		if i > 1 {
			s += ","
		}
		s += " title=$" + fmt.Sprint(i)
		i++
	}
	threadUp := models.Thread{}

	if i > 1 {
		sqlRow := "UPDATE threads SET" + s + " WHERE slug=$" + fmt.Sprint(i) + " RETURNING *;"
		values = append(values, thread.Slug)
		err := models.DB.QueryRow(sqlRow, values...).
			Scan(&threadUp.Author, &threadUp.Created, &threadUp.Forum, &threadUp.ID, &threadUp.Message, &threadUp.Slug, &threadUp.Title, &threadUp.Votes)

		return threadUp, err
	} else {
		err := models.DB.QueryRow("SELECT author, created, forum, id, message, slug, title, votes FROM threads WHERE slug ILIKE $1;", thread.Slug).
			Scan(&threadUp.Author, &threadUp.Created, &threadUp.Forum, &threadUp.ID, &threadUp.Message, &threadUp.Slug, &threadUp.Title, &threadUp.Votes)

		return threadUp, err
	}
}

func InsertVote(vote models.Vote) error {
	_, err := models.DB.Exec("INSERT INTO votes(nickname, voice, thread) VALUES ($1, $2, NULLIF($3, 0));", vote.Nickname, vote.Voice, vote.Thread)

	return err
}

func UpdateVote(vote models.Vote) error {
	_, err := models.DB.Exec("UPDATE votes SET voice=$1 WHERE nickname=$2 AND thread=$3;", vote.Voice, vote.Nickname, vote.Thread)

	return err
}
