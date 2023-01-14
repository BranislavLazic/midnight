package model

type ServiceID int64

type Service struct {
	ID                   ServiceID `gorm:"primaryKey;autoIncrement"`
	Name                 string    `gorm:"type:VARCHAR(255)"`
	URL                  string
	CheckIntervalSeconds int
}

type ServiceRepository interface {
	GetAll() ([]Service, error)
}
