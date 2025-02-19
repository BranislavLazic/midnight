package postgres

import (
	"github.com/branislavlazic/midnight/model"
)

func (psr *Repository) CreateService(service *model.Service) (model.ServiceID, error) {
	resp := psr.db.Create(&service)
	return service.ID, resp.Error
}

func (psr *Repository) UpdateService(service *model.Service) error {
	return psr.db.Save(&service).Error
}

func (psr *Repository) GetAllServices() ([]model.Service, error) {
	services := make([]model.Service, 0)
	resp := psr.db.Order("id desc").Preload("Environment").Find(&services)
	return services, resp.Error
}

func (psr *Repository) GetServiceById(ID model.ServiceID) (*model.Service, error) {
	var service *model.Service
	resp := psr.db.Model(&model.Service{}).Preload("Environment").First(&service, ID)
	return service, resp.Error
}

func (psr *Repository) ServiceExistsByURL(URL string) bool {
	var service *model.Service
	return psr.db.First(&service, "url = ?", URL).Error == nil
}

func (psr *Repository) DeleteServiceById(ID model.ServiceID) error {
	service, err := psr.GetServiceById(ID)
	if err != nil {
		return err
	}
	return psr.db.Delete(&service, ID).Error
}

func (psr *Repository) DeleteAllServices() error {
	return psr.db.Delete(&model.Service{}, "id > 0").Error
}
