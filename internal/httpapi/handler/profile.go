package handler

import (
	"github.com/rinatkh/homework_go_15/internal/httpapi/authctx"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Me возвращает identity, которую middleware уже проверил по access JWT.
func Me(w http.ResponseWriter, r *http.Request) {
	id, ok := authctx.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"id": id})
}
