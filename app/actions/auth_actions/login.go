package auth_actions

import (
	"go-forum-backend/app/models/user_models"
	"go-forum-backend/utils/response"
	"go-forum-backend/utils/rules"
)

func Login(credentials user_models.Credentials) ([]rules.ValidationError, *user_models.User) {
	var errs []rules.ValidationError

	rules.StringMinLength(credentials.Email, 5, "email", &errs)
	rules.StringMinLength(credentials.Password, 6, "password", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	existing := user_models.GetUserByEmail(credentials.Email)
	if existing == nil || !existing.CheckPassword(credentials.Password) {
		errs = append(errs, rules.ValidationError{Field: "email", Message: response.ErrAuthFailed})
		return errs, nil
	}

	return nil, existing
}
