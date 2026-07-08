package event_handlers

import (
	"net/http"
	"strconv"

	"go-forum-backend/app/models/event_models"
	"go-forum-backend/utils/jwt"
	"go-forum-backend/utils/log"
	"go-forum-backend/utils/response"
)

func setEventStatus(w http.ResponseWriter, r *http.Request, status string) {
	log.Api(r)

	if jwt.RoleFromToken(r.Header.Get("Authorization")) != "administrator" {
		response.NewErrorMessage(w, response.ErrForbidden, http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}
	if event_models.GetEventByID(id) == nil {
		response.NewErrorMessage(w, response.ErrEventNotFound, http.StatusNotFound)
		return
	}
	if err := event_models.SetStatus(id, status); err != nil {
		response.NewErrorMessage(w, response.ErrEventNotFound, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, map[string]any{"id": id, "status": status}, "")
}

func ValidateEventHandler(w http.ResponseWriter, r *http.Request) {
	setEventStatus(w, r, "validated")
}

func RejectEventHandler(w http.ResponseWriter, r *http.Request) {
	setEventStatus(w, r, "rejected")
}
