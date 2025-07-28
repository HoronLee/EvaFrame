package service

import (
	"crypto/md5"
	"fmt"

	"evaframe/internal/dao/gorm"
	"evaframe/internal/models"
	"evaframe/pkg/config"
	"evaframe/pkg/jwt"
	"evaframe/pkg/validator"

	"go.uber.org/zap"
)

type UserService struct {
	userDAO   *gorm.UserDAO
	jwt       *jwt.JWT
	validator *validator.Validator
	logger    *zap.Logger
	config    *config.Config
}

func NewUserService(
	userDAO *gorm.UserDAO,
	jwt *jwt.JWT,
	validator *validator.Validator,
	logger *zap.Logger,
	config *config.Config,
) *UserService {
	return &UserService{
		userDAO:   userDAO,
		jwt:       jwt,
		validator: validator,
		logger:    logger,
		config:    config,
	}
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (s *UserService) Register(req *RegisterRequest) (*models.User, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	// 检查邮箱是否已存在
	if _, err := s.userDAO.GetByEmail(req.Email); err == nil {
		return nil, fmt.Errorf("email already exists")
	}

	// 密码加密（简单MD5，生产环境应使用bcrypt）
	hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.userDAO.Create(user); err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return nil, err
	}

	s.logger.Info("user registered successfully", zap.String("email", req.Email))
	return user, nil
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	// 查找用户
	user, err := s.userDAO.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 验证密码
	hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))
	if user.Password != hashedPassword {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 生成JWT token
	token, err := s.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		s.logger.Error("failed to generate token", zap.Error(err))
		return nil, err
	}

	s.logger.Info("user logged in successfully", zap.String("email", req.Email))
	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userDAO.GetByID(id)
}

func (s *UserService) ListUsers(offset, limit int) ([]*models.User, error) {
	return s.userDAO.List(offset, limit)
}
