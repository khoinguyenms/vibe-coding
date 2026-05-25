package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/vibe-be/internal/model"
	"github.com/vibe-be/internal/service"
	"github.com/vibe-be/pkg/logger"
	"github.com/vibe-be/pkg/response"
)

type userHandler struct {
	svc service.UserService
	log *logger.Logger
}

func NewUserHandler(svc service.UserService, log *logger.Logger) UserHandler {
	return &userHandler{svc: svc, log: log.Named("user-handler")}
}

func (h *userHandler) Create(ctx *gin.Context) {
	var req model.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	user, err := h.svc.Create(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	user.Password = ""
	response.Success(ctx, http.StatusCreated, user)
}

func (h *userHandler) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(ctx, http.StatusNotFound, "user not found")
			return
		}
		h.log.Ctx(ctx.Request.Context()).Error("get user", zap.Error(err))
		response.Error(ctx, http.StatusInternalServerError, "internal error")
		return
	}

	user.Password = ""
	response.Success(ctx, http.StatusOK, user)
}

func (h *userHandler) List(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	offset, _ := strconv.Atoi(ctx.Query("offset"))

	users, err := h.svc.List(ctx.Request.Context(), int32(limit), int32(offset))
	if err != nil {
		h.log.Ctx(ctx.Request.Context()).Error("list users", zap.Error(err))
		response.Error(ctx, http.StatusInternalServerError, "internal error")
		return
	}

	for i := range users {
		users[i].Password = ""
	}
	response.Success(ctx, http.StatusOK, users)
}
