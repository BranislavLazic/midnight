package model

import "strings"

type ServiceID int64

type Service struct {
	ID                   ServiceID      `gorm:"primaryKey;autoIncrement" json:"id"`
	EnvironmentID        *EnvironmentID `json:"-"`
	Environment          *Environment   `gorm:"foreignKey:EnvironmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"environment,omitempty"`
	Name                 string         `gorm:"type:VARCHAR(255)" json:"name"`
	URL                  string         `json:"url"`
	ResponseBody         string         `json:"responseBody"`
	CheckIntervalSeconds int            `json:"checkIntervalSeconds"`
}

type ServiceRequest struct {
	Name                 string         `json:"name" validate:"max=255"`
	EnvironmentID        *EnvironmentID `json:"environmentId"`
	URL                  string         `json:"url" validate:"required,max=4096"`
	ResponseBody         string         `json:"responseBody" validate:"max=8192"`
	CheckIntervalSeconds int            `json:"checkIntervalSeconds" validate:"required,max=1000000"`
}

func (sr *ServiceRequest) Sanitize() {
	sr.Name = strings.TrimSpace(sr.Name)
	sr.URL = strings.TrimSpace(sr.URL)
	sr.ResponseBody = strings.TrimSpace(sr.ResponseBody)
}

func (sr *ServiceRequest) ToPersistentService(ID ServiceID, env *Environment) *Service {
	return &Service{
		ID:                   ID,
		Name:                 sr.Name,
		Environment:          env,
		URL:                  sr.URL,
		ResponseBody:         sr.ResponseBody,
		CheckIntervalSeconds: sr.CheckIntervalSeconds,
	}
}

type ServiceRepository interface {
	Create(service *Service) (ServiceID, error)
	Update(service *Service) error
	GetAll() ([]Service, error)
	GetById(ID ServiceID) (*Service, error)
	DeleteById(ID ServiceID) error
	DeleteAll() error
	ExistsByURL(URL string) bool
}
