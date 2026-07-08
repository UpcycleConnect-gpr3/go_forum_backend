package talk_handlers

import (
	"encoding/json"
	"go-forum-backend/app/actions/talk_actions"
	"go-forum-backend/app/models/talk_models"
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

func GetTalksHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	page, limit := parsePage(r)
	talks := talk_models.GetAllTalks(page, limit)
	response.NewSuccessData(w, talks, "")
}

func GetTalkHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	talk := talk_models.GetTalkByID(id)
	if talk == nil {
		response.NewErrorMessage(w, response.ErrTalkNotFound, http.StatusNotFound)
		return
	}

	response.NewSuccessData(w, talk, "")
}

func CreateTalkHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var dto talk_actions.CreateTalkDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, talk := talk_actions.CreateTalk(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, talk, "")
}

func UpdateTalkHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if talk_models.GetTalkByID(id) == nil {
		response.NewErrorMessage(w, response.ErrTalkNotFound, http.StatusNotFound)
		return
	}

	var dto talk_actions.UpdateTalkDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, talk := talk_actions.UpdateTalk(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, talk, "")
}

func DeleteTalkHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if talk_models.GetTalkByID(id) == nil {
		response.NewErrorMessage(w, response.ErrTalkNotFound, http.StatusNotFound)
		return
	}

	talk_models.DeleteTalk(id)
	response.NewSuccessMessage(w, "Talk deleted")
}

func GetTalkMessagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if talk_models.GetTalkByID(id) == nil {
		response.NewErrorMessage(w, response.ErrTalkNotFound, http.StatusNotFound)
		return
	}

	messages := talk_models.GetTalkMessages(id)
	response.NewSuccessData(w, messages, "")
}

func LinkTalkMessageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if talk_models.GetTalkByID(id) == nil {
		response.NewErrorMessage(w, response.ErrTalkNotFound, http.StatusNotFound)
		return
	}

	var body struct {
		MessageID int `json:"message_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	talk_models.LinkMessage(id, body.MessageID)
	response.NewSuccessMessage(w, "Message assigned")
}

func UnlinkTalkMessageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	messageID, err := strconv.Atoi(r.PathValue("message_id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid message_id", http.StatusBadRequest)
		return
	}

	talk_models.UnlinkMessage(id, messageID)
	response.NewSuccessMessage(w, "Message removed")
}

func GetTalkUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if talk_models.GetTalkByID(id) == nil {
		response.NewErrorMessage(w, response.ErrTalkNotFound, http.StatusNotFound)
		return
	}

	users := talk_models.GetTalkUsers(id)
	response.NewSuccessData(w, users, "")
}

func LinkTalkUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if talk_models.GetTalkByID(id) == nil {
		response.NewErrorMessage(w, response.ErrTalkNotFound, http.StatusNotFound)
		return
	}

	var body struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	talk_models.LinkUser(id, body.UserID)
	response.NewSuccessMessage(w, "User assigned")
}

func UnlinkTalkUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	userID := r.PathValue("user_id")
	talk_models.UnlinkUser(id, userID)
	response.NewSuccessMessage(w, "User removed")
}
