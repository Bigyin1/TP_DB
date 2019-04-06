package db

import (
	"database/sql"
	"errors"
	"fmt"
	"gohw/internal/models"
	rerrors "gohw/internal/return_errors"
)

func (db *Database) GetForumBySlug(slug string) (forum models.Forum, err error) {

	sqlQuery := `SELECT owner, title, slug, threads, posts
				FROM forums f
				WHERE lower(f.slug)=lower($1);`

	row := db.DB.QueryRow(sqlQuery, slug)

	err = row.Scan(&forum.Author, &forum.Title, &forum.Slug, &forum.Threads, &forum.Posts)
	if err != nil {
		fmt.Printf("GetForumBySlug error: %s\n", err.Error())
		return
	}
	return
}

func (db *Database) CreateForum(forum *models.Forum) (err error) {

	var tx *sql.Tx
	tx, err = db.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var us models.User

	if us, err = db.GetUserByName(forum.Author); err != nil {
		return errors.New(rerrors.UserNotFound)
	}
	forum.Author = us.Nickname
	sqlInsert := `INSERT INTO forums(owner, title, slug)
				VALUES ($1, $2, $3);`

	if _, err = tx.Exec(sqlInsert, forum.Author, forum.Title, forum.Slug); err != nil {
		fmt.Printf("CreateForum err %s\n", err.Error())
		f, e := db.GetForumBySlug(forum.Slug)

		if e != nil {
			return e
		}
		*forum = f
		return errors.New(rerrors.AlreadyExist)
	}
	tx.Commit()
	return errors.New(rerrors.Ok)
}

func (db *Database) CountForums() (count int, err error) {

	sqlQuery := `SELECT COUNT(*) FROM forums;`

	row := db.DB.QueryRow(sqlQuery)

	if err = row.Scan(&count); err != nil {
		fmt.Printf("CountForums error: %s", err.Error())
	}
	return
}
