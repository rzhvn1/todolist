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
	store types.TaskStore
	userStore types.UserStore
}

func NewHandler(store types.TaskStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/tasks/user/{user_id}", auth.WithJWTAuth(h.handleGetTasksByUserID, h.userStore)).Methods(http.MethodGet)

	router.HandleFunc("/tasks", auth.WithJWTAuth(h.handleCreateTask, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/tasks/{task_id}", auth.WithJWTAuth(h.handleUpdateTask, h.userStore)).Methods(http.MethodPut)
}

func (h *Handler) handleGetTasksByUserID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	str, ok := vars["user_id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	userID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	user, err := h.store.GetTasksByUserID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, user)
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

