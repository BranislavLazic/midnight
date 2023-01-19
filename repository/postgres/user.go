package postgres

import (
	"github.com/branislavlazic/midnight/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(user *model.User) (model.UserID, error) {
	resp := ur.db.Create(&user)
	return user.ID, resp.Error
}

func (ur *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user *model.User
	resp := ur.db.First(&user, "email = ?", email)
	return user, resp.Error
}
