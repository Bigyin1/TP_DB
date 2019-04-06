package db

import (
	"database/sql"
	"fmt"
	"gohw/internal/models"
)

func (db *Database) InsertOrUpdateVote(vote models.Vote, thread *models.Thread) (err error) {

	var tx *sql.Tx
	tx, err = db.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sqlQuery := `SELECT voice
				FROM votes
				WHERE nickname=$1 and thread=$2`

	row := tx.QueryRow(sqlQuery, vote.Nickname, thread.ID)

	oldVoice := 0
	if err = row.Scan(&oldVoice); err != nil && err != sql.ErrNoRows {
		return
	}

	if err == sql.ErrNoRows {
		sqlQuery := `INSERT INTO votes(nickname, voice, thread)
					VALUES($1, $2, $3);`
		if _, err = tx.Exec(sqlQuery, vote.Nickname, vote.Voice, thread.ID); err != nil {
			return
		}
	} else {
		sqlQuery := `UPDATE votes
					SET voice = $1
					WHERE nickname = $2 AND thread = $3;`
		if _, err = tx.Exec(sqlQuery, vote.Voice, vote.Nickname, thread.ID); err != nil {
			return
		}
	}

	thread.Votes += (vote.Voice - oldVoice)
	sqlUpdate := `UPDATE threads
				SET votes = $1
				WHERE id = $2;`
	if _, err = tx.Exec(sqlUpdate, thread.Votes, thread.ID); err != nil {
		return
	}
	tx.Commit()
	return
}
