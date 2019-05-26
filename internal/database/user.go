package db

import (
	"database/sql"
	"fmt"
	"gohw/internal/models"
)

func (db *Database) InsertNewUser(user *models.User) (err error) {

	var tx *sql.Tx
	tx, err = db.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("InsertNewUser start")
	sqlInsert := `INSERT INTO users(nickname, fullname, about, email)
	VALUES ( $1, $2, $3, $4 );`
	if _, err = tx.Exec(sqlInsert, user.Nickname, user.Fullname, user.About, user.Email); err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	return
}

func (db *Database) GetAllCollisionUsers(user *models.User) (usrs models.Users, err error) {
	var tx *sql.Tx
	tx, err = db.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("GetAllCollisionUsers start")
	sqlQuery := `SELECT nickname, fullname, about, email
				FROM users u
				WHERE u.nickname=$1 OR u.email=$2;`
	rows, err := tx.Query(sqlQuery, user.Nickname, user.Email)

	if err != nil {
		fmt.Println(err.Error())
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

	fmt.Printf("GetUserByName start, username: %s\n", name)
	var tx *sql.Tx
	tx, err = db.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	sqlQuery := `SELECT nickname, fullname, about, email
				FROM users
				WHERE lower(nickname) like lower($1);`

	row := tx.QueryRow(sqlQuery, name)
	fmt.Printf("GetUserByName query done\n")
	if err = row.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
		fmt.Printf("GetUserByName error: %s\n", err.Error())
		return
	}
	// err = tx.Commit()
	// if err != nil {
	// 	return
	// }
	fmt.Printf("GetUserByName end\n")
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

func (db *Database) GetUsersByForum(query models.URLQuery) (users models.Users, err error) {

	sqlQuery := `SELECT u.nickname, u.fullname, u.about, u.email
				FROM users u
				JOIN UsersForum f on f.userNickname=u.nickname
				WHERE f.forum=$1`

	if query.Since != "" {
		if query.Desc {
			sqlQuery += fmt.Sprintf(" AND u.nickname < '%s'", query.Since)
		} else {
			sqlQuery += fmt.Sprintf(" AND u.nickname > '%s'", query.Since)
		}
	}

	if query.Desc {
		sqlQuery += ` ORDER BY 1 DESC`
	} else {
		sqlQuery += ` ORDER BY 1 ASC`
	}

	sqlQuery += ` LIMIT $2;`
	rows, err := db.DB.Query(sqlQuery, query.Slug, query.Limit)
	if err != nil {
		return
	}
	users = make(models.Users, 0)
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

func (db *Database) CountUsers() (count int, err error) {

	sqlQuery := `SELECT COUNT(*) FROM users;`

	row := db.DB.QueryRow(sqlQuery)

	if err = row.Scan(&count); err != nil {
		fmt.Printf("CountUsers error: %s", err.Error())
	}
	return
}
