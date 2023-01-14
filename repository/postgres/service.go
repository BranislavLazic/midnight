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

func (psr *ServiceRepository) Create(service *model.Service) (model.ServiceID, error) {
	resp := psr.db.Create(&service)
	return service.ID, resp.Error
}

func (psr *ServiceRepository) GetAll() ([]model.Service, error) {
	services := make([]model.Service, 0)
	resp := psr.db.Order("id desc").Find(&services)
	return services, resp.Error
}

func (psr *ServiceRepository) ExistsByURL(URL string) bool {
	var service *model.Service
	return psr.db.First(&service, "url = ?", URL).Error == nil
}
