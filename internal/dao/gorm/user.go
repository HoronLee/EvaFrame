package gorm

import (
	"evaframe/internal/models"

	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (d *UserDAO) Create(user *models.User) error {
	return d.db.Create(user).Error
}

func (d *UserDAO) Update(user *models.User) error {
	return d.db.Save(user).Error
}

func (d *UserDAO) Delete(id uint) error {
	return d.db.Delete(&models.User{}, id).Error
}

func (d *UserDAO) List(offset, limit int) ([]*models.User, error) {
	var users []*models.User
	err := d.db.Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

func (d *UserDAO) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := d.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDAO) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := d.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
