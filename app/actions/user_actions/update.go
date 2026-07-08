package user_actions

import (
	"go-forum-backend/app/models/user_models"
	"go-forum-backend/utils/rules"
)

type UpdateUserActionDTO struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func UpdateUser(id string, dto UpdateUserActionDTO) ([]rules.ValidationError, *user_models.User) {
	var errs []rules.ValidationError

	if dto.Username != "" {
		rules.StringMinLength(dto.Username, 3, "username", &errs)
	}
	if len(errs) > 0 {
		return errs, nil
	}

	user := user_models.UpdateUser(id, user_models.UpdateUserDTO{
		Username:  dto.Username,
		Firstname: dto.Firstname,
		Lastname:  dto.Lastname,
	})

	return nil, user
}
