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

func (db *Database) GetThreadByID(id int) (thread models.Thread, err error) {

	sqlQuery := `SELECT id, created, slug, message, title, votes, forum, author
				FROM threads
				WHERE id=$1;`

	row, err := db.DB.Query(sqlQuery, id)
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

func (db *Database) GetForumThreads(query models.URLQuery) (threads models.Threads, err error) {

	sqlQuery := `SELECT id, title, author, forum, message, votes, slug, created
				FROM threads t
				WHERE forum=$1`

	if query.Since != "" {
		if query.Desc {
			sqlQuery += ` AND t.created < ` + query.Since
		} else {
			sqlQuery += ` AND t.created > ` + query.Since
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

	for rows.Next() {
		t := &models.Thread{}
		if err = rows.Scan(&t.ID, &t.Title, &t.Author, &t.Forum,
			&t.Message, &t.Votes, &t.Slug, &t.Created); err != nil {
			return
		}
		threads = append(threads, t)
	}
	rows.Close()
	return
}
