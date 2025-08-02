package service

import (
	"crypto/md5"
	"fmt"

	"evaframe/internal/models"
	"evaframe/pkg/config"
	"evaframe/pkg/jwt"
	"evaframe/pkg/logger"
	"evaframe/pkg/validator"
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
	logger    *logger.Logger
	config    *config.Config
}

func NewUserService(
	userDAO UserDAO,
	jwt *jwt.JWT,
	validator *validator.Validator,
	logger *logger.Logger,
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
		s.logger.LogIf("failed to create user", err)
		return nil, err
	}

	s.logger.InfoString("user", "user registered successfully", email)
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
		s.logger.LogIf("failed to generate token", err)
		return nil, "", err
	}

	s.logger.InfoString("user", "user logged in successfully", email)
	return user, token, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userDAO.GetByID(id)
}

func (s *UserService) ListUsers(offset, limit int) ([]*models.User, error) {
	return s.userDAO.List(offset, limit)
}
