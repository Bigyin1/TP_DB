package db

import (
	"database/sql"
	"fmt"
	"gohw/internal/models"
)

func (db *Database) GetPostByID(id int) (post models.Post, err error) {

	fmt.Println("GetPostByID start")
	var tx *sql.Tx
	tx, err = db.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sqlQuery := `SELECT p.author, p.created, p.forum, p.is_edited, p.id, p.message, p.parent, p.thread
				FROM posts as p where p.id = $1;`

	row := tx.QueryRow(sqlQuery, id)
	err = row.Scan(&post.Author, &post.Created, &post.Forum,
		&post.IsEdited, &post.ID, &post.Message, &post.Parent, &post.Thread)

	return
}

func (db *Database) UpdatePost(post *models.Post) (err error) {

	fmt.Println("UpdatePost start")
	var tx *sql.Tx
	tx, err = db.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sqlQuery := `UPDATE posts
				SET message = COALESCE(NULLIF($1, ''), message),
					is_edited = $2
				WHERE posts.id = $3;
				`
	_, err = tx.Exec(sqlQuery, post.Message, post.IsEdited, post.ID)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = tx.Commit()
	return
}

func (db *Database) CountPosts() (count int, err error) {

	sqlQuery := `SELECT COUNT(*) FROM posts;`

	row := db.DB.QueryRow(sqlQuery)

	if err = row.Scan(&count); err != nil {
		fmt.Printf("CountPosts error: %s", err.Error())
	}
	return
}

func (db *Database) CreatePost(post *models.Post, tx *sql.Tx) (err error) {

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sqlQuery := `INSERT into posts(author, forum, message, parent, thread, created)
				VALUES($1, $2, $3, $4, $5, $6)
				RETURNING id;`

	row := tx.QueryRow(sqlQuery, post.Author, post.Forum, post.Message, post.Parent, post.Thread, post.Created)

	sqlForumUsers := `INSERT INTO UsersForum(forum, userNickname)
					 VALUES ($1, $2)
					 ON CONFLICT (forum, userNickname) DO NOTHING;`

	if err = row.Scan(&post.ID); err != nil {
		fmt.Printf("CreatePost error: %s", err.Error())
		return
	}

	_, err = tx.Exec(sqlForumUsers, post.Forum, post.Author)
	if err != nil {
		fmt.Printf("CreatePost error: %s\n", err.Error())
	}
	return
}
