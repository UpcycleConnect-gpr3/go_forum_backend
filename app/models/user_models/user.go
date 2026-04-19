package user_models

import (
	"database/sql"
	"fmt"
	"go-forum-backend/database"
	"go-forum-backend/utils/log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const TABLE = "USERS"

type User struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	password  string    `db:"password" json:"-"`
	Email     string    `json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserDTO struct {
	Username  string
	Firstname string
	Lastname  string
}

func (u *User) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.password = string(hashed)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	return err == nil
}

func GetAllUsers(page, limit int) []User {
	action := "SELECT " + TABLE + " (paginated)"
	offset := (page - 1) * limit

	rows, err := database.Forum.Query(
		"SELECT id, username, firstname, lastname, email, created_at, updated_at FROM "+TABLE+" LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		log.Database(action, err)
		return []User{}
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Username, &u.Firstname, &u.Lastname, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			log.Database(action, err)
			continue
		}
		users = append(users, u)
	}
	return users
}

func GetUserByID(id string) *User {
	user := User{}
	action := fmt.Sprintf("SELECT "+TABLE+" WHERE id : %s", id)

	row := database.Forum.QueryRow(
		"SELECT id, username, firstname, lastname, email, created_at, updated_at FROM "+TABLE+" WHERE id = ?",
		id,
	)

	err := row.Scan(&user.Id, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Database(action, err)
		return nil
	}

	return &user
}

func GetUserByEmail(email string) *User {
	user := User{}
	action := fmt.Sprintf("SELECT "+TABLE+" WHERE email : %s", email)

	row := database.Forum.QueryRow(
		"SELECT id, email, password FROM "+TABLE+" WHERE email = ?",
		email,
	)

	err := row.Scan(&user.Id, &user.Email, &user.password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Database(action, err)
		return nil
	}

	return &user
}

func CreateUser(credentials Credentials) *User {
	action := fmt.Sprintf("INSERT INTO "+TABLE+" : %s", credentials.Email)

	u := User{}
	if err := u.SetPassword(credentials.Password); err != nil {
		log.Database(action, err)
		return nil
	}

	id := uuid.New()
	_, err := database.Forum.Exec(
		"INSERT INTO "+TABLE+" (id, email, password) VALUES (?, ?, ?)",
		id.String(), credentials.Email, u.password,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}

	return GetUserByID(id.String())
}

func UpdateUser(id string, dto UpdateUserDTO) *User {
	action := fmt.Sprintf("UPDATE "+TABLE+" WHERE id : %s", id)

	_, err := database.Forum.Exec(
		"UPDATE "+TABLE+" SET username = ?, firstname = ?, lastname = ? WHERE id = ?",
		dto.Username, dto.Firstname, dto.Lastname, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}

	return GetUserByID(id)
}

func DeleteUser(id string) {
	action := fmt.Sprintf("DELETE FROM "+TABLE+" WHERE id : %s", id)

	_, err := database.Forum.Exec("DELETE FROM "+TABLE+" WHERE id = ?", id)
	if err != nil {
		log.Database(action, err)
	}
}

func GetUserMessages(userID string) []MessageSummary {
	action := fmt.Sprintf("SELECT MESSAGES JOIN USER_MESSAGE WHERE user_id : %s", userID)

	rows, err := database.Forum.Query(
		"SELECT m.id, m.content, m.created_at FROM MESSAGES m JOIN USER_MESSAGE um ON m.id = um.message_id WHERE um.user_id = ?",
		userID,
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

type MessageSummary struct {
	Id        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
