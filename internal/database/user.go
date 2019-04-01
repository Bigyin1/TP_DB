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
		fmt.Printf("GetAllCollisionUsers error : %s\n", err.Error())
		return
	}

	for rows.Next() {
		u := &models.User{}
		if err = rows.Scan(&u.Nickname, &u.Fullname, &u.About, &u.Email); err != nil {
			fmt.Printf("GetAllCollisionUsers error : %s\n", err.Error())
			return
		}
		usrs = append(usrs, u)
	}
	rows.Close()
	return
}
