package task

import (
	"database/sql"
	"fmt"
	"todo/types"
	"todo/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetTaskByID(taskID int) (*types.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks WHERE id = ?", taskID)
	if err != nil {
		return nil, err
	}

	t := new(types.Task)
	for rows.Next() {
		t, err = scanRowsIntoTask(rows)
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (s *Store) GetPaginatedTasks(pagination utils.PaginationParams) ([]types.Task, int, error) {
	// get total count
	var total int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&total); err != nil {
		return nil, 0, err
	}

	// dynamic query with safe params
	query := fmt.Sprintf(`
	SELECT id, user_id, title, description, status, priority, due_date, created_at, updated_at
	FROM tasks
	ORDER BY %s %s
	LIMIT ? OFFSET ?`, pagination.SortBy, pagination.Order)

	rows, err := s.db.Query(query, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// map db rows to struct
	var tasks []types.Task
	for rows.Next() {
		t, err := scanRowsIntoTask(rows)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, *t)
	}

	return tasks, total, nil
}

func (s *Store) CreateTask(task types.CreateTaskPayload) error {
	if task.Status == "" {
		task.Status = "pending"
	}

	_, err := s.db.Exec(
		"INSERT INTO tasks (user_id, title, description, status, priority, due_date) VALUES (?, ?, ?, ?, ?, ?)",
		task.UserID, task.Title, task.Description, task.Status, task.Priority, task.DueDate)

	return err
}

func (s *Store) UpdateTask(taskID int, task types.UpdateTaskPayload) error {
	_, err := s.db.Exec(
		"UPDATE tasks SET user_id = ?, title = ?, description = ?, status = ?, priority = ?, due_date = ? WHERE id = ?",
		task.UserID, task.Title, task.Description, task.Status, task.Priority, task.DueDate, taskID)

	return err
}

func (s *Store) DeleteTask(taskID int) (int64, error) {
	result, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
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
