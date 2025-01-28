package task

import (
	// "fmt"
	"fmt"
	"net/http"
	"strconv"
	"todo/types"
	"todo/utils"
	"todo/services/auth"

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

