package service

import (
	"crypto/md5"
	"fmt"

	"evaframe/internal/models"
	"evaframe/pkg/config"
	"evaframe/pkg/jwt"
	"evaframe/pkg/validator"

	"go.uber.org/zap"
)

// UserDAO 接口定义 - Service 层定义需要的数据访问方法
type UserDAO interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	List(offset, limit int) ([]*models.User, error)
}

type UserService struct {
	userDAO   UserDAO
	jwt       *jwt.JWT
	validator *validator.Validator
	logger    *zap.Logger
	config    *config.Config
}

func NewUserService(
	userDAO UserDAO,
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

// 业务逻辑方法 - 直接使用领域对象
func (s *UserService) CreateUser(name, email, password string) (*models.User, error) {
	// 检查邮箱是否已存在
	if _, err := s.userDAO.GetByEmail(email); err == nil {
		return nil, fmt.Errorf("email already exists")
	}

	// 密码加密（简单MD5，生产环境应使用bcrypt）
	hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(password)))

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	if err := s.userDAO.Create(user); err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return nil, err
	}

	s.logger.Info("user registered successfully", zap.String("email", email))
	return user, nil
}

func (s *UserService) AuthenticateUser(email, password string) (*models.User, string, error) {
	// 查找用户
	user, err := s.userDAO.GetByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// 验证密码
	hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	if user.Password != hashedPassword {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// 生成JWT token
	token, err := s.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		s.logger.Error("failed to generate token", zap.Error(err))
		return nil, "", err
	}

	s.logger.Info("user logged in successfully", zap.String("email", email))
	return user, token, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userDAO.GetByID(id)
}

func (s *UserService) ListUsers(offset, limit int) ([]*models.User, error) {
	return s.userDAO.List(offset, limit)
}
