package db

func (db *Database) ServiceClear() (err error) {

	sqlRow := `TRUNCATE posts, users, threads, forums, votes CASCADE;`

	_, err = db.DB.Exec(sqlRow)

	return
}
