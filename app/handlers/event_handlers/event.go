package event_handlers

import (
	"encoding/json"
	"go-forum-backend/app/actions/event_actions"
	"go-forum-backend/app/models/event_models"
	"go-forum-backend/utils/log"
	"go-forum-backend/utils/response"
	"net/http"
	"strconv"
)

func parsePage(r *http.Request) (int, int) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 20
	}
	return page, limit
}

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	page, limit := parsePage(r)
	events := event_models.GetAllEvents(page, limit)
	response.NewSuccessData(w, events, "")
}

func GetEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	event := event_models.GetEventByID(id)
	if event == nil {
		response.NewErrorMessage(w, response.ErrEventNotFound, http.StatusNotFound)
		return
	}

	response.NewSuccessData(w, event, "")
}

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var dto event_actions.CreateEventActionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, event := event_actions.CreateEvent(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, event, "")
}

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if event_models.GetEventByID(id) == nil {
		response.NewErrorMessage(w, response.ErrEventNotFound, http.StatusNotFound)
		return
	}

	var dto event_actions.UpdateEventActionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, event := event_actions.UpdateEvent(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, event, "")
}

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if event_models.GetEventByID(id) == nil {
		response.NewErrorMessage(w, response.ErrEventNotFound, http.StatusNotFound)
		return
	}

	event_models.DeleteEvent(id)
	response.NewSuccessMessage(w, "Event deleted")
}

func GetEventUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if event_models.GetEventByID(id) == nil {
		response.NewErrorMessage(w, response.ErrEventNotFound, http.StatusNotFound)
		return
	}

	users := event_models.GetEventUsers(id)
	response.NewSuccessData(w, users, "")
}

func LinkEventUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if event_models.GetEventByID(id) == nil {
		response.NewErrorMessage(w, response.ErrEventNotFound, http.StatusNotFound)
		return
	}

	var body struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	event_models.LinkUser(id, body.UserID)
	response.NewSuccessMessage(w, "User assigned")
}

func UnlinkEventUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	userID := r.PathValue("user_id")
	event_models.UnlinkUser(id, userID)
	response.NewSuccessMessage(w, "User removed")
}
