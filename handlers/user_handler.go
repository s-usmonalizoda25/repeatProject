package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"go.uber.org/zap"
	"project/internal/models"
	"project/internal/service"
	"project/pkg/errs"
	"project/pkg/logger"
)
type UserHandler struct {
	logger      *logger.Logger
	userService *service.UserService
}
type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}
func New(logger *logger.Logger, userService *service.UserService) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
	}
}
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAll(r.Context())
	if err != nil {
		h.logger.Error("h.userService.GetAll", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		h.logger.Error("h.userService.Encode", zap.Error(err))
	}
}
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Ищем юзера через сервис
	userResp, err := h.userService.GetById(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrValidation):
			http.Error(w, "invalid request", http.StatusBadRequest)
		case errors.Is(err, errs.ErrUserNotFound):
			http.Error(w, "user not found", http.StatusNotFound)
		default:
			h.logger.Error("h.userService.GetByID", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	response := user{
		ID:   id,
		Name: userResp.Name,
		Age:  userResp.Age,
	}
	if err = json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("h.userService.GetByID Encode", zap.Error(err))
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var userRequest user
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.userService.Create(r.Context(), &models.User{
		ID:   userRequest.ID,
		Name: userRequest.Name,
		Age:  userRequest.Age,
	})
	if err != nil {
		if errors.Is(err, errs.ErrValidation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.logger.Error("h.userService.Create", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var userRequest user
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	err = h.userService.Update(r.Context(), &models.User{
		ID:   userRequest.ID,
		Name: userRequest.Name,
		Age:  userRequest.Age,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}