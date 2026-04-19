package project_actions

import (
	"go-forum-backend/app/models/project_models"
	"go-forum-backend/utils/rules"
)

type CreateProjectDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateProject(dto CreateProjectDTO) ([]rules.ValidationError, *project_models.Project) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Title, 1, "title", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	project := project_models.CreateProject(project_models.CreateProjectDTO{
		Title:       dto.Title,
		Description: dto.Description,
	})

	return nil, project
}
