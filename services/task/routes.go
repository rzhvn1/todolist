package task

import (
	// "fmt"
	"fmt"
	"net/http"
	"strconv"
	"todo/services/auth"
	"todo/types"
	"todo/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.TaskStore
	userStore types.UserStore
}

func NewHandler(store types.TaskStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/tasks", auth.WithJWTAuth(h.handleGetTasks, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/tasks", auth.WithJWTAuth(h.handleCreateTask, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/tasks/{task_id}", auth.WithJWTAuth(h.handleUpdateTask, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/tasks/{task_id}", auth.WithJWTAuth(h.handleDeleteTask, h.userStore)).Methods(http.MethodDelete)
}

func (h *Handler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	allowedSortFields := []string{"user_id", "status", "priority", "due_date"}
	pagination, err := utils.ParsePaginationParams(r, allowedSortFields)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// get paginated data from store
	tasks, total, err := h.store.GetPaginatedTasks(pagination)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get tasks: %v", err))
		return
	}

	utils.WritePaginatedResponse(w, pagination.Page, pagination.Limit, total, tasks)
}

func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var task types.CreateTaskPayload
	if err := utils.ParseJSON(r, &task); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(task); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	err := h.store.CreateTask(task)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, task)
}

func (h *Handler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var task types.UpdateTaskPayload

	vars := mux.Vars(r)
	str, ok := vars["task_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing task ID"))
		return
	}

	taskID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid task ID"))
		return
	}
	existingTask, err := h.store.GetTaskByID(taskID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	if err := utils.ParseJSON(r, &task); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(task); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	if task.UserID != nil {
		_, err := h.userStore.GetUserByID(*task.UserID)
		if err != nil {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user not found"))
			return
		}

	}

	if task.Title == nil {
		task.Title = &existingTask.Title
	}
	if task.Description == nil {
		task.Description = &existingTask.Description
	}
	if task.Status == nil {
		task.Status = &existingTask.Status
	}
	if task.Priority == nil {
		task.Priority = &existingTask.Priority
	}
	if task.DueDate == nil {
		task.DueDate = existingTask.DueDate
	}

	err = h.store.UpdateTask(taskID, task)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	updatedTask, _ := h.store.GetTaskByID(taskID)
	utils.WriteJson(w, http.StatusOK, updatedTask)
}

func (h *Handler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["task_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing task ID"))
		return
	}

	taskID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid task ID"))
		return
	}

	rowsAffected, err := h.store.DeleteTask(taskID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete task: %v", err))
		return
	}
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("task not found"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
