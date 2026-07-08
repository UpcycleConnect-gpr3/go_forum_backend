package message_actions

import (
	"go-forum-backend/app/models/message_models"
	"go-forum-backend/utils/rules"
)

type UpdateMessageDTO struct {
	Content  string `json:"content"`
	FilePath string `json:"file_path"`
}

func UpdateMessage(id int, dto UpdateMessageDTO) ([]rules.ValidationError, *message_models.Message) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Content, 1, "content", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	message := message_models.UpdateMessage(id, message_models.UpdateMessageDTO{
		Content:  dto.Content,
		FilePath: dto.FilePath,
	})

	return nil, message
}
