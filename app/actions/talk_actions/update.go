package talk_actions

import (
	"go-forum-backend/app/models/talk_models"
	"go-forum-backend/utils/rules"
)

type UpdateTalkDTO struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}

func validateUpdate(dto UpdateTalkDTO) []rules.ValidationError {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.StringMaxLength(dto.Title, 255, "title", &errs)
	rules.StringMinLength(dto.Status, 1, "status", &errs)

	return errs
}

func UpdateTalk(id int, dto UpdateTalkDTO) ([]rules.ValidationError, *talk_models.Talk) {
	errs := validateUpdate(dto)
	if len(errs) > 0 {
		return errs, nil
	}

	talk := talk_models.UpdateTalk(id, talk_models.UpdateTalkDTO{
		Title:  dto.Title,
		Status: dto.Status,
	})

	return nil, talk
}
