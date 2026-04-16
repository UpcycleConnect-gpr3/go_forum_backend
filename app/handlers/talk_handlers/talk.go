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

func CreateTalkHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var dto talk_actions.TalkDTO
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

	page, limit := parsePage(r)
	messages := talk_models.GetTalkMessages(id, page, limit)
	response.NewSuccessData(w, messages, "")
}

func CreateTalkMessageHandler(w http.ResponseWriter, r *http.Request) {
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

	var dto talk_actions.MessageDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, msg := talk_actions.CreateTalkMessage(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, msg, "")
}

func GetTalkCategoriesHandler(w http.ResponseWriter, r *http.Request) {
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

	categories := talk_models.GetTalkCategories(id)
	response.NewSuccessData(w, categories, "")
}

func LinkTalkCategoryHandler(w http.ResponseWriter, r *http.Request) {
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
		CategoryID int `json:"category_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	talk_models.LinkCategory(id, body.CategoryID)
	response.NewSuccessMessage(w, "Category linked")
}

func UnlinkTalkCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(r.PathValue("category_id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid category_id", http.StatusBadRequest)
		return
	}

	talk_models.UnlinkCategory(id, categoryID)
	response.NewSuccessMessage(w, "Category unlinked")
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

func GetTalkEventsHandler(w http.ResponseWriter, r *http.Request) {
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

	events := talk_models.GetTalkEvents(id)
	response.NewSuccessData(w, events, "")
}

func LinkTalkEventHandler(w http.ResponseWriter, r *http.Request) {
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
		EventID int `json:"event_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	talk_models.LinkEvent(id, body.EventID)
	response.NewSuccessMessage(w, "Event linked")
}

func UnlinkTalkEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	eventID, err := strconv.Atoi(r.PathValue("event_id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid event_id", http.StatusBadRequest)
		return
	}

	talk_models.UnlinkEvent(id, eventID)
	response.NewSuccessMessage(w, "Event unlinked")
}

func GetTalkProjectsHandler(w http.ResponseWriter, r *http.Request) {
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

	projects := talk_models.GetTalkProjects(id)
	response.NewSuccessData(w, projects, "")
}

func LinkTalkProjectHandler(w http.ResponseWriter, r *http.Request) {
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
		ProjectID int `json:"project_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	talk_models.LinkProject(id, body.ProjectID)
	response.NewSuccessMessage(w, "Project linked")
}

func UnlinkTalkProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	projectID, err := strconv.Atoi(r.PathValue("project_id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid project_id", http.StatusBadRequest)
		return
	}

	talk_models.UnlinkProject(id, projectID)
	response.NewSuccessMessage(w, "Project unlinked")
}
