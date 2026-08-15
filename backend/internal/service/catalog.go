package service

import (
	"context"
	"fmt"

	"batam-medhub/internal/model"

	"gorm.io/gorm"
)

// CatalogService provides operations for querying supported medical services.
type CatalogService struct {
	db *gorm.DB
}

// NewCatalogService constructs a CatalogService with the given database handle.
func NewCatalogService(db *gorm.DB) *CatalogService {
	return &CatalogService{db: db}
}

// MedicalServiceItem represents a supported medical check-up or screening service.
type MedicalServiceItem struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Description *string `json:"description"`
	Active      bool    `json:"active"`
	Synthetic   bool    `json:"synthetic"`
	Source      string  `json:"source"`
}

// MedicalServiceListResponse wraps a list of active medical services.
type MedicalServiceListResponse struct {
	Services []MedicalServiceItem `json:"services"`
}

// ListMedicalServices queries and returns all active medical services.
func (s *CatalogService) ListMedicalServices(ctx context.Context) (*MedicalServiceListResponse, error) {
	var records []model.MedicalService
	if err := s.db.WithContext(ctx).Where("active = ?", true).Order("code ASC").Find(&records).Error; err != nil {
		return nil, fmt.Errorf("list medical services: %w", err)
	}

	items := make([]MedicalServiceItem, len(records))
	for i, r := range records {
		items[i] = MedicalServiceItem{
			Code:        r.Code,
			Name:        r.Name,
			Category:    r.Category,
			Description: r.Description,
			Active:      r.Active,
			Synthetic:   r.Synthetic,
			Source:      r.Source,
		}
	}

	return &MedicalServiceListResponse{Services: items}, nil
}
