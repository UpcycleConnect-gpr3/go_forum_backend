package project_models

import (
	"database/sql"
	"fmt"
	"go-forum-backend/database"
	"go-forum-backend/utils/log"
	"time"
)

const TABLE = "PROJECTS"

type Project struct {
	Id          int            `json:"id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
}

type CreateProjectDTO struct {
	Title       string
	Description string
}

type UpdateProjectDTO struct {
	Title       string
	Description string
}

func GetAllProjects(page, limit int) []Project {
	action := "SELECT " + TABLE + " (paginated)"
	offset := (page - 1) * limit

	rows, err := database.Forum.Query(
		"SELECT id, title, description, created_at, updated_at FROM "+TABLE+" LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		log.Database(action, err)
		return []Project{}
	}
	defer rows.Close()

	projects := []Project{}
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.Id, &p.Title, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			log.Database(action, err)
			continue
		}
		projects = append(projects, p)
	}
	return projects
}

func GetProjectByID(id int) *Project {
	project := Project{}
	action := fmt.Sprintf("SELECT "+TABLE+" WHERE id : %d", id)

	row := database.Forum.QueryRow(
		"SELECT id, title, description, created_at, updated_at FROM "+TABLE+" WHERE id = ?",
		id,
	)

	err := row.Scan(&project.Id, &project.Title, &project.Description, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Database(action, err)
		return nil
	}

	return &project
}

func CreateProject(dto CreateProjectDTO) *Project {
	action := "INSERT INTO " + TABLE

	result, err := database.Forum.Exec(
		"INSERT INTO "+TABLE+" (title, description) VALUES (?, ?)",
		dto.Title,
		sql.NullString{String: dto.Description, Valid: dto.Description != ""},
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

	return GetProjectByID(int(id))
}

func UpdateProject(id int, dto UpdateProjectDTO) *Project {
	action := fmt.Sprintf("UPDATE "+TABLE+" WHERE id : %d", id)

	_, err := database.Forum.Exec(
		"UPDATE "+TABLE+" SET title = ?, description = ? WHERE id = ?",
		dto.Title,
		sql.NullString{String: dto.Description, Valid: dto.Description != ""},
		id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}

	return GetProjectByID(id)
}

func DeleteProject(id int) {
	action := fmt.Sprintf("DELETE FROM "+TABLE+" WHERE id : %d", id)

	_, err := database.Forum.Exec("DELETE FROM "+TABLE+" WHERE id = ?", id)
	if err != nil {
		log.Database(action, err)
	}
}
