package category_actions

import (
	"go-forum-backend/app/models/category_models"
	"go-forum-backend/utils/rules"
)

type UpdateCategoryDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func UpdateCategory(id int, dto UpdateCategoryDTO) ([]rules.ValidationError, *category_models.Category) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	category := category_models.UpdateCategory(id, category_models.UpdateCategoryDTO{
		Name:        dto.Name,
		Description: dto.Description,
	})

	return nil, category
}
