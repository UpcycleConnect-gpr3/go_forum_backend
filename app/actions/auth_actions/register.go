package auth_actions

import (
	"go-forum-backend/app/models/user_models"
	"go-forum-backend/utils/rules"
)

type RegisterDTO struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func Register(dto RegisterDTO) ([]rules.ValidationError, *user_models.User) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Email, 5, "email", &errs)
	rules.StringMinLength(dto.Password, 6, "password", &errs)
	rules.StringMaxLength(dto.Password, 72, "password", &errs)
	rules.StringMinLength(dto.Username, 3, "username", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	if user_models.GetUserByEmail(dto.Email) != nil {
		errs = append(errs, rules.ValidationError{Field: "email", Message: "email already taken"})
		return errs, nil
	}

	user := user_models.CreateUser(user_models.Credentials{
		Email:    dto.Email,
		Password: dto.Password,
	})

	return nil, user
}
