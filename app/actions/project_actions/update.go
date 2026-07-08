package project_actions

import (
	"go-forum-backend/app/models/project_models"
	"go-forum-backend/utils/rules"
)

type UpdateProjectDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func UpdateProject(id int, dto UpdateProjectDTO) ([]rules.ValidationError, *project_models.Project) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Title, 1, "title", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	project := project_models.UpdateProject(id, project_models.UpdateProjectDTO{
		Title:       dto.Title,
		Description: dto.Description,
	})

	return nil, project
}
