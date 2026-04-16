package project_models

import (
	"database/sql"
	"fmt"
	"go-forum-backend/database"
	"go-forum-backend/utils/log"
)

const TABLE = "PROJECTS"

type Project struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetProjectByID(id int) *Project {
	project := Project{}
	action := fmt.Sprintf("SELECT "+TABLE+" WHERE id : %d", id)

	row := database.Forum.QueryRow("SELECT id, name FROM "+TABLE+" WHERE id = ?", id)

	err := row.Scan(&project.Id, &project.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Database(action, err)
		return nil
	}

	return &project
}
