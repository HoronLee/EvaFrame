package gorm

import (
	"evaframe/internal/models"
	"evaframe/internal/service"

	"gorm.io/gorm"
)

// UserDAOImpl 实现 service.UserDAO 接口
type UserDAOImpl struct {
	db *gorm.DB
}

// NewUserDAO 返回接口类型
func NewUserDAO(db *gorm.DB) service.UserDAO {
	return &UserDAOImpl{db: db}
}

func (d *UserDAOImpl) Create(user *models.User) error {
	return d.db.Create(user).Error
}

func (d *UserDAOImpl) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := d.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDAOImpl) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := d.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDAOImpl) List(offset, limit int) ([]*models.User, error) {
	var users []*models.User
	err := d.db.Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}
