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

