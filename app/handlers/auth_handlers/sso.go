package auth_handlers

import (
	"encoding/json"
	"go-forum-backend/app/middleware/auth_middleware"
	"go-forum-backend/app/models/user_models"
	"go-forum-backend/utils/log"
	"go-forum-backend/utils/response"
	"net/http"
	"strings"
)

type SSOLoginDTO struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

func SSOLoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	userId := auth_middleware.GetUserId(r.Context())

	var dto SSOLoginDTO
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&dto)
	}

	username := dto.Username
	if username == "" && dto.Email != "" {
		username = strings.SplitN(dto.Email, "@", 2)[0]
	}
	if username == "" {
		username = userId
	}
	email := dto.Email
	if email == "" {
		email = userId
	}

	user_models.EnsureUser(userId, username, dto.Firstname, dto.Lastname, email)

	response.NewSuccessData(w, map[string]string{
		"id":        userId,
		"username":  username,
		"firstname": dto.Firstname,
		"lastname":  dto.Lastname,
		"email":     email,
	}, "")
}
