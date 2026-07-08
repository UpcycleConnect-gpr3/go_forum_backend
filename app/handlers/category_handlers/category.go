package category_handlers

import (
	"encoding/json"
	"go-forum-backend/app/actions/category_actions"
	"go-forum-backend/app/models/category_models"
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

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	page, limit := parsePage(r)
	categories := category_models.GetAllCategories(page, limit)
	response.NewSuccessData(w, categories, "")
}

func GetCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	category := category_models.GetCategoryByID(id)
	if category == nil {
		response.NewErrorMessage(w, response.ErrCategoryNotFound, http.StatusNotFound)
		return
	}

	response.NewSuccessData(w, category, "")
}

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var dto category_actions.CreateCategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, category := category_actions.CreateCategory(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, category, "")
}

func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if category_models.GetCategoryByID(id) == nil {
		response.NewErrorMessage(w, response.ErrCategoryNotFound, http.StatusNotFound)
		return
	}

	var dto category_actions.UpdateCategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, category := category_actions.UpdateCategory(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, category, "")
}

func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if category_models.GetCategoryByID(id) == nil {
		response.NewErrorMessage(w, response.ErrCategoryNotFound, http.StatusNotFound)
		return
	}

	category_models.DeleteCategory(id)
	response.NewSuccessMessage(w, "Category deleted")
}

func GetCategoryTalksHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if category_models.GetCategoryByID(id) == nil {
		response.NewErrorMessage(w, response.ErrCategoryNotFound, http.StatusNotFound)
		return
	}

	talks := category_models.GetCategoryTalks(id)
	response.NewSuccessData(w, talks, "")
}
