package rep

import (
	"fmt"
	"github.com/lib/pq"
	"subd/dz/models"
)

func Create(user models.User) error {
	_, err := models.DB.Exec("INSERT INTO users(about, email, fullname, nickname) VALUES ($1, $2, $3, $4);",
		user.About, user.Email, user.Fullname, user.Nickname)

	return err
}

func CheckByEmail(email string) bool {
	var mail string
	err := models.DB.QueryRow("SELECT email FROM users WHERE email ILIKE $1;", email).Scan(&mail)
	return err == nil
}

func CheckByNickname(nickname string) bool {
	var name string
	err := models.DB.QueryRow("SELECT nickname FROM users WHERE nickname ILIKE $1;", nickname).Scan(&name)
	return err == nil
}

func ConflictUsers(email, nickname string) (*models.Users, error) {
	users := models.Users{}
	user := models.User{}

	rows, err := models.DB.Query("SELECT about, email, fullname, nickname FROM users WHERE nickname ILIKE $2 OR email ILIKE $1;", email, nickname)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}

func FindByNickname(nickname string) (*models.User, error) {
	user := models.User{}

	err := models.DB.QueryRow("SELECT about, email, fullname, nickname FROM users WHERE nickname ILIKE $1;", nickname).
		Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

	return &user, err
}

func UpdateUser(nickname string, userUpdate models.UserUpdate) (*models.User, error) {
	values := make([]interface{}, 0, 3)
	i := 1
	var s string
	if userUpdate.About != "" {
		values = append(values, userUpdate.About)
		s += " about=$" + fmt.Sprint(i)
		i++
	}
	if userUpdate.Email != "" {
		values = append(values, userUpdate.Email)

		if i > 1 {
			s += ","
		}
		s += " email=$" + fmt.Sprint(i)
		i++
	}
	if userUpdate.Fullname != "" {
		values = append(values, userUpdate.Fullname)
		if i > 1 {
			s += ","
		}
		s += " fullname=$" + fmt.Sprint(i)
		i++
	}
	user := models.User{}

	user.Nickname = nickname

	if i > 1 {
		sqlRow := "UPDATE users SET" + s + " WHERE nickname ILIKE $" + fmt.Sprint(i) + " RETURNING fullname, about, email;"
		values = append(values, nickname)
		err := models.DB.QueryRow(sqlRow, values...).Scan(&user.Fullname, &user.About, &user.Email)
		return &user, err
	} else {
		err := models.DB.QueryRow("SELECT about, email, fullname FROM users WHERE nickname ILIKE $1;", nickname).
			Scan(&user.About, &user.Email, &user.Fullname)
		return &user, err
	}
}

func FindByForum(slug, since string, limit int, desc bool) (*models.Users, error) {
	//users := models.Users{}
	//sqlRek := "SELECT users.about, users.email, users.fullname, users.nickname FROM forum_users " +
	//	"JOIN users ON users.nickname = forum_users.author WHERE slug = $1 "
	//
	//values := make([]interface{}, 0, 3)
	//
	//values = append(values, slug)
	//i := 2
	//if since != "" {
	//	if desc {
	//		sqlRek += `AND LOWER(users.nickname) < LOWER($2) COLLATE "C" `
	//	} else {
	//		sqlRek += `AND LOWER(users.nickname) > LOWER($2) COLLATE "C" `
	//	}
	//	values = append(values, since)
	//	i++
	//}
	//sqlRek += `ORDER BY LOWER(users.nickname) COLLATE "C" `
	//if desc {
	//	sqlRek += `DESC `
	//}
	//if limit != 0{
	//	sqlRek += `LIMIT $` + fmt.Sprint(i) + `;`
	//	values = append(values, limit)
	//}else{
	//	sqlRek += `;`
	//}
	//rows, err := models.DB.Query(sqlRek, values...)
	//
	//if err != nil {
	//	return nil, err
	//}
	//defer rows.Close()
	//user := models.User{}
	//for rows.Next() {
	//	err = rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	//	if err != nil {
	//		return nil, err
	//	}
	//	users = append(users, user)
	//}
	////if len(users) == 0{
	////	users := models.Users{}
	////	user := models.User{}
	////	var usernames []string
	////
	////	rows, err := models.DB.Query("SELECT author FROM threads WHERE forum ILIKE $1 UNION SELECT author FROM posts WHERE forum ILIKE $1;", slug)
	////	if err != nil {
	////		return nil, err
	////	}
	////	defer rows.Close()
	////
	////	var u string
	////	for rows.Next() {
	////		err := rows.Scan(&u)
	////		if err != nil {
	////			return nil, err
	////		}
	////		usernames = append(usernames, u)
	////	}
	////
	////	values := make([]interface{}, 0, 3)
	////	sqlRek := `SELECT about, email, fullname, nickname FROM users WHERE nickname = ANY($1) `
	////	values = append(values, pq.Array(usernames))
	////	i := 2
	////	if since != "" {
	////		if desc {
	////			sqlRek += `AND LOWER(nickname) < LOWER($2) COLLATE "C" `
	////		} else {
	////			sqlRek += `AND LOWER(nickname) > LOWER($2) COLLATE "C" `
	////		}
	////		values = append(values, since)
	////		i++
	////	}
	////	sqlRek += `ORDER BY LOWER(nickname) COLLATE "C" `
	////	if desc {
	////		sqlRek += `DESC `
	////	}
	////	sqlRek += `LIMIT $` + fmt.Sprint(i) + `;`
	////	values = append(values, limit)
	////
	////	rows, err = models.DB.Query(sqlRek, values...)
	////
	////	if err != nil {
	////		return nil, err
	////	}
	////	defer rows.Close()
	////
	////	for rows.Next() {
	////		err = rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	////		if err != nil {
	////			return nil, err
	////		}
	////		users = append(users, user)
	////	}
	////
	////	return &users, nil
	////}
	//return &users, nil

	users := models.Users{}
	user := models.User{}
	var usernames []string

	rows, err := models.DB.Query("SELECT author FROM threads WHERE forum ILIKE $1 UNION SELECT author FROM posts WHERE forum ILIKE $1;", slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u string
	for rows.Next() {
		err := rows.Scan(&u)
		if err != nil {
			return nil, err
		}
		usernames = append(usernames, u)
	}

	values := make([]interface{}, 0, 3)
	sqlRek := `SELECT about, email, fullname, nickname FROM users WHERE nickname = ANY($1) `
	values = append(values, pq.Array(usernames))
	i := 2
	if since != "" {
		if desc {
			sqlRek += `AND LOWER(nickname) < LOWER($2) COLLATE "C" `
		} else {
			sqlRek += `AND LOWER(nickname) > LOWER($2) COLLATE "C" `
		}
		values = append(values, since)
		i++
	}
	sqlRek += `ORDER BY LOWER(nickname) COLLATE "C" `
	if desc {
		sqlRek += `DESC `
	}
	sqlRek += `LIMIT $` + fmt.Sprint(i) + `;`
	values = append(values, limit)

	rows, err = models.DB.Query(sqlRek, values...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}
