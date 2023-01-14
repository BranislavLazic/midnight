package model

import "strings"

type ServiceID int64

type Service struct {
	ID                   ServiceID `gorm:"primaryKey;autoIncrement"`
	Name                 string    `gorm:"type:VARCHAR(255)"`
	URL                  string
	CheckIntervalSeconds int
}

type CreateServiceRequest struct {
	Name                 string `json:"name" validate:"max=255"`
	URL                  string `json:"url" validate:"required,max=4096"`
	CheckIntervalSeconds int    `json:"checkIntervalSeconds" validate:"required,max=1000000"`
}

func (csr *CreateServiceRequest) Sanitize() {
	csr.Name = strings.TrimSpace(csr.Name)
	csr.URL = strings.TrimSpace(csr.URL)
}

func (csr *CreateServiceRequest) ToPersistentService() *Service {
	return &Service{
		Name:                 csr.Name,
		URL:                  csr.URL,
		CheckIntervalSeconds: csr.CheckIntervalSeconds,
	}
}

type ServiceRepository interface {
	Create(service *Service) (ServiceID, error)
	GetAll() ([]Service, error)
	ExistsByURL(URL string) bool
}
