package storage

import (
	"github.com/cesc1802/onboarding-and-volunteer-service/feature/user/domain"

	"gorm.io/gorm"
)

type ApplicantRequestRepository struct {
	DB *gorm.DB
}

func NewApplicantRequestRepository(db *gorm.DB) *ApplicantRequestRepository {
	return &ApplicantRequestRepository{DB: db}
}

func (r *ApplicantRequestRepository) CreateApplicantRequest(request *domain.ApplicantRequestDomain) error {
	return r.DB.Create(request).Error
}