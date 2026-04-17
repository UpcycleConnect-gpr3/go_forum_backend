package server

import (
	"go-forum-backend/app/handlers/auth_handlers"
	"go-forum-backend/app/handlers/category_handlers"
	"go-forum-backend/app/handlers/event_handlers"
	"go-forum-backend/app/handlers/message_handlers"
	"go-forum-backend/app/handlers/metric_handlers"
	"go-forum-backend/app/handlers/project_handlers"
	"go-forum-backend/app/handlers/talk_handlers"
	"go-forum-backend/app/handlers/user_handlers"
	"go-forum-backend/app/middleware/auth_middleware"
	"go-forum-backend/app/middleware/ratelimit_middleware"
	"go-forum-backend/app/middleware/source_middleware"
	"go-forum-backend/config"
	"go-forum-backend/database"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	log "github.com/thedataflows/go-lib-log"
)

func initialize() {
	logger := log.NewLoggerBuilder().WithLogLevel(zerolog.DebugLevel).WithBufferSize(10000).WithRateLimit(1000).WithGroupWindow(2 * time.Second).WithLogFormat(log.LOG_FORMAT_JSON).WithOutput(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Build()
	defer logger.Close()

	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal().Err(err).Msg("Error loading .env file")
	}

	// Config Initialization
	config.InitDatabase()

	err = database.Forum.Ping()
	if err != nil {
		logger.Fatal().Err(err).Msg("(DATABASE)")
	}
}

