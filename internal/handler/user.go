package handler

import (
	"strconv"

	"evaframe/internal/service"
	"evaframe/pkg/response"
	"evaframe/pkg/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService *service.UserService
	val   *validator.Validator
	logger      *zap.Logger
}

func NewUserHandler(userService *service.UserService, validator *validator.Validator, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		val:   validator,
		logger:      logger,
	}
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证请求数据
	if err := h.val.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用业务逻辑层
	user, err := h.userService.CreateUser(req.Name, req.Email, req.Password)
	if err != nil {
		h.logger.Error("register failed", zap.Error(err))
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, user)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User  any    `json:"user"`
	Token string `json:"token"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证请求数据
	if err := h.val.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用业务逻辑层
	user, token, err := h.userService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		h.logger.Error("login failed", zap.Error(err))
		response.Error(c, 401, err.Error())
		return
	}

	// 构造响应
	result := LoginResponse{
		User:  user,
		Token: token,
	}

	response.Success(c, result)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		h.logger.Error("get profile failed", zap.Error(err))
		response.Error(c, 404, "user not found")
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, err := h.userService.ListUsers(offset, limit)
	if err != nil {
		h.logger.Error("list users failed", zap.Error(err))
		response.InternalError(c, "failed to list users")
		return
	}

	response.Success(c, users)
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine, authMiddleware gin.HandlerFunc) {
	api := r.Group("/api/v1")
	{
		// 公开路由
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)

		// 需要认证的路由
		auth := api.Group("/", authMiddleware)
		{
			auth.GET("/profile", h.GetProfile)
			auth.GET("/users", h.ListUsers)
		}
	}
}
