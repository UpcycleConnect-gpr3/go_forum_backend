package category_actions

import (
	"go-forum-backend/app/models/category_models"
	"go-forum-backend/utils/rules"
)

type CreateCategoryDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateCategory(dto CreateCategoryDTO) ([]rules.ValidationError, *category_models.Category) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	category := category_models.CreateCategory(category_models.CreateCategoryDTO{
		Name:        dto.Name,
		Description: dto.Description,
	})

	return nil, category
}
