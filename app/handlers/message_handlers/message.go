package message_handlers

import (
	"encoding/json"
	"go-forum-backend/app/actions/message_actions"
	"go-forum-backend/app/models/message_models"
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

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	page, limit := parsePage(r)
	messages := message_models.GetAllMessages(page, limit)
	response.NewSuccessData(w, messages, "")
}

func CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var dto message_actions.MessageDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, message := message_actions.CreateMessage(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, message, "")
}

func GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	message := message_models.GetMessageByID(id)
	if message == nil {
		response.NewErrorMessage(w, response.ErrMessageNotFound, http.StatusNotFound)
		return
	}

	response.NewSuccessData(w, message, "")
}

func UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if message_models.GetMessageByID(id) == nil {
		response.NewErrorMessage(w, response.ErrMessageNotFound, http.StatusNotFound)
		return
	}

	var dto message_actions.UpdateMessageDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, message := message_actions.UpdateMessage(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, message, "")
}

func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if message_models.GetMessageByID(id) == nil {
		response.NewErrorMessage(w, response.ErrMessageNotFound, http.StatusNotFound)
		return
	}

	message_models.DeleteMessage(id)
	response.NewSuccessMessage(w, "Message deleted")
}

func GetMessageUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if message_models.GetMessageByID(id) == nil {
		response.NewErrorMessage(w, response.ErrMessageNotFound, http.StatusNotFound)
		return
	}

	users := message_models.GetMessageUsers(id)
	response.NewSuccessData(w, users, "")
}

func LinkMessageUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if message_models.GetMessageByID(id) == nil {
		response.NewErrorMessage(w, response.ErrMessageNotFound, http.StatusNotFound)
		return
	}

	var body struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	message_models.LinkUser(id, body.UserID)
	response.NewSuccessMessage(w, "User assigned")
}

func UnlinkMessageUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	userID := r.PathValue("user_id")
	message_models.UnlinkUser(id, userID)
	response.NewSuccessMessage(w, "User removed")
}
