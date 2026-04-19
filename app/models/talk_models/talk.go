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
	Id         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CategoryID int       `db:"category_id" json:"category_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type CreateTalkDTO struct {
	Title      string
	Content    string
	CategoryID int
}

type UpdateTalkDTO struct {
	Title   string
	Content string
}

type MessageSummary struct {
	Id        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type UserSummary struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func GetAllTalks(page, limit int) []Talk {
	action := "SELECT " + TABLE + " (paginated)"
	offset := (page - 1) * limit

	rows, err := database.Forum.Query(
		"SELECT id, title, content, category_id, created_at, updated_at FROM "+TABLE+" LIMIT ? OFFSET ?",
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
		if err := rows.Scan(&t.Id, &t.Title, &t.Content, &t.CategoryID, &t.CreatedAt, &t.UpdatedAt); err != nil {
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
		"SELECT id, title, content, category_id, created_at, updated_at FROM "+TABLE+" WHERE id = ?",
		id,
	)

	err := row.Scan(&talk.Id, &talk.Title, &talk.Content, &talk.CategoryID, &talk.CreatedAt, &talk.UpdatedAt)
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
		"INSERT INTO "+TABLE+" (title, content, category_id) VALUES (?, ?, ?)",
		dto.Title, dto.Content, dto.CategoryID,
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
		"UPDATE "+TABLE+" SET title = ?, content = ? WHERE id = ?",
		dto.Title, dto.Content, id,
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

func GetTalkMessages(talkID int) []MessageSummary {
	action := fmt.Sprintf("SELECT MESSAGES JOIN TALK_MESSAGE WHERE talk_id : %d", talkID)

	rows, err := database.Forum.Query(
		"SELECT m.id, m.content, m.created_at FROM MESSAGES m JOIN TALK_MESSAGE tm ON m.id = tm.message_id WHERE tm.talk_id = ?",
		talkID,
	)
	if err != nil {
		log.Database(action, err)
		return []MessageSummary{}
	}
	defer rows.Close()

	messages := []MessageSummary{}
	for rows.Next() {
		var m MessageSummary
		if err := rows.Scan(&m.Id, &m.Content, &m.CreatedAt); err != nil {
			log.Database(action, err)
			continue
		}
		messages = append(messages, m)
	}
	return messages
}

func LinkMessage(talkID, messageID int) {
	action := fmt.Sprintf("INSERT INTO TALK_MESSAGE talk_id:%d message_id:%d", talkID, messageID)

	_, err := database.Forum.Exec(
		"INSERT IGNORE INTO TALK_MESSAGE (talk_id, message_id) VALUES (?, ?)",
		talkID, messageID,
	)
	if err != nil {
		log.Database(action, err)
	}
}

func UnlinkMessage(talkID, messageID int) {
	action := fmt.Sprintf("DELETE FROM TALK_MESSAGE talk_id:%d message_id:%d", talkID, messageID)

	_, err := database.Forum.Exec(
		"DELETE FROM TALK_MESSAGE WHERE talk_id = ? AND message_id = ?",
		talkID, messageID,
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
