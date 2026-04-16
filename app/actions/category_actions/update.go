package category_actions

import (
	"go-forum-backend/app/models/category_models"
	"go-forum-backend/utils/rules"
)

func validateUpdate(dto CategoryDTO) []rules.ValidationError {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMaxLength(dto.Name, 255, "name", &errs)

	return errs
}

func UpdateCategory(id int, dto CategoryDTO) ([]rules.ValidationError, *category_models.Category) {
	errs := validateUpdate(dto)
	if len(errs) > 0 {
		return errs, nil
	}

	category := category_models.UpdateCategory(id, category_models.CategoryDTO{
		Name:        dto.Name,
		Description: dto.Description,
	})

	return nil, category
}
