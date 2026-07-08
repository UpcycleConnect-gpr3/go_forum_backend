package event_actions

import (
	"go-forum-backend/app/models/event_models"
	"go-forum-backend/utils/rules"
	"time"
)

type CreateEventActionDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StartAt     string `json:"start_at"`
	EndAt       string `json:"end_at"`
}

func CreateEvent(dto CreateEventActionDTO) ([]rules.ValidationError, *event_models.Event) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Title, 1, "title", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	startAt, err := time.Parse(time.RFC3339, dto.StartAt)
	if err != nil {
		errs = append(errs, rules.ValidationError{Field: "start_at", Message: "start_at must be a valid RFC3339 date"})
		return errs, nil
	}

	endAt, err := time.Parse(time.RFC3339, dto.EndAt)
	if err != nil {
		errs = append(errs, rules.ValidationError{Field: "end_at", Message: "end_at must be a valid RFC3339 date"})
		return errs, nil
	}

	event := event_models.CreateEvent(event_models.CreateEventDTO{
		Title:       dto.Title,
		Description: dto.Description,
		StartAt:     startAt,
		EndAt:       endAt,
	})

	return nil, event
}
