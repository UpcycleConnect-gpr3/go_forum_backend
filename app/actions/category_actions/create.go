package category_actions

import (
	"go-forum-backend/app/models/category_models"
	"go-forum-backend/utils/rules"
)

type CategoryDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func validateCreate(dto CategoryDTO) []rules.ValidationError {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)

	return errs
}

func CreateCategory(dto CategoryDTO) ([]rules.ValidationError, *category_models.Category) {
	errs := validateCreate(dto)
	if len(errs) > 0 {
		return errs, nil
	}

	category := category_models.CreateCategory(category_models.CategoryDTO{
		Name:        dto.Name,
		Description: dto.Description,
	})

	return nil, category
}
