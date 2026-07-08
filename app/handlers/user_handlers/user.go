package user_handlers

import (
	"encoding/json"
	"go-forum-backend/app/actions/user_actions"
	"go-forum-backend/app/models/user_models"
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

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	page, limit := parsePage(r)
	users := user_models.GetAllUsers(page, limit)
	response.NewSuccessData(w, users, "")
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := r.PathValue("id")

	user := user_models.GetUserByID(id)
	if user == nil {
		response.NewErrorMessage(w, response.ErrUserNotFound, http.StatusNotFound)
		return
	}

	response.NewSuccessData(w, user, "")
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := r.PathValue("id")

	if user_models.GetUserByID(id) == nil {
		response.NewErrorMessage(w, response.ErrUserNotFound, http.StatusNotFound)
		return
	}

	var dto user_actions.UpdateUserActionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, user := user_actions.UpdateUser(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, user, "")
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := r.PathValue("id")

	if user_models.GetUserByID(id) == nil {
		response.NewErrorMessage(w, response.ErrUserNotFound, http.StatusNotFound)
		return
	}

	user_models.DeleteUser(id)
	response.NewSuccessMessage(w, "User deleted")
}

func GetUserMessagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := r.PathValue("id")

	if user_models.GetUserByID(id) == nil {
		response.NewErrorMessage(w, response.ErrUserNotFound, http.StatusNotFound)
		return
	}

	messages := user_models.GetUserMessages(id)
	response.NewSuccessData(w, messages, "")
}
