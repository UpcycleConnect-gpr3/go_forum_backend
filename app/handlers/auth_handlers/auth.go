package auth_handlers

import (
	"encoding/json"
	"go-forum-backend/app/actions/auth_actions"
	"go-forum-backend/app/models/user_models"
	"go-forum-backend/utils/jwt"
	"go-forum-backend/utils/log"
	"go-forum-backend/utils/response"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var credentials user_models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, user := auth_actions.Login(credentials)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	token, err := jwt.GenerateJWT(user.Id.String())
	if err != nil {
		response.NewErrorMessage(w, response.ErrGenerateToken, http.StatusInternalServerError)
		return
	}

	response.NewSuccessData(w, map[string]string{"bearer_token": token}, "")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var dto auth_actions.RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, user := auth_actions.Register(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    user,
	})
}
