package message_models

import (
	"database/sql"
	"fmt"
	"go-forum-backend/database"
	"go-forum-backend/utils/log"
	"time"
)

const TABLE = "MESSAGES"

type Message struct {
	Id        int            `json:"id"`
	Content   string         `json:"content"`
	FilePath  sql.NullString `db:"file_path" json:"file_path"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}

type CreateMessageDTO struct {
	Content  string
	FilePath string
}

type UpdateMessageDTO struct {
	Content string
}

type UserSummary struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func GetAllMessages(page, limit int) []Message {
	action := "SELECT " + TABLE + " (paginated)"
	offset := (page - 1) * limit

	rows, err := database.Forum.Query(
		"SELECT id, content, file_path, created_at, updated_at FROM "+TABLE+" LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		log.Database(action, err)
		return []Message{}
	}
	defer rows.Close()

	messages := []Message{}
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.Id, &m.Content, &m.FilePath, &m.CreatedAt, &m.UpdatedAt); err != nil {
			log.Database(action, err)
			continue
		}
		messages = append(messages, m)
	}
	return messages
}

func GetMessageByID(id int) *Message {
	message := Message{}
	action := fmt.Sprintf("SELECT "+TABLE+" WHERE id : %d", id)

	row := database.Forum.QueryRow(
		"SELECT id, content, file_path, created_at, updated_at FROM "+TABLE+" WHERE id = ?",
		id,
	)

	err := row.Scan(&message.Id, &message.Content, &message.FilePath, &message.CreatedAt, &message.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Database(action, err)
		return nil
	}

	return &message
}

func CreateMessage(dto CreateMessageDTO) *Message {
	action := "INSERT INTO " + TABLE

	result, err := database.Forum.Exec(
		"INSERT INTO "+TABLE+" (content, file_path) VALUES (?, ?)",
		dto.Content, sql.NullString{String: dto.FilePath, Valid: dto.FilePath != ""},
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

	return GetMessageByID(int(id))
}

func UpdateMessage(id int, dto UpdateMessageDTO) *Message {
	action := fmt.Sprintf("UPDATE "+TABLE+" WHERE id : %d", id)

	_, err := database.Forum.Exec(
		"UPDATE "+TABLE+" SET content = ? WHERE id = ?",
		dto.Content, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}

	return GetMessageByID(id)
}

func DeleteMessage(id int) {
	action := fmt.Sprintf("DELETE FROM "+TABLE+" WHERE id : %d", id)

	_, err := database.Forum.Exec("DELETE FROM "+TABLE+" WHERE id = ?", id)
	if err != nil {
		log.Database(action, err)
	}
}

func GetMessageUsers(messageID int) []UserSummary {
	action := fmt.Sprintf("SELECT USERS JOIN USER_MESSAGE WHERE message_id : %d", messageID)

	rows, err := database.Forum.Query(
		"SELECT u.id, u.username, u.firstname, u.lastname FROM USERS u JOIN USER_MESSAGE um ON u.id = um.user_id WHERE um.message_id = ?",
		messageID,
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

func LinkUser(messageID int, userID string) {
	action := fmt.Sprintf("INSERT INTO USER_MESSAGE message_id:%d user_id:%s", messageID, userID)

	_, err := database.Forum.Exec(
		"INSERT IGNORE INTO USER_MESSAGE (user_id, message_id) VALUES (?, ?)",
		userID, messageID,
	)
	if err != nil {
		log.Database(action, err)
	}
}

func UnlinkUser(messageID int, userID string) {
	action := fmt.Sprintf("DELETE FROM USER_MESSAGE message_id:%d user_id:%s", messageID, userID)

	_, err := database.Forum.Exec(
		"DELETE FROM USER_MESSAGE WHERE user_id = ? AND message_id = ?",
		userID, messageID,
	)
	if err != nil {
		log.Database(action, err)
	}
}