func Start() {
	initialize()

	limiterLow := ratelimit_middleware.NewRateLimiter(10, 1*time.Minute)
	limiterMedium := ratelimit_middleware.NewRateLimiter(30, 1*time.Minute)
	limiterHigh := ratelimit_middleware.NewRateLimiter(60, 1*time.Minute)

	containerApp := source_middleware.Container("app")

	logger := log.NewLoggerBuilder().WithLogLevel(zerolog.DebugLevel).WithBufferSize(10000).WithRateLimit(1000).WithGroupWindow(2 * time.Second).WithOutput(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Build()
	defer logger.Close()

	// ── Route 1: Health ───────────────────────────────────────────────────────
	http.HandleFunc("GET /health/{$}", limiterLow.RateLimit(containerApp(metric_handlers.Health)))

	// ── Routes 2-3: Auth ──────────────────────────────────────────────────────
	http.HandleFunc("POST /auth/login/{$}", limiterLow.RateLimit(containerApp(auth_handlers.LoginHandler)))
	http.HandleFunc("POST /auth/register/{$}", limiterLow.RateLimit(containerApp(auth_handlers.RegisterHandler)))

	// ── Routes 4-8: Users ─────────────────────────────────────────────────────
	http.HandleFunc("GET /users/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(user_handlers.GetUsersHandler))))
	http.HandleFunc("GET /users/{id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(user_handlers.GetUserHandler))))
	http.HandleFunc("PUT /users/{id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(user_handlers.UpdateUserHandler))))
	http.HandleFunc("DELETE /users/{id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(user_handlers.DeleteUserHandler))))
	http.HandleFunc("GET /users/{id}/messages/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(user_handlers.GetUserMessagesHandler))))

	// ── Routes 9-14: Categories ───────────────────────────────────────────────
	http.HandleFunc("GET /categories/{$}", limiterHigh.RateLimit(containerApp(category_handlers.GetCategoriesHandler)))
	http.HandleFunc("GET /categories/{id}/{$}", limiterHigh.RateLimit(containerApp(category_handlers.GetCategoryHandler)))
	http.HandleFunc("POST /categories/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(category_handlers.CreateCategoryHandler))))
	http.HandleFunc("PUT /categories/{id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(category_handlers.UpdateCategoryHandler))))
	http.HandleFunc("DELETE /categories/{id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(category_handlers.DeleteCategoryHandler))))
	http.HandleFunc("GET /categories/{id}/talks/{$}", limiterHigh.RateLimit(containerApp(category_handlers.GetCategoryTalksHandler)))

	// ── Routes 15-25: Talks ───────────────────────────────────────────────────
	http.HandleFunc("GET /talks/{$}", limiterHigh.RateLimit(containerApp(talk_handlers.GetTalksHandler)))
	http.HandleFunc("GET /talks/{id}/{$}", limiterHigh.RateLimit(containerApp(talk_handlers.GetTalkHandler)))
	http.HandleFunc("POST /talks/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(talk_handlers.CreateTalkHandler))))
	http.HandleFunc("PUT /talks/{id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(talk_handlers.UpdateTalkHandler))))
	http.HandleFunc("DELETE /talks/{id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(talk_handlers.DeleteTalkHandler))))
	http.HandleFunc("GET /talks/{id}/messages/{$}", limiterHigh.RateLimit(containerApp(talk_handlers.GetTalkMessagesHandler)))
	http.HandleFunc("POST /talks/{id}/messages/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(talk_handlers.LinkTalkMessageHandler))))
	http.HandleFunc("DELETE /talks/{id}/messages/{message_id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(talk_handlers.UnlinkTalkMessageHandler))))
	http.HandleFunc("GET /talks/{id}/users/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(talk_handlers.GetTalkUsersHandler))))
	http.HandleFunc("POST /talks/{id}/users/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(talk_handlers.LinkTalkUserHandler))))
	http.HandleFunc("DELETE /talks/{id}/users/{user_id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(talk_handlers.UnlinkTalkUserHandler))))

	// ── Routes 26-33: Messages ────────────────────────────────────────────────
	http.HandleFunc("GET /messages/{$}", limiterHigh.RateLimit(containerApp(message_handlers.GetMessagesHandler)))
	http.HandleFunc("GET /messages/{id}/{$}", limiterHigh.RateLimit(containerApp(message_handlers.GetMessageHandler)))
	http.HandleFunc("POST /messages/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(message_handlers.CreateMessageHandler))))
	http.HandleFunc("PUT /messages/{id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(message_handlers.UpdateMessageHandler))))
	http.HandleFunc("DELETE /messages/{id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(message_handlers.DeleteMessageHandler))))
	http.HandleFunc("GET /messages/{id}/users/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(message_handlers.GetMessageUsersHandler))))
	http.HandleFunc("POST /messages/{id}/users/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(message_handlers.LinkMessageUserHandler))))
	http.HandleFunc("DELETE /messages/{id}/users/{user_id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(message_handlers.UnlinkMessageUserHandler))))

	// ── Routes 34-41: Events ──────────────────────────────────────────────────
	http.HandleFunc("GET /events/{$}", limiterHigh.RateLimit(containerApp(event_handlers.GetEventsHandler)))
	http.HandleFunc("GET /events/{id}/{$}", limiterHigh.RateLimit(containerApp(event_handlers.GetEventHandler)))
	http.HandleFunc("POST /events/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(event_handlers.CreateEventHandler))))
	http.HandleFunc("PUT /events/{id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(event_handlers.UpdateEventHandler))))
	http.HandleFunc("DELETE /events/{id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(event_handlers.DeleteEventHandler))))
	http.HandleFunc("GET /events/{id}/users/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(event_handlers.GetEventUsersHandler))))
	http.HandleFunc("POST /events/{id}/users/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(event_handlers.LinkEventUserHandler))))
	http.HandleFunc("DELETE /events/{id}/users/{user_id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(event_handlers.UnlinkEventUserHandler))))

	// ── Routes 42-46: Projects ────────────────────────────────────────────────
	http.HandleFunc("GET /projects/{$}", limiterHigh.RateLimit(containerApp(project_handlers.GetProjectsHandler)))
	http.HandleFunc("GET /projects/{id}/{$}", limiterHigh.RateLimit(containerApp(project_handlers.GetProjectHandler)))
	http.HandleFunc("POST /projects/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(project_handlers.CreateProjectHandler))))
	http.HandleFunc("PUT /projects/{id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(project_handlers.UpdateProjectHandler))))
	http.HandleFunc("DELETE /projects/{id}/{$}", limiterMedium.RateLimit(containerApp(auth_middleware.IsAuth(project_handlers.DeleteProjectHandler))))

	logger.Info().Msg("Listening at http://localhost:" + os.Getenv("APP_PORT"))
	err := http.ListenAndServe(":"+os.Getenv("APP_PORT"), nil)
	if err != nil {
		return
	}
}
