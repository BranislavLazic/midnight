package model

import "strings"

type EnvironmentID int64

type Environment struct {
	ID   EnvironmentID `gorm:"primaryKey,autoIncrement" json:"id"`
	Name string        `gorm:"type:VARCHAR(255)"`
}

type EnvironmentRequest struct {
	Name string `json:"name"`
}

func (cer *EnvironmentRequest) Sanitize() {
	cer.Name = strings.TrimSpace(cer.Name)
}

func (cer *EnvironmentRequest) ToPersistentEnvironment() *Environment {
	return &Environment{Name: cer.Name}
}

type EnvironmentRepository interface {
	Create(env *Environment) (EnvironmentID, error)
	GetAll() ([]Environment, error)
	DeleteAll() error
}
