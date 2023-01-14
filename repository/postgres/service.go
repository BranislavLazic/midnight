package postgres

import (
	"github.com/branislavlazic/midnight/model"
	"gorm.io/gorm"
)

type ServiceRepository struct {
	db *gorm.DB
}

func NewPgServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}
func (psr *ServiceRepository) GetAll() ([]model.Service, error) {
	services := make([]model.Service, 0)
	resp := psr.db.Order("id desc").Find(&services)
	return services, resp.Error
}
