package auth_handlers

import (
	"go-forum-backend/app/actions/user_actions"
	"go-forum-backend/app/models/user_models"
	"go-forum-backend/utils/jwt"
	"go-forum-backend/utils/log"
	"go-forum-backend/utils/response"
	"net/http"
)

// LoginHandler provisions the local forum user from a token issued by the
// central auth service. It never generates a token itself: it verifies the
// bearer token, then returns the matching user, creating it on first login.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		response.NewErrorMessage(w, response.ErrAuthTokenRequired, http.StatusUnauthorized)
		return
	}

	userId, err := jwt.VerifyJWT(tokenString)
	if err != nil {
		response.NewErrorMessage(w, response.ErrInvalidAuthToken, http.StatusUnauthorized)
		return
	}

	if existing := user_models.GetUserByID(userId); existing != nil {
		response.NewSuccessData(w, existing, "")
		return
	}

	validationErrors, newUser := user_actions.CreateUserFromToken(userId)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}
	if newUser == nil {
		response.NewErrorMessage(w, "Could not provision user", http.StatusInternalServerError)
		return
	}

	response.NewSuccessData(w, newUser, "")
}
