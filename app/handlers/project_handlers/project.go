package project_handlers

import (
	"encoding/json"
	"go-forum-backend/app/actions/project_actions"
	"go-forum-backend/app/models/project_models"
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

func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	page, limit := parsePage(r)
	projects := project_models.GetAllProjects(page, limit)
	response.NewSuccessData(w, projects, "")
}

func GetProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	project := project_models.GetProjectByID(id)
	if project == nil {
		response.NewErrorMessage(w, response.ErrProjectNotFound, http.StatusNotFound)
		return
	}

	response.NewSuccessData(w, project, "")
}

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var dto project_actions.CreateProjectDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, project := project_actions.CreateProject(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, project, "")
}

func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if project_models.GetProjectByID(id) == nil {
		response.NewErrorMessage(w, response.ErrProjectNotFound, http.StatusNotFound)
		return
	}

	var dto project_actions.UpdateProjectDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, project := project_actions.UpdateProject(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, project, "")
}

func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if project_models.GetProjectByID(id) == nil {
		response.NewErrorMessage(w, response.ErrProjectNotFound, http.StatusNotFound)
		return
	}

	project_models.DeleteProject(id)
	response.NewSuccessMessage(w, "Project deleted")
}
