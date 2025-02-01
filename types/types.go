package types

import (
	"time"
	"todo/utils"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type UserStore interface {
	GetUserByID(userID int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(User) error
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterUserPayload struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
}

type Task struct {
	ID          int        `json:"id"`
	UserID      *int       `json:"user_id"` // Fixed tag
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`   // pending, in_progress, completed
	Priority    int        `json:"priority"` // 1 - low, 2 - medium, 3 - high
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TaskStore interface {
	GetTaskByID(taskID int) (*Task, error)
	GetPaginatedTasks(pagination utils.PaginationParams) ([]Task, int, error)
	CreateTask(task CreateTaskPayload) error
	UpdateTask(taskID int, task UpdateTaskPayload) error
	DeleteTask(taskID int) (int64, error)
}

type CreateTaskPayload struct {
	UserID      *int       `json:"user_id"`
	Title       string     `json:"title" validate:"required"`
	Description *string     `json:"description"`
	Status      string     `json:"status,omitempty" validate:"omitempty,oneof=pending in_progress completed"`
	Priority    int        `json:"priority" validate:"required"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskPayload struct {
	UserID      *int       `json:"user_id,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty" validate:"omitempty,oneof=pending in_progress completed"`
	Priority    *int       `json:"priority,omitempty" validate:"omitempty,oneof=1 2 3"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}
