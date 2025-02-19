package postgres

import (
	"github.com/branislavlazic/midnight/model"
)

func (er *Repository) CreateEnvironment(env *model.Environment) (model.EnvironmentID, error) {
	resp := er.db.Create(&env)
	return env.ID, resp.Error
}

func (er *Repository) GetAllEnvironments() ([]model.Environment, error) {
	envs := make([]model.Environment, 0)
	resp := er.db.Order("name asc").Find(&envs)
	return envs, resp.Error
}

func (er *Repository) GetEnvironmentByID(ID model.EnvironmentID) (*model.Environment, error) {
	var environment *model.Environment
	resp := er.db.First(&environment, ID)
	return environment, resp.Error
}

func (er *Repository) DeleteAllEnvironments() error {
	return er.db.Delete(&model.Environment{}, "id > 0").Error
}
