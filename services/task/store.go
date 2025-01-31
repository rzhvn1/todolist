package task

import (
	"database/sql"
	"fmt"
	"todo/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetTasksByUserID(userID int) ([]*types.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	tasks := make([]*types.Task, 0)
	for rows.Next() {
		t, err := scanRowsIntoTask(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (s *Store) CreateTask(task types.CreateTaskPayload) error {
	if task.Status == "" {
		task.Status = "pending"
	}
	
	_, err := s.db.Exec(
		"INSERT INTO tasks (user_id, title, description, status, priority, due_date) VALUES (?, ?, ?, ?, ?, ?)",
		task.UserID, task.Title, task.Description, task.Status, task.Priority, task.DueDate)
	fmt.Println("HEre")
	if err != nil {
		return err
	}
	fmt.Println("NOT HERE")
	
	return nil
}

func scanRowsIntoTask(rows *sql.Rows) (*types.Task, error) {
	task := new(types.Task)

	err := rows.Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return task, nil
}
