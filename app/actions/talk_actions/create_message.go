package talk_actions

import (
	"go-forum-backend/app/models/talk_models"
	"go-forum-backend/utils/rules"
)

type MessageDTO struct {
	Content  string `json:"content"`
	FilePath string `json:"file_path"`
}

func validateMessage(dto MessageDTO) []rules.ValidationError {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Content, 1, "content", &errs)

	return errs
}

func CreateTalkMessage(talkID int, dto MessageDTO) ([]rules.ValidationError, *talk_models.MessageSummary) {
	errs := validateMessage(dto)
	if len(errs) > 0 {
		return errs, nil
	}

	msg := talk_models.CreateTalkMessage(talkID, dto.Content, dto.FilePath)
	return nil, msg
}
