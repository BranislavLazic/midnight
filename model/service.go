package model

import "strings"

type ServiceID int64

type Service struct {
	ID                   ServiceID `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                 string    `gorm:"type:VARCHAR(255)" json:"name"`
	URL                  string    `json:"url"`
	CheckIntervalSeconds int       `json:"checkIntervalSeconds"`
}

type ServiceRequest struct {
	Name                 string `json:"name" validate:"max=255"`
	URL                  string `json:"url" validate:"required,max=4096"`
	CheckIntervalSeconds int    `json:"checkIntervalSeconds" validate:"required,max=1000000"`
}

func (sr *ServiceRequest) Sanitize() {
	sr.Name = strings.TrimSpace(sr.Name)
	sr.URL = strings.TrimSpace(sr.URL)
}

func (sr *ServiceRequest) ToPersistentService(ID ServiceID) *Service {
	return &Service{
		ID:                   ID,
		Name:                 sr.Name,
		URL:                  sr.URL,
		CheckIntervalSeconds: sr.CheckIntervalSeconds,
	}
}

type ServiceRepository interface {
	Create(service *Service) (ServiceID, error)
	Update(service *Service) error
	GetAll() ([]Service, error)
	GetById(ID ServiceID) (*Service, error)
	DeleteById(ID ServiceID) error
	ExistsByURL(URL string) bool
}
