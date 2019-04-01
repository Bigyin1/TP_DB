package api

import (
	db "gohw/internal/database"
)

type Handler struct {
	db *db.Database
}

func Init(DB *db.Database) (handler *Handler) {
	handler = &Handler{
		db: DB,
	}
	return
}

func GetHandler() (handler *Handler, err error) {

	var DB *db.Database
	if DB, err = db.Init(); err != nil {
		return
	}
	handler = Init(DB)
	return
}
