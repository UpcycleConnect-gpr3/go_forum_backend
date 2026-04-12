package metric_handlers

import (
	"go-forum-backend/utils/log"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
}
