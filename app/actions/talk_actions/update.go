package talk_actions

import (
	"go-forum-backend/app/models/talk_models"
	"go-forum-backend/utils/rules"
)

type UpdateTalkDTO struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func UpdateTalk(id int, dto UpdateTalkDTO) ([]rules.ValidationError, *talk_models.Talk) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.StringMinLength(dto.Content, 1, "content", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	talk := talk_models.UpdateTalk(id, talk_models.UpdateTalkDTO{
		Title:   dto.Title,
		Content: dto.Content,
	})

	return nil, talk
}
