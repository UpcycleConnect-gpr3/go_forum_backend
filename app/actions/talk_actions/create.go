package talk_actions

import (
	"go-forum-backend/app/models/talk_models"
	"go-forum-backend/utils/rules"
)

type CreateTalkDTO struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	CategoryID int    `json:"category_id"`
}

func CreateTalk(dto CreateTalkDTO) ([]rules.ValidationError, *talk_models.Talk) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.StringMinLength(dto.Content, 1, "content", &errs)
	rules.IntMinLength(dto.CategoryID, 1, "category_id", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	talk := talk_models.CreateTalk(talk_models.CreateTalkDTO{
		Title:      dto.Title,
		Content:    dto.Content,
		CategoryID: dto.CategoryID,
	})

	return nil, talk
}
