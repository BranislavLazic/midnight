package postgres

import (
	"github.com/branislavlazic/midnight/model"
)

func (ur *Repository) CreateUser(user *model.User) (model.UserID, error) {
	resp := ur.db.Create(&user)
	return user.ID, resp.Error
}

func (ur *Repository) GetUserByEmail(email string) (*model.User, error) {
	var user *model.User
	resp := ur.db.First(&user, "email = ?", email)
	return user, resp.Error
}
