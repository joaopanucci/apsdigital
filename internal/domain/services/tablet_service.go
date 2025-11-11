package services

import (
	"context"
	"fmt"
	"time"

	"apsdigital/internal/domain/entities"
	"apsdigital/internal/domain/repositories"
	
	"github.com/google/uuid"
)

type TabletService struct {
	tabletRepo repositories.TabletRepository
	userRepo   repositories.UserRepository
}

func NewTabletService(tabletRepo repositories.TabletRepository, userRepo repositories.UserRepository) *TabletService {
	return &TabletService{
		tabletRepo: tabletRepo,
		userRepo:   userRepo,
	}
}

func (s *TabletService) GetAll(ctx context.Context, filters map[string]interface{}) ([]*entities.Tablet, error) {
	tablets, err := s.tabletRepo.List(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get tablets: %w", err)
	}
	
	return tablets, nil
}

func (s *TabletService) GetByID(ctx context.Context, id int) (*entities.Tablet, error) {
	tablet, err := s.tabletRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tablet: %w", err)
	}
	
	return tablet, nil
}

func (s *TabletService) GetByUserCPF(ctx context.Context, cpf string) ([]*entities.Tablet, error) {
	tablets, err := s.tabletRepo.GetByUserCPF(ctx, cpf)
	if err != nil {
		return nil, fmt.Errorf("failed to get tablets by CPF: %w", err)
	}
	
	return tablets, nil
}

func (s *TabletService) GetByAssignedUser(ctx context.Context, userID uuid.UUID) ([]*entities.Tablet, error) {
	tablets, err := s.tabletRepo.GetByAssignedUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tablets by user: %w", err)
	}
	
	return tablets, nil
}

func (s *TabletService) Create(ctx context.Context, tablet *entities.Tablet) error {
	// Validate that the tablet doesn't already exist with the same serial number
	existing, _ := s.tabletRepo.List(ctx, map[string]interface{}{
		"serial_number": tablet.SerialNumber,
	})
	
	for _, existingTablet := range existing {
		if existingTablet.SerialNumber == tablet.SerialNumber {
			return fmt.Errorf("tablet with serial number '%s' already exists", tablet.SerialNumber)
		}
	}
	
	if err := s.tabletRepo.Create(ctx, tablet); err != nil {
		return fmt.Errorf("failed to create tablet: %w", err)
	}
	
	return nil
}

func (s *TabletService) AssignToUser(ctx context.Context, tabletID int, userID uuid.UUID, userCPF string) error {
	// Get tablet
	tablet, err := s.tabletRepo.GetByID(ctx, tabletID)
	if err != nil {
		return fmt.Errorf("tablet not found")
	}
	
	// Check if tablet is available
	if tablet.Status != entities.TabletStatusAvailable {
		return fmt.Errorf("tablet is not available for assignment")
	}
	
	// Verify user exists and is active
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	
	if user.Status != entities.UserStatusActive {
		return fmt.Errorf("user is not active")
	}
	
	// Assign tablet to user
	tablet.Status = entities.TabletStatusAssigned
	tablet.AssignedUserID = &userID
	tablet.UserCPF = &userCPF
	now := time.Now()
	tablet.AssignedAt = &now
	
	if err := s.tabletRepo.Update(ctx, tablet); err != nil {
		return fmt.Errorf("failed to assign tablet: %w", err)
	}
	
	return nil
}

func (s *TabletService) ReturnTablet(ctx context.Context, tabletID int) error {
	// Get tablet
	tablet, err := s.tabletRepo.GetByID(ctx, tabletID)
	if err != nil {
		return fmt.Errorf("tablet not found")
	}
	
	// Check if tablet is assigned
	if tablet.Status != entities.TabletStatusAssigned {
		return fmt.Errorf("tablet is not currently assigned")
	}
	
	// Return tablet (make it available)
	tablet.Status = entities.TabletStatusAvailable
	tablet.AssignedUserID = nil
	tablet.UserCPF = nil
	tablet.AssignedAt = nil
	
	if err := s.tabletRepo.Update(ctx, tablet); err != nil {
		return fmt.Errorf("failed to return tablet: %w", err)
	}
	
	return nil
}

func (s *TabletService) MarkAsMaintenance(ctx context.Context, tabletID int) error {
	// Get tablet
	tablet, err := s.tabletRepo.GetByID(ctx, tabletID)
	if err != nil {
		return fmt.Errorf("tablet not found")
	}
	
	// Mark as maintenance
	tablet.Status = entities.TabletStatusMaintenance
	
	// If it was assigned, unassign it
	if tablet.Status == entities.TabletStatusAssigned {
		tablet.AssignedUserID = nil
		tablet.UserCPF = nil
		tablet.AssignedAt = nil
	}
	
	if err := s.tabletRepo.Update(ctx, tablet); err != nil {
		return fmt.Errorf("failed to mark tablet as maintenance: %w", err)
	}
	
	return nil
}

func (s *TabletService) Update(ctx context.Context, tablet *entities.Tablet) error {
	// Check if tablet exists
	_, err := s.tabletRepo.GetByID(ctx, tablet.ID)
	if err != nil {
		return fmt.Errorf("tablet not found")
	}
	
	if err := s.tabletRepo.Update(ctx, tablet); err != nil {
		return fmt.Errorf("failed to update tablet: %w", err)
	}
	
	return nil
}

func (s *TabletService) Delete(ctx context.Context, id int) error {
	// Check if tablet exists
	_, err := s.tabletRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("tablet not found")
	}
	
	if err := s.tabletRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete tablet: %w", err)
	}
	
	return nil
}