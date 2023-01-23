package postgres

import (
	"github.com/branislavlazic/midnight/model"
	"gorm.io/gorm"
)

type EnvironmentRepository struct {
	db *gorm.DB
}

func NewEnvironmentRepository(db *gorm.DB) *EnvironmentRepository {
	return &EnvironmentRepository{db: db}
}

func (er *EnvironmentRepository) Create(env *model.Environment) (model.EnvironmentID, error) {
	resp := er.db.Create(&env)
	return env.ID, resp.Error
}

func (er *EnvironmentRepository) GetAll() ([]model.Environment, error) {
	envs := make([]model.Environment, 0)
	resp := er.db.Order("name asc").Find(&envs)
	return envs, resp.Error
}

func (er *EnvironmentRepository) GetDefault() (*model.Environment, error) {
	var environment *model.Environment
	resp := er.db.First(&environment, "name = ?", model.DefaultEnvironmentName)
	return environment, resp.Error
}

func (er *EnvironmentRepository) DeleteAll() error {
	return er.db.Delete(&model.Environment{}, "id > 0").Error
}
