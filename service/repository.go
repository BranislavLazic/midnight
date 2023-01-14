package service

import "gorm.io/gorm"

type PgServiceRepository struct {
	db *gorm.DB
}

func NewPgServiceRepository(db *gorm.DB) *PgServiceRepository {
	return &PgServiceRepository{db: db}
}
func (psr *PgServiceRepository) GetAll() ([]Service, error) {
	services := make([]Service, 0)
	resp := psr.db.Order("id desc").Find(&services)
	return services, resp.Error
}
