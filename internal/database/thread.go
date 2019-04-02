package db

import (
	"errors"
	"fmt"
	"gohw/internal/models"
	rerrors "gohw/internal/return_errors"
)

func (db *Database) GetThreadBySlug(slug string) (thread models.Thread, err error) {

	sqlQuery := `SELECT id, created, slug, message, title, votes, forum, author
				FROM threads
				WHERE slug=$1;`

	row, err := db.DB.Query(sqlQuery, slug)
	if err != nil {
		return
	}

	row.Next()
	err = row.Scan(&thread.ID, &thread.Created, &thread.Slug, &thread.Message,
		&thread.Title, &thread.Votes, &thread.Forum, &thread.Author)

	return
}

func (db *Database) CreateThread(thread *models.Thread) (err error) {

	if _, err = db.GetUserByName(thread.Author); err != nil {
		return errors.New(rerrors.UserNotFound)
	}
	if _, err = db.GetForumBySlug(thread.Forum); err != nil {
		return errors.New(rerrors.UserNotFound)
	}

	sqlInsert := `INSERT INTO threads(author, forum, message, slug, title, created)
				VALUES ($1, $2, $3, $4, $5, $6);`

	if _, err = db.DB.Exec(sqlInsert, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Title, thread.Created); err != nil {
		t, e := db.GetThreadBySlug(thread.Slug)
		if e != nil {
			fmt.Printf("CreateThread error: %s\n", e.Error())
		}
		*thread = t
		return errors.New(rerrors.AlreadyExist)
	}

	t, _ := db.GetThreadBySlug(thread.Slug)
	*thread = t
	return errors.New("No Error")
}
