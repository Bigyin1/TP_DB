package db

import (
	"fmt"
	"gohw/internal/models"
)

func (db *Database) GetPostByID(id int) (post models.Post, err error) {

	sqlQuery := `SELECT p.author, p.created, p.forum, p.is_edited, p.id, p.message, p.parent, p.thread
				FROM posts as p where p.id = $1;`

	row := db.DB.QueryRow(sqlQuery, id)
	err = row.Scan(&post.Author, &post.Created, &post.Forum,
		&post.IsEdited, &post.ID, &post.Message, &post.Parent, &post.Thread)

	return
}

func (db *Database) UpdatePost(post *models.Post) (err error) {

	sqlQuery := `UPDATE posts
				SET message = $1,
					is_edited = $2
				WHERE posts.id = $3;
				`
	_, err = db.DB.Query(sqlQuery, post.Message, post.IsEdited, post.ID)
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
