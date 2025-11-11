package services

import (
	"context"
	"fmt"

	"github.com/joaopanucci/apsdigital/internal/domain/entities"
	"github.com/joaopanucci/apsdigital/internal/domain/repositories"
)

type ProfessionService struct {
	professionRepo repositories.ProfessionRepository
}

func NewProfessionService(professionRepo repositories.ProfessionRepository) *ProfessionService {
	return &ProfessionService{
		professionRepo: professionRepo,
	}
}

func (s *ProfessionService) CreateProfession(ctx context.Context, profession *entities.Profession) error {
	// Check if profession already exists
	existing, _ := s.professionRepo.GetByName(ctx, profession.Name)
	if existing != nil {
		return fmt.Errorf("profession with name '%s' already exists", profession.Name)
	}

	return s.professionRepo.Create(ctx, profession)
}

func (s *ProfessionService) GetProfessionByID(ctx context.Context, id int) (*entities.Profession, error) {
	return s.professionRepo.GetByID(ctx, id)
}

func (s *ProfessionService) GetAllProfessions(ctx context.Context) ([]*entities.Profession, error) {
	return s.professionRepo.GetAll(ctx)
}

func (s *ProfessionService) UpdateProfession(ctx context.Context, profession *entities.Profession) error {
	// Check if profession exists
	existing, err := s.professionRepo.GetByID(ctx, profession.ID)
	if err != nil {
		return fmt.Errorf("profession not found")
	}

	// Check if name is being changed to an existing one
	if existing.Name != profession.Name {
		nameExists, _ := s.professionRepo.GetByName(ctx, profession.Name)
		if nameExists != nil {
			return fmt.Errorf("profession with name '%s' already exists", profession.Name)
		}
	}

	return s.professionRepo.Update(ctx, profession)
}

func (s *ProfessionService) DeleteProfession(ctx context.Context, id int) error {
	// Check if profession exists
	_, err := s.professionRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("profession not found")
	}

	return s.professionRepo.Delete(ctx, id)
}