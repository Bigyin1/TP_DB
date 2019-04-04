package db

import (
	"fmt"
	"gohw/internal/models"
)

func (db *Database) InsertNewUser(user *models.User) (err error) {

	sqlInsert := `INSERT INTO users(nickname, fullname, about, email)
	VALUES ( $1, $2, $3, $4 );`
	_, err = db.DB.Exec(sqlInsert, user.Nickname, user.Fullname, user.About, user.Email)

	return
}

func (db *Database) GetAllCollisionUsers(user *models.User) (usrs models.Users, err error) {
	sqlQuery := `SELECT nickname, fullname, about, email
				FROM users u
				WHERE u.nickname=$1 OR u.email=$2;`
	rows, err := db.DB.Query(sqlQuery, user.Nickname, user.Email)
	if err != nil {
		return
	}

	for rows.Next() {
		u := &models.User{}
		if err = rows.Scan(&u.Nickname, &u.Fullname, &u.About, &u.Email); err != nil {
			return
		}
		usrs = append(usrs, u)
	}
	rows.Close()
	return
}

func (db *Database) GetUserByName(name string) (user models.User, err error) {

	sqlQuery := `SELECT nickname, fullname, about, email
				FROM users u
				WHERE nickname=$1;`

	row, err := db.DB.Query(sqlQuery, name)
	if err != nil {
		return
	}

	row.Next()
	if err = row.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
		return
	}
	return
}

func (db *Database) UpdateProfile(user *models.User) (err error) {

	sqlQuery := `UPDATE users
	SET
		email    = COALESCE(NULLIF($1, ''), email),
		about    = COALESCE(NULLIF($2, ''), about),
		fullname = COALESCE(NULLIF($3, ''), fullname)
	WHERE nickname=$4
	RETURNING fullname, about, email;`

	row, err := db.DB.Query(sqlQuery, user.Email, user.About, user.Fullname, user.Nickname)
	if err != nil {
		return
	}
	row.Next()
	if err = row.Scan(&user.Fullname, &user.About, &user.Email); err != nil {
		return
	}
	return
}

func (db *Database) GetUsersByForum(query models.ForumUsersQuery) (users models.Users, err error) {

	sqlQuery := `SELECT DISTINCT u.nickname, u.fullname, u.about, u.email
				FROM forums f
				JOIN threads t on f.slug=t.forum
				JOIN users u on t.author=u.nickname
				WHERE f.slug like $1`

	if query.Desc {
		sqlQuery += ` AND u.nickname < $2`
	} else {
		sqlQuery += ` AND u.nickname > $2`
	}
	sqlQuery += ` UNION DISTINCT SELECT DISTINCT u.nickname, u.fullname, u.about, u.email
	FROM forums f
	JOIN posts p on f.slug=p.forum
	JOIN users u on p.author=u.nickname
	WHERE f.slug like $1`

	if query.Desc {
		sqlQuery += ` AND u.nickname < $2`
	} else {
		sqlQuery += ` AND u.nickname > $2`
	}

	if query.Desc {
		sqlQuery += ` ORDER BY 1 DESC`
	} else {
		sqlQuery += ` ORDER BY 1 ASC`
	}

	sqlQuery += ` LIMIT $3;`

	fmt.Println(sqlQuery)
	rows, err := db.DB.Query(sqlQuery, query.Slug, query.Since, query.Limit)
	if err != nil {
		return
	}

	for rows.Next() {
		u := &models.User{}
		if err = rows.Scan(&u.Nickname, &u.Fullname, &u.About, &u.Email); err != nil {
			return
		}
		users = append(users, u)
	}
	rows.Close()

	return
}
