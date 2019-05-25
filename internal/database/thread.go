package db

import (
	"database/sql"
	"errors"
	"fmt"
	"gohw/internal/models"
	rerrors "gohw/internal/return_errors"
)

func (db *Database) GetThreadBySlug(slug string) (thread models.Thread, err error) {

	var tx *sql.Tx
	var slugBuf sql.NullString
	tx, err = db.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("GetThreadBySlug start")
	sqlQuery := `SELECT id, created, slug, message, title, votes, forum, author
				FROM threads
				WHERE slug=$1;`

	row, err := tx.Query(sqlQuery, slug)
	if err != nil {
		return
	}

	row.Next()
	err = row.Scan(&thread.ID, &thread.Created, &slugBuf, &thread.Message,
		&thread.Title, &thread.Votes, &thread.Forum, &thread.Author)
	thread.Slug = slugBuf.String

	return
}

func (db *Database) GetThreadByID(id int) (thread models.Thread, err error) {

	var tx *sql.Tx
	var slugBuf sql.NullString
	tx, err = db.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("GetThreadByID start")
	sqlQuery := `SELECT id, created, slug, message, title, votes, forum, author
				FROM threads
				WHERE id=$1;`

	row, err := tx.Query(sqlQuery, id)
	if err != nil {
		fmt.Println("GetThreadByID ", err.Error())
		return
	}

	row.Next()
	err = row.Scan(&thread.ID, &thread.Created, &slugBuf, &thread.Message,
		&thread.Title, &thread.Votes, &thread.Forum, &thread.Author)
	thread.Slug = slugBuf.String
	if err != nil {
		fmt.Println("GetThreadByID ", err.Error())
	}
	return
}

func (db *Database) CreateThread(thread *models.Thread) (err error) {

	var tx *sql.Tx
	tx, err = db.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var (
		forum models.Forum
	)
	if _, err = db.GetUserByName(thread.Author); err != nil {
		fmt.Println(err.Error())
		return errors.New(rerrors.UserNotFound)
	}
	if forum, err = db.GetForumBySlug(thread.Forum); err != nil {
		fmt.Println(err.Error())
		return errors.New(rerrors.UserNotFound)
	}

	thread.Forum = forum.Slug

	if thread.Created != "" {
		sqlInsert := `INSERT INTO threads (author, forum, message, slug, title, created)
				VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6) RETURNING id, forum;`

		row := tx.QueryRow(sqlInsert, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Title, thread.Created)
		if err = row.Scan(&thread.ID, &thread.Forum); err != nil {
			fmt.Printf("CreateThread error: %s\n", err.Error())
			t, e := db.GetThreadBySlug(thread.Slug)
			if e != nil {
				fmt.Printf("CreateThread error: %s\n", e.Error())
			}
			*thread = t
			return errors.New(rerrors.AlreadyExist)
		}
	} else {
		sqlInsert := `INSERT INTO threads (author, forum, message, slug, title)
				VALUES ($1, $2, $3, NULLIF($4, ''), $5) RETURNING id, forum;`

		row := tx.QueryRow(sqlInsert, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Title)
		if err = row.Scan(&thread.ID, &thread.Forum); err != nil {
			fmt.Printf("CreateThread error: %s\n", err.Error())
			t, e := db.GetThreadBySlug(thread.Slug)
			if e != nil {
				fmt.Printf("CreateThread error: %s\n", e.Error())
			}
			*thread = t
			return errors.New(rerrors.AlreadyExist)
		}
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	return errors.New("No Error")
}

func (db *Database) GetForumThreads(query models.URLQuery) (threads models.Threads, err error) {

	var tx *sql.Tx
	tx, err = db.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sqlQuery := `SELECT id, title, author, forum, message, votes, slug, created
				FROM threads t
				WHERE forum=$1`

	if query.Since != "" {
		if query.Desc {
			sqlQuery += fmt.Sprintf(" AND created <= '%s'", query.Since)
		} else {
			sqlQuery += fmt.Sprintf(" AND created >= '%s'", query.Since)
		}
	}
	if query.Desc {
		sqlQuery += ` ORDER BY created DESC`
	} else {
		sqlQuery += ` ORDER BY created ASC`
	}

	sqlQuery += ` LIMIT $2;`

	rows, err := tx.Query(sqlQuery, query.Slug, query.Limit)
	if err != nil {
		return
	}

	threads = make(models.Threads, 0)
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

func (db *Database) CountThreads() (count int, err error) {

	var tx *sql.Tx
	tx, err = db.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sqlQuery := `SELECT COUNT(*) FROM threads;`

	row := tx.QueryRow(sqlQuery)

	if err = row.Scan(&count); err != nil {
		fmt.Printf("CountThreads error: %s", err.Error())
	}
	return
}

func (db *Database) UpdateThread(old *models.Thread, new *models.Thread) (err error) {

	var tx *sql.Tx
	tx, err = db.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sqlQuery := `UPDATE threads
				SET title    = COALESCE(NULLIF($1, ''), title),
				message    = COALESCE(NULLIF($2, ''), message)
				WHERE id = $3
				RETURNING title, message;`

	row := tx.QueryRow(sqlQuery, new.Title, new.Message, old.ID)

	if err = row.Scan(&old.Title, &old.Message); err != nil {
		fmt.Printf("UpdateThreads error: %s\n", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	return
}

func (db *Database) GetThreadPosts(posts *models.Posts, query models.URLQuery) (err error) {

	var tx *sql.Tx
	tx, err = db.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sqlQuery := `SELECT id, parent, author, message, is_edited, forum, thread, created
				FROM posts t
				WHERE`

	sign := ">"

	if query.Since != "" {
		if query.Desc {
			sign = "<"
		}
	}
	switch query.Sort {
	case "tree":
		sqlQuery += ` thread = $1`
		if query.Since != "" {
			sqlQuery += fmt.Sprintf(" AND path %s (SELECT path from posts where id = '%s')", sign, query.Since)
		}
		if query.Desc {
			sqlQuery += ` ORDER BY path DESC`
		} else {
			sqlQuery += ` ORDER BY path ASC`
		}

		sqlQuery += ` LIMIT $2;`
	case "parent_tree":
		sqlQuery += ` branch in (SELECT p.id FROM posts as p WHERE p.thread = $1 AND parent = 0`
		if query.Since != "" {
			sqlQuery += fmt.Sprintf(" and p.id %s (SELECT branch FROM Posts WHERE id = '%s')", sign, query.Since)
		}
		sqlQuery += ` order by p.id`
		if query.Desc {
			sqlQuery += ` desc`
		}
		sqlQuery += ` LIMIT $2`
		if query.Desc {
			sqlQuery += " ) ORDER by t.branch desc, t.path;"
		} else {
			sqlQuery += " ) ORDER by t.path;"
		}
	case "flat":
		sqlQuery += ` thread = $1`
		if query.Since != "" {
			sqlQuery += fmt.Sprintf(" AND id %s %s", sign, query.Since)
		}
		if query.Desc {
			sqlQuery += ` ORDER BY created DESC, id DESC`
		} else {
			sqlQuery += ` ORDER BY created ASC, id ASC`
		}

		sqlQuery += ` LIMIT $2;`

	}

	rows, err := tx.Query(sqlQuery, query.ID, query.Limit)
	if err != nil {
		fmt.Printf("GetThreadsPosts error: %s\n", err.Error())
		return
	}
	*posts = make(models.Posts, 0)
	for rows.Next() {
		post := &models.Post{}
		if err = rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
			return
		}
		*posts = append(*posts, post)
	}
	rows.Close()

	return
}
