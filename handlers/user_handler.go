package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"project/internal/models"
	"project/internal/service"
	"project/pkg/errs"
	"project/pkg/logger"

	"go.uber.org/zap"
)

type UserHandler struct {
	logger      *logger.Logger
	userService *service.UserService
}

type userResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	IsActive bool   `json:"is_active"`
}

type loginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func New(logger *logger.Logger, userService *service.UserService) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
	}
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.CleanGetAll(r.Context())
	if err != nil {
		h.logger.Error("h.userService.CleanGetAll", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetAllArchive(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllArchive(r.Context())
	if err != nil {
		h.logger.Error("h.userService.GetAllArchive", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Login: Попытка входа с некорректным JSON")
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Info("Login: Попытка входа пользователя", zap.String("username", req.Name))

	err := h.userService.Login(r.Context(), req.Name, req.Password)
	if err != nil {
		h.logger.Warn("Login: Неудачная попытка входа", zap.String("username", req.Name), zap.Error(err))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	h.logger.Info("Login: Успешный вход в систему", zap.String("username", req.Name))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

func (h *UserHandler) SoftDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	h.logger.Info("SoftDelete: Запрос на деактивацию пользователя", zap.Int("user_id", id))

	err = h.userService.SoftDelete(r.Context(), id)
	if err != nil {
		h.logger.Warn("SoftDelete: Ошибка деактивации", zap.Int("user_id", id), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info("SoftDelete: Пользователь успешно деактивирован", zap.Int("user_id", id))
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userResp, err := h.userService.GetByID(r.Context(), id)
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

	response := userResponse{
		ID:       id,
		Name:     userResp.Name,
		Age:      userResp.Age,
		IsActive: userResp.IsActive,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var userRequest models.User
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.userService.Create(r.Context(), &userRequest)
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
	var userRequest models.User
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.userService.Update(r.Context(), &userRequest)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrValidation):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, errs.ErrUserNotFound):
			http.Error(w, "user not found", http.StatusNotFound)
		default:
			h.logger.Error("h.userService.Update", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
}
