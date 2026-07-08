package user_actions

import (
	"go-forum-backend/app/models/user_models"
	"go-forum-backend/utils/rules"
)

// CreateUserFromToken creates a local shadow user from the authenticated
// userId carried by the central auth token. Username and email default to the
// userId (a UUID) so the UNIQUE/NOT NULL constraints are satisfied without a
// password; the profile is enriched later through UpdateUser.
func CreateUserFromToken(userId string) ([]rules.ValidationError, *user_models.User) {
	var errs []rules.ValidationError

	rules.StringMinLength(userId, 1, "userId", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	var user user_models.User
	err := user.Create(user_models.CreateUserDTO{
		Id:        userId,
		Username:  userId,
		Firstname: "",
		Lastname:  "",
		Email:     userId,
	})
	if err != nil {
		return nil, nil
	}

	return nil, user_models.GetUserByID(userId)
}
