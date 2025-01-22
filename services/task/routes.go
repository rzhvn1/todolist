package task

import (
	// "fmt"
	"todo/types"

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
}