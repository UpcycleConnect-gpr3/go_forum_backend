package message_actions

import (
	"go-forum-backend/app/models/message_models"
	"go-forum-backend/utils/rules"
)

type CreateMessageDTO struct {
	Content  string `json:"content"`
	FilePath string `json:"file_path"`
}

func CreateMessage(dto CreateMessageDTO) ([]rules.ValidationError, *message_models.Message) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Content, 1, "content", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	message := message_models.CreateMessage(message_models.CreateMessageDTO{
		Content:  dto.Content,
		FilePath: dto.FilePath,
	})

	return nil, message
}
