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
		fmt.Println(err.Error())
		return errors.New(rerrors.UserNotFound)
	}
	if _, err = db.GetForumBySlug(thread.Forum); err != nil {
		fmt.Println(err.Error())
		return errors.New(rerrors.UserNotFound)
	}

	if thread.Created != "" {
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
	} else {
		sqlInsert := `INSERT INTO threads(author, forum, message, slug, title)
				VALUES ($1, $2, $3, $4, $5);`

		if _, err = db.DB.Exec(sqlInsert, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Title); err != nil {
			t, e := db.GetThreadBySlug(thread.Slug)
			if e != nil {
				fmt.Printf("CreateThread error: %s\n", e.Error())
			}
			*thread = t
			return errors.New(rerrors.AlreadyExist)
		}
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

	sqlQuery := `SELECT COUNT(*) FROM threads;`

	row := db.DB.QueryRow(sqlQuery)

	if err = row.Scan(&count); err != nil {
		fmt.Printf("CountThreads error: %s", err.Error())
	}
	return
}

func (db *Database) UpdateThread(old *models.Thread, new *models.Thread) (err error) {

	sqlQuery := `UPDATE threads
				SET title    = COALESCE(NULLIF($1, ''), title),
				message    = COALESCE(NULLIF($2, ''), message)
				WHERE id = $3
				RETURNING title, message;`

	row := db.DB.QueryRow(sqlQuery, new.Title, new.Message, old.ID)

	if err = row.Scan(&old.Title, &old.Message); err != nil {
		fmt.Printf("UpdateThreads error: %s\n", err.Error())
	}

	return
}

func (db *Database) GetThreadPosts(posts *models.Posts, query models.URLQuery) (err error) {

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
			sqlQuery += fmt.Sprintf(" AND path %s (SELECT path from posts where id = %s", sign, query.Since)
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
			sqlQuery += ` ORDER BY created DESC`
		} else {
			sqlQuery += ` ORDER BY created ASC`
		}

		sqlQuery += ` LIMIT $2;`

	}

	rows, err := db.DB.Query(sqlQuery, query.ID, query.Limit)
	if err != nil {
		fmt.Printf("GetThreadsPosts error: %s\n", err.Error())
		return
	}
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
