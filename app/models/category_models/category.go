package category_models

import (
	"database/sql"
	"fmt"
	"go-forum-backend/database"
	"go-forum-backend/utils/log"
	"time"
)

const TABLE = "CATEGORIES"

type Category struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type CreateCategoryDTO struct {
	Name        string
	Description string
}

type UpdateCategoryDTO struct {
	Name        string
	Description string
}

type TalkSummary struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

func GetAllCategories(page, limit int) []Category {
	action := "SELECT " + TABLE + " (paginated)"
	offset := (page - 1) * limit

	rows, err := database.Forum.Query(
		"SELECT id, name, description, created_at, updated_at FROM "+TABLE+" LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		log.Database(action, err)
		return []Category{}
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.Id, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			log.Database(action, err)
			continue
		}
		categories = append(categories, c)
	}
	return categories
}

func GetCategoryByID(id int) *Category {
	category := Category{}
	action := fmt.Sprintf("SELECT "+TABLE+" WHERE id : %d", id)

	row := database.Forum.QueryRow(
		"SELECT id, name, description, created_at, updated_at FROM "+TABLE+" WHERE id = ?",
		id,
	)

	err := row.Scan(&category.Id, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Database(action, err)
		return nil
	}

	return &category
}

func CreateCategory(dto CreateCategoryDTO) *Category {
	action := "INSERT INTO " + TABLE

	result, err := database.Forum.Exec(
		"INSERT INTO "+TABLE+" (name, description) VALUES (?, ?)",
		dto.Name, dto.Description,
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

	return GetCategoryByID(int(id))
}

func UpdateCategory(id int, dto UpdateCategoryDTO) *Category {
	action := fmt.Sprintf("UPDATE "+TABLE+" WHERE id : %d", id)

	_, err := database.Forum.Exec(
		"UPDATE "+TABLE+" SET name = ?, description = ? WHERE id = ?",
		dto.Name, dto.Description, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}

	return GetCategoryByID(id)
}

func DeleteCategory(id int) {
	action := fmt.Sprintf("DELETE FROM "+TABLE+" WHERE id : %d", id)

	_, err := database.Forum.Exec("DELETE FROM "+TABLE+" WHERE id = ?", id)
	if err != nil {
		log.Database(action, err)
	}
}

func GetCategoryTalks(categoryID int) []TalkSummary {
	action := fmt.Sprintf("SELECT TALKS WHERE category_id : %d", categoryID)

	rows, err := database.Forum.Query(
		"SELECT id, title, created_at FROM TALKS WHERE category_id = ?",
		categoryID,
	)
	if err != nil {
		log.Database(action, err)
		return []TalkSummary{}
	}
	defer rows.Close()

	talks := []TalkSummary{}
	for rows.Next() {
		var t TalkSummary
		if err := rows.Scan(&t.Id, &t.Title, &t.CreatedAt); err != nil {
			log.Database(action, err)
			continue
		}
		talks = append(talks, t)
	}
	return talks
}
