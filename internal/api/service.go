package api

import (
	"fmt"
	"gohw/internal/models"
	"net/http"
)

func (h *Handler) ServiceClear(rw http.ResponseWriter, r *http.Request) {

	if err := h.db.ServiceClear(); err != nil {
		fmt.Printf("ServiceClear error: %s", err.Error())
		response(rw, http.StatusInternalServerError, nil)
		return
	}
	response(rw, http.StatusOK, nil)

}

func (h *Handler) ServiceInfo(rw http.ResponseWriter, r *http.Request) {
	var (
		service models.Service
		err     error
	)

	if service.User, err = h.db.CountUsers(); err != nil {
		response(rw, http.StatusInternalServerError, nil)
		return
	}
	if service.Post, err = h.db.CountPosts(); err != nil {
		response(rw, http.StatusInternalServerError, nil)
		return
	}
	if service.Forum, err = h.db.CountForums(); err != nil {
		response(rw, http.StatusInternalServerError, nil)
		return
	}
	if service.Thread, err = h.db.CountThreads(); err != nil {
		response(rw, http.StatusInternalServerError, nil)
		return
	}
	response(rw, http.StatusOK, service)

}
