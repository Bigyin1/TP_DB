package api

import (
	"fmt"
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
