package event_models

import (
	"database/sql"
	"fmt"
	"go-forum-backend/database"
	"go-forum-backend/utils/log"
)

const TABLE = "EVENTS"

type Event struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"`
}

func GetEventByID(id int) *Event {
	event := Event{}
	action := fmt.Sprintf("SELECT "+TABLE+" WHERE id : %d", id)

	row := database.Forum.QueryRow("SELECT id, title, date FROM "+TABLE+" WHERE id = ?", id)

	err := row.Scan(&event.Id, &event.Title, &event.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Database(action, err)
		return nil
	}

	return &event
}
