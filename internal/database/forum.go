package db

import (
	"errors"
	"fmt"
	"gohw/internal/models"
	rerrors "gohw/internal/return_errors"
)

func (db *Database) GetForumBySlug(slug string) (forum models.Forum, err error) {

	sqlQuery := `SELECT owner, title, slug, threads, posts
				FROM forums f
				WHERE f.slug=$1`

	row, err := db.DB.Query(sqlQuery, slug)
	if err != nil {
		return
	}

	row.Next()
	err = row.Scan(&forum.Author, &forum.Title, &forum.Slug, &forum.Threads, &forum.Posts)
	return
}

func (db *Database) CreateForum(forum *models.Forum) (err error) {

	if _, err = db.GetUserByName(forum.Author); err != nil {
		return errors.New(rerrors.UserNotFound)
	}
	sqlInsert := `INSERT INTO forums(owner, title, slug)
				VALUES ($1, $2, $3);`

	if _, err = db.DB.Exec(sqlInsert, forum.Author, forum.Title, forum.Slug); err != nil {
		f, e := db.GetForumBySlug(forum.Slug)
		if e != nil {
			fmt.Printf("CreateForum error: %s", e.Error())
		}
		*forum = f
		return errors.New(rerrors.AlreadyExist)
	}
	return errors.New("No Error")
}

func (db *Database) CountForums() (count int, err error) {

	sqlQuery := `SELECT COUNT(*) FROM forums;`

	row := db.DB.QueryRow(sqlQuery)

	if err = row.Scan(&count); err != nil {
		fmt.Printf("CountForums error: %s", err.Error())
	}
	return
}
