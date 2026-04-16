package talk_models

import (
	"database/sql"
	"fmt"
	"go-forum-backend/database"
	"go-forum-backend/utils/log"
	"time"
)

const TABLE = "TALKS"

type Talk struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type CreateTalkDTO struct {
	Title       string
	Type        string
	Status      string
	Description string
}

type UpdateTalkDTO struct {
	Title  string
	Status string
}

type MessageSummary struct {
	Id        int            `json:"id"`
	Content   string         `json:"content"`
	FilePath  sql.NullString `json:"file_path"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type CategorySummary struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UserSummary struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type EventSummary struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"`
}

type ProjectSummary struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllTalks(page, limit int) []Talk {
	action := "SELECT " + TABLE + " (paginated)"
	offset := (page - 1) * limit

	rows, err := database.Forum.Query(
		"SELECT id, title, type, status, description, created_at, updated_at FROM "+TABLE+" LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		log.Database(action, err)
		return []Talk{}
	}
	defer rows.Close()

	talks := []Talk{}
	for rows.Next() {
		var t Talk
		if err := rows.Scan(&t.Id, &t.Title, &t.Type, &t.Status, &t.Description, &t.CreatedAt, &t.UpdatedAt); err != nil {
			log.Database(action, err)
			continue
		}
		talks = append(talks, t)
	}
	return talks
}

func GetTalkByID(id int) *Talk {
	talk := Talk{}
	action := fmt.Sprintf("SELECT "+TABLE+" WHERE id : %d", id)

	row := database.Forum.QueryRow(
		"SELECT id, title, type, status, description, created_at, updated_at FROM "+TABLE+" WHERE id = ?",
		id,
	)

	err := row.Scan(&talk.Id, &talk.Title, &talk.Type, &talk.Status, &talk.Description, &talk.CreatedAt, &talk.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Database(action, err)
		return nil
	}

	return &talk
}

