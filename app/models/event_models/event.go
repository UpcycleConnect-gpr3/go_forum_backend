package event_models

import (
	"database/sql"
	"fmt"
	"go-forum-backend/database"
	"go-forum-backend/utils/log"
	"time"
)

const TABLE = "EVENTS"

type Event struct {
	Id          int            `json:"id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	StartAt     time.Time      `db:"start_at" json:"start_at"`
	EndAt       time.Time      `db:"end_at" json:"end_at"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
}

type CreateEventDTO struct {
	Title       string
	Description string
	StartAt     time.Time
	EndAt       time.Time
}

type UpdateEventDTO struct {
	Title       string
	Description string
	StartAt     time.Time
	EndAt       time.Time
}

type UserSummary struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func GetAllEvents(page, limit int) []Event {
	action := "SELECT " + TABLE + " (paginated)"
	offset := (page - 1) * limit

	rows, err := database.Forum.Query(
		"SELECT id, title, description, start_at, end_at, created_at, updated_at FROM "+TABLE+" LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		log.Database(action, err)
		return []Event{}
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.Id, &e.Title, &e.Description, &e.StartAt, &e.EndAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			log.Database(action, err)
			continue
		}
		events = append(events, e)
	}
	return events
}

func GetEventByID(id int) *Event {
	event := Event{}
	action := fmt.Sprintf("SELECT "+TABLE+" WHERE id : %d", id)

	row := database.Forum.QueryRow(
		"SELECT id, title, description, start_at, end_at, created_at, updated_at FROM "+TABLE+" WHERE id = ?",
		id,
	)

	err := row.Scan(&event.Id, &event.Title, &event.Description, &event.StartAt, &event.EndAt, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Database(action, err)
		return nil
	}

	return &event
}

func CreateEvent(dto CreateEventDTO) *Event {
	action := "INSERT INTO " + TABLE

	result, err := database.Forum.Exec(
		"INSERT INTO "+TABLE+" (title, description, start_at, end_at) VALUES (?, ?, ?, ?)",
		dto.Title,
		sql.NullString{String: dto.Description, Valid: dto.Description != ""},
		dto.StartAt,
		dto.EndAt,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Database(action, err)
		return nil
	}

	return GetEventByID(int(id))
}

func UpdateEvent(id int, dto UpdateEventDTO) *Event {
	action := fmt.Sprintf("UPDATE "+TABLE+" WHERE id : %d", id)

	_, err := database.Forum.Exec(
		"UPDATE "+TABLE+" SET title = ?, description = ?, start_at = ?, end_at = ? WHERE id = ?",
		dto.Title,
		sql.NullString{String: dto.Description, Valid: dto.Description != ""},
		dto.StartAt,
		dto.EndAt,
		id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}

	return GetEventByID(id)
}

func DeleteEvent(id int) {
	action := fmt.Sprintf("DELETE FROM "+TABLE+" WHERE id : %d", id)

	_, err := database.Forum.Exec("DELETE FROM "+TABLE+" WHERE id = ?", id)
	if err != nil {
		log.Database(action, err)
	}
}

func GetEventUsers(eventID int) []UserSummary {
	action := fmt.Sprintf("SELECT USERS JOIN USER_EVENT WHERE event_id : %d", eventID)

	rows, err := database.Forum.Query(
		"SELECT u.id, u.username, u.firstname, u.lastname FROM USERS u JOIN USER_EVENT ue ON u.id = ue.user_id WHERE ue.event_id = ?",
		eventID,
	)
	if err != nil {
		log.Database(action, err)
		return []UserSummary{}
	}
	defer rows.Close()

	users := []UserSummary{}
	for rows.Next() {
		var u UserSummary
		if err := rows.Scan(&u.Id, &u.Username, &u.Firstname, &u.Lastname); err != nil {
			log.Database(action, err)
			continue
		}
		users = append(users, u)
	}
	return users
}

func LinkUser(eventID int, userID string) {
	action := fmt.Sprintf("INSERT INTO USER_EVENT event_id:%d user_id:%s", eventID, userID)

	_, err := database.Forum.Exec(
		"INSERT IGNORE INTO USER_EVENT (user_id, event_id) VALUES (?, ?)",
		userID, eventID,
	)
	if err != nil {
		log.Database(action, err)
	}
}

func UnlinkUser(eventID int, userID string) {
	action := fmt.Sprintf("DELETE FROM USER_EVENT event_id:%d user_id:%s", eventID, userID)

	_, err := database.Forum.Exec(
		"DELETE FROM USER_EVENT WHERE user_id = ? AND event_id = ?",
		userID, eventID,
	)
	if err != nil {
		log.Database(action, err)
	}
}
