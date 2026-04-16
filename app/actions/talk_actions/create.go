package talk_actions

import (
	"go-forum-backend/app/models/talk_models"
	"go-forum-backend/utils/rules"
)

type TalkDTO struct {
	Title       string `json:"title"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

func validateCreate(dto TalkDTO) []rules.ValidationError {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.StringMaxLength(dto.Title, 255, "title", &errs)
	rules.StringMinLength(dto.Type, 1, "type", &errs)
	rules.StringMinLength(dto.Status, 1, "status", &errs)

	return errs
}

func CreateTalk(dto TalkDTO) ([]rules.ValidationError, *talk_models.Talk) {
	errs := validateCreate(dto)
	if len(errs) > 0 {
		return errs, nil
	}

	talk := talk_models.CreateTalk(talk_models.CreateTalkDTO{
		Title:       dto.Title,
		Type:        dto.Type,
		Status:      dto.Status,
		Description: dto.Description,
	})

	return nil, talk
}