func CreateTalk(dto CreateTalkDTO) *Talk {
	action := "INSERT INTO " + TABLE

	result, err := database.Forum.Exec(
		"INSERT INTO "+TABLE+" (title, type, status, description) VALUES (?, ?, ?, ?)",
		dto.Title, dto.Type, dto.Status, dto.Description,
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

	return GetTalkByID(int(id))
}

func UpdateTalk(id int, dto UpdateTalkDTO) *Talk {
	action := fmt.Sprintf("UPDATE "+TABLE+" WHERE id : %d", id)

	_, err := database.Forum.Exec(
		"UPDATE "+TABLE+" SET title = ?, status = ? WHERE id = ?",
		dto.Title, dto.Status, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}

	return GetTalkByID(id)
}

func DeleteTalk(id int) {
	action := fmt.Sprintf("DELETE FROM "+TABLE+" WHERE id : %d", id)

	_, err := database.Forum.Exec("DELETE FROM "+TABLE+" WHERE id = ?", id)
	if err != nil {
		log.Database(action, err)
	}
}

func GetTalkMessages(talkID, page, limit int) []MessageSummary {
	action := fmt.Sprintf("SELECT MESSAGES JOIN MESSAGE_TALK WHERE talk_id : %d", talkID)
	offset := (page - 1) * limit

	rows, err := database.Forum.Query(
		"SELECT m.id, m.content, m.file_path, m.created_at, m.updated_at FROM MESSAGES m JOIN MESSAGE_TALK mt ON m.id = mt.message_id WHERE mt.talk_id = ? LIMIT ? OFFSET ?",
		talkID, limit, offset,
	)
	if err != nil {
		log.Database(action, err)
		return []MessageSummary{}
	}
	defer rows.Close()

	messages := []MessageSummary{}
	for rows.Next() {
		var m MessageSummary
		if err := rows.Scan(&m.Id, &m.Content, &m.FilePath, &m.CreatedAt, &m.UpdatedAt); err != nil {
			log.Database(action, err)
			continue
		}
		messages = append(messages, m)
	}
	return messages
}

func CreateTalkMessage(talkID int, content, filePath string) *MessageSummary {
	action := fmt.Sprintf("INSERT MESSAGES + MESSAGE_TALK for talk_id : %d", talkID)

	result, err := database.Forum.Exec(
		"INSERT INTO MESSAGES (content, file_path) VALUES (?, ?)",
		content, sql.NullString{String: filePath, Valid: filePath != ""},
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}

	messageID, err := result.LastInsertId()
	if err != nil {
		log.Database(action, err)
		return nil
	}

	_, err = database.Forum.Exec(
		"INSERT IGNORE INTO MESSAGE_TALK (message_id, talk_id) VALUES (?, ?)",
		messageID, talkID,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}

	msg := &MessageSummary{}
	row := database.Forum.QueryRow(
		"SELECT id, content, file_path, created_at, updated_at FROM MESSAGES WHERE id = ?",
		messageID,
	)
	if err := row.Scan(&msg.Id, &msg.Content, &msg.FilePath, &msg.CreatedAt, &msg.UpdatedAt); err != nil {
		log.Database(action, err)
		return nil
	}

	return msg
}

func GetTalkCategories(talkID int) []CategorySummary {
	action := fmt.Sprintf("SELECT CATEGORIES JOIN CATEGORY_TALK WHERE talk_id : %d", talkID)

	rows, err := database.Forum.Query(
		"SELECT c.id, c.name, c.description FROM CATEGORIES c JOIN CATEGORY_TALK ct ON c.id = ct.category_id WHERE ct.talk_id = ?",
		talkID,
	)
	if err != nil {
		log.Database(action, err)
		return []CategorySummary{}
	}
	defer rows.Close()

	categories := []CategorySummary{}
	for rows.Next() {
		var c CategorySummary
		if err := rows.Scan(&c.Id, &c.Name, &c.Description); err != nil {
			log.Database(action, err)
			continue
		}
		categories = append(categories, c)
	}
	return categories
}

func LinkCategory(talkID, categoryID int) {
	action := fmt.Sprintf("INSERT INTO CATEGORY_TALK talk_id:%d category_id:%d", talkID, categoryID)

	_, err := database.Forum.Exec(
		"INSERT IGNORE INTO CATEGORY_TALK (category_id, talk_id) VALUES (?, ?)",
		categoryID, talkID,
	)
	if err != nil {
		log.Database(action, err)
	}
}

func UnlinkCategory(talkID, categoryID int) {
	action := fmt.Sprintf("DELETE FROM CATEGORY_TALK talk_id:%d category_id:%d", talkID, categoryID)

	_, err := database.Forum.Exec(
		"DELETE FROM CATEGORY_TALK WHERE category_id = ? AND talk_id = ?",
		categoryID, talkID,
	)
	if err != nil {
		log.Database(action, err)
	}
}

func GetTalkUsers(talkID int) []UserSummary {
	action := fmt.Sprintf("SELECT USERS JOIN USER_TALK WHERE talk_id : %d", talkID)

	rows, err := database.Forum.Query(
		"SELECT u.id, u.username, u.firstname, u.lastname FROM USERS u JOIN USER_TALK ut ON u.id = ut.user_id WHERE ut.talk_id = ?",
		talkID,
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

func LinkUser(talkID int, userID string) {
	action := fmt.Sprintf("INSERT INTO USER_TALK talk_id:%d user_id:%s", talkID, userID)

	_, err := database.Forum.Exec(
		"INSERT IGNORE INTO USER_TALK (user_id, talk_id) VALUES (?, ?)",
		userID, talkID,
	)
	if err != nil {
		log.Database(action, err)
	}
}

func UnlinkUser(talkID int, userID string) {
	action := fmt.Sprintf("DELETE FROM USER_TALK talk_id:%d user_id:%s", talkID, userID)

	_, err := database.Forum.Exec(
		"DELETE FROM USER_TALK WHERE user_id = ? AND talk_id = ?",
		userID, talkID,
	)
	if err != nil {
		log.Database(action, err)
	}
}

func GetTalkEvents(talkID int) []EventSummary {
	action := fmt.Sprintf("SELECT EVENTS JOIN TALK_EVENT WHERE talk_id : %d", talkID)

	rows, err := database.Forum.Query(
		"SELECT e.id, e.title, e.date FROM EVENTS e JOIN TALK_EVENT te ON e.id = te.event_id WHERE te.talk_id = ?",
		talkID,
	)
	if err != nil {
		log.Database(action, err)
		return []EventSummary{}
	}
	defer rows.Close()

	events := []EventSummary{}
	for rows.Next() {
		var e EventSummary
		if err := rows.Scan(&e.Id, &e.Title, &e.Date); err != nil {
			log.Database(action, err)
			continue
		}
		events = append(events, e)
	}
	return events
}

func LinkEvent(talkID, eventID int) {
	action := fmt.Sprintf("INSERT INTO TALK_EVENT talk_id:%d event_id:%d", talkID, eventID)

	_, err := database.Forum.Exec(
		"INSERT IGNORE INTO TALK_EVENT (talk_id, event_id) VALUES (?, ?)",
		talkID, eventID,
	)
	if err != nil {
		log.Database(action, err)
	}
}

func UnlinkEvent(talkID, eventID int) {
	action := fmt.Sprintf("DELETE FROM TALK_EVENT talk_id:%d event_id:%d", talkID, eventID)

	_, err := database.Forum.Exec(
		"DELETE FROM TALK_EVENT WHERE talk_id = ? AND event_id = ?",
		talkID, eventID,
	)
	if err != nil {
		log.Database(action, err)
	}
}

func GetTalkProjects(talkID int) []ProjectSummary {
	action := fmt.Sprintf("SELECT PROJECTS JOIN TALK_PROJECT WHERE talk_id : %d", talkID)

	rows, err := database.Forum.Query(
		"SELECT p.id, p.name FROM PROJECTS p JOIN TALK_PROJECT tp ON p.id = tp.project_id WHERE tp.talk_id = ?",
		talkID,
	)
	if err != nil {
		log.Database(action, err)
		return []ProjectSummary{}
	}
	defer rows.Close()

	projects := []ProjectSummary{}
	for rows.Next() {
		var p ProjectSummary
		if err := rows.Scan(&p.Id, &p.Name); err != nil {
			log.Database(action, err)
			continue
		}
		projects = append(projects, p)
	}
	return projects
}

func LinkProject(talkID, projectID int) {
	action := fmt.Sprintf("INSERT INTO TALK_PROJECT talk_id:%d project_id:%d", talkID, projectID)

	_, err := database.Forum.Exec(
		"INSERT IGNORE INTO TALK_PROJECT (talk_id, project_id) VALUES (?, ?)",
		talkID, projectID,
	)
	if err != nil {
		log.Database(action, err)
	}
}

func UnlinkProject(talkID, projectID int) {
	action := fmt.Sprintf("DELETE FROM TALK_PROJECT talk_id:%d project_id:%d", talkID, projectID)

	_, err := database.Forum.Exec(
		"DELETE FROM TALK_PROJECT WHERE talk_id = ? AND project_id = ?",
		talkID, projectID,
	)
	if err != nil {
		log.Database(action, err)
	}
}
