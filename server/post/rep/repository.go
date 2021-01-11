package rep

import (
	"fmt"
	"subd/dz/models"
)

//func InsertPost(p *models.Post) error {
//	err := models.DB.QueryRow("INSERT INTO posts(author, created, forum, message, parent, thread) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;",
//		p.Author, p.Created, p.Forum, p.Message, p.Parent, p.Thread).
//		Scan(&p.Author, &p.Created, &p.Forum, &p.ID, &p.IsEdited, &p.Message, &p.Parent, &p.Thread, &p.Path)
//
//	return err
//}

func InsertPosts(pos *models.Posts, id int, forum string) (*models.Posts, error) {
	posts := models.Posts{}
	if len(*pos) == 0 {
		return &posts, nil
	}
	query := "INSERT INTO posts(author, forum, message, thread, parent, created) VALUES "
	i := 1
	values := make([]interface{}, 0, 2*len(posts))
	for _, post := range *pos {
		query += fmt.Sprintf("('%s', '%s', '%s', %d, $%d, $%d), ",
			post.Author, forum, post.Message, id, i, i+1)
		i += 2
		values = append(values, post.Parent, post.Created)
	}

	rows, err := models.DB.Query(query[:len(query)-2]+" RETURNING *;", values...)
	if err != nil {
		return nil, err
	}

	p := models.Post{}

	for rows.Next() {
		err = rows.Scan(&p.Author, &p.Created, &p.Forum, &p.ID, &p.IsEdited, &p.Message, &p.Parent, &p.Thread, &p.Path)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	return &posts, nil
}

//func CheckPostByThread(post int) bool {
//	var count int
//	err := models.DB.QueryRow("SELECT thread FROM posts WHERE thread = $1;", post).Scan(&count)
//	return err == nil
//}

func FindPosts(author string, limit, since int, sort string, desc bool) (*models.Posts, error) {
	posts := models.Posts{}

	i := 2
	values := make([]interface{}, 0, 3)
	sqlRec := "SELECT author, created, forum, id, message, parent, thread FROM posts WHERE "

	values = append(values, author)
	if sort == "flat" || sort == "" {
		sqlRec += "author ILIKE $1 "
		if since != 0 {
			if desc {
				sqlRec += "AND id < $2 "
			} else {
				sqlRec += "AND id > $2 "
			}
			i++
			values = append(values, since)
		}
		if desc {
			sqlRec += "ORDER BY created DESC, id DESC LIMIT $" + fmt.Sprint(i) + ";"
		} else {
			sqlRec += "ORDER BY created ASC, id LIMIT $" + fmt.Sprint(i) + ";"
		}
	} else if sort == "tree" {
		sqlRec += "author ILIKE $1 "
		if since != 0 {
			if desc {
				sqlRec += "AND PATH < (SELECT path FROM posts WHERE id = $2) "
			} else {
				sqlRec += "AND PATH > (SELECT path FROM posts WHERE id = $2) "
			}
			i++
			values = append(values, since)
		}
		if desc {
			sqlRec += "ORDER BY path DESC, id  DESC LIMIT $" + fmt.Sprint(i) + ";"
		} else {
			sqlRec += "ORDER BY path, id LIMIT $" + fmt.Sprint(i) + ";"
		}
	} else {
		sqlRec += "path[1] IN (SELECT id FROM posts WHERE author = $1 AND parent IS NULL "
		if since != 0 {
			if desc {
				sqlRec += "AND PATH[1] < (SELECT path[1] FROM posts WHERE id = $2) "
			} else {
				sqlRec += "AND PATH[1] > (SELECT path[1] FROM posts WHERE id = $2) "
			}
			i++
			values = append(values, since)
		}
		if desc {
			sqlRec += "ORDER BY id DESC LIMIT $" + fmt.Sprint(i) + ")ORDER BY path[1] DESC, path, id;"
		} else {
			sqlRec += "ORDER BY id LIMIT $" + fmt.Sprint(i) + ")ORDER BY path, id;"
		}
	}
	values = append(values, limit)

	rows, err := models.DB.Query(sqlRec, values...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	p := models.Post{}

	for rows.Next() {
		err = rows.Scan(&p.Author, &p.Created, &p.Forum, &p.ID, &p.Message, &p.Parent, &p.Thread)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	return &posts, nil
}

func FindByID(id int) (*models.Post, error) {
	post := models.Post{}
	err := models.DB.QueryRow("SELECT author, created, forum, id, is_edited, message, parent, thread FROM posts WHERE id = $1;", id).
		Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)

	return &post, err
}

func UpdatePost(post *models.Post, postUpdate models.PostUpdate) error {
	if postUpdate.Message != "" && postUpdate.Message != post.Message {
		err := models.DB.QueryRow("UPDATE posts SET message=$1, is_edited=true WHERE id=$2 RETURNING *;", postUpdate.Message, post.ID).
			Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &post.Path)
		return err
	}
	return nil
}
