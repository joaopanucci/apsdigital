package services

import (
	"context"
	"fmt"

	"apsdigital/internal/domain/entities"
	"apsdigital/internal/domain/repositories"
)

type MunicipalityService struct {
	municipalityRepo repositories.MunicipalityRepository
}

func NewMunicipalityService(municipalityRepo repositories.MunicipalityRepository) *MunicipalityService {
	return &MunicipalityService{
		municipalityRepo: municipalityRepo,
	}
}

func (s *MunicipalityService) GetAll(ctx context.Context) ([]*entities.Municipality, error) {
	municipalities, err := s.municipalityRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get municipalities: %w", err)
	}
	
	return municipalities, nil
}

func (s *MunicipalityService) GetByID(ctx context.Context, id int) (*entities.Municipality, error) {
	municipality, err := s.municipalityRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get municipality: %w", err)
	}
	
	return municipality, nil
}

func (s *MunicipalityService) GetByName(ctx context.Context, name string) (*entities.Municipality, error) {
	municipality, err := s.municipalityRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get municipality by name: %w", err)
	}
	
	return municipality, nil
}

func (s *MunicipalityService) Create(ctx context.Context, municipality *entities.Municipality) error {
	// Check if municipality with this name already exists
	existing, _ := s.municipalityRepo.GetByName(ctx, municipality.Name)
	if existing != nil {
		return fmt.Errorf("municipality with name '%s' already exists", municipality.Name)
	}
	
	if err := s.municipalityRepo.Create(ctx, municipality); err != nil {
		return fmt.Errorf("failed to create municipality: %w", err)
	}
	
	return nil
}

func (s *MunicipalityService) Update(ctx context.Context, municipality *entities.Municipality) error {
	// Check if municipality exists
	existing, err := s.municipalityRepo.GetByID(ctx, municipality.ID)
	if err != nil {
		return fmt.Errorf("municipality not found")
	}
	
	// Check if another municipality with this name exists (excluding current one)
	if existing.Name != municipality.Name {
		nameCheck, _ := s.municipalityRepo.GetByName(ctx, municipality.Name)
		if nameCheck != nil && nameCheck.ID != municipality.ID {
			return fmt.Errorf("municipality with name '%s' already exists", municipality.Name)
		}
	}
	
	if err := s.municipalityRepo.Update(ctx, municipality); err != nil {
		return fmt.Errorf("failed to update municipality: %w", err)
	}
	
	return nil
}

func (s *MunicipalityService) Delete(ctx context.Context, id int) error {
	// Check if municipality exists
	_, err := s.municipalityRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("municipality not found")
	}
	
	if err := s.municipalityRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete municipality: %w", err)
	}
	
	return nil
}