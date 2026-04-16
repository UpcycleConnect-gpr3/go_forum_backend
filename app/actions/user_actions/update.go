package user_actions

import (
	"go-forum-backend/app/models/user_models"
	"go-forum-backend/utils/rules"
)

type UpdateUserDTO struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func validateUpdate(dto UpdateUserDTO) []rules.ValidationError {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Username, 3, "username", &errs)
	rules.StringMaxLength(dto.Username, 50, "username", &errs)
	rules.StringMinLength(dto.Firstname, 1, "firstname", &errs)
	rules.StringMinLength(dto.Lastname, 1, "lastname", &errs)

	return errs
}

func UpdateUser(id string, dto UpdateUserDTO) ([]rules.ValidationError, *user_models.User) {
	errs := validateUpdate(dto)
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
