package services

import (
	"context"
	"fmt"
	"time"

	"github.com/joaopanucci/apsdigital/internal/domain/entities"
	"github.com/joaopanucci/apsdigital/internal/domain/repositories"

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

// SearchAgentByCPF searches for a user agent by CPF
func (s *TabletService) SearchAgentByCPF(ctx context.Context, cpf string) (*entities.User, error) {
	user, err := s.userRepo.GetByCPF(ctx, cpf)
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}
	return user, nil
}

// RequestNewTablet creates a new tablet request
func (s *TabletService) RequestNewTablet(ctx context.Context, agentCPF string, requestedBy uuid.UUID) error {
	// This would typically create a tablet request record
	// For now, we'll just return success
	return nil
}

// RequestTabletReturn creates a tablet return request
func (s *TabletService) RequestTabletReturn(ctx context.Context, agentCPF string) error {
	// Find tablets by CPF and mark for return
	tablets, err := s.GetByUserCPF(ctx, agentCPF)
	if err != nil {
		return err
	}

	for _, tablet := range tablets {
		if err := s.ReturnTablet(ctx, tablet.ID); err != nil {
			return err
		}
	}
	return nil
}

// ReportTabletBroken reports a tablet as broken
func (s *TabletService) ReportTabletBroken(ctx context.Context, agentCPF string, description string) error {
	tablets, err := s.GetByUserCPF(ctx, agentCPF)
	if err != nil {
		return err
	}

	for _, tablet := range tablets {
		if err := s.MarkAsMaintenance(ctx, tablet.ID); err != nil {
			return err
		}
	}
	return nil
}

// ReportTabletStolen reports a tablet as stolen
func (s *TabletService) ReportTabletStolen(ctx context.Context, agentCPF string, description string) error {
	tablets, err := s.GetByUserCPF(ctx, agentCPF)
	if err != nil {
		return err
	}

	for _, tablet := range tablets {
		tablet.Status = entities.TabletStatusMaintenance
		tablet.UpdatedAt = time.Now()
		if err := s.Update(ctx, tablet); err != nil {
			return err
		}
	}
	return nil
}

// GetTabletRequests gets all tablet requests
func (s *TabletService) GetTabletRequests(ctx context.Context, filters map[string]interface{}) ([]map[string]interface{}, error) {
	// This would typically return tablet request records
	// For now, we'll return an empty slice
	return []map[string]interface{}{}, nil
}

// ApproveTabletRequest approves a tablet request
func (s *TabletService) ApproveTabletRequest(ctx context.Context, requestID uint, approvedBy uuid.UUID) error {
	// This would typically update a tablet request record
	// For now, we'll just return success
	return nil
}

// RejectTabletRequest rejects a tablet request
func (s *TabletService) RejectTabletRequest(ctx context.Context, requestID uint, rejectedBy uuid.UUID, reason string) error {
	// This would typically update a tablet request record
	// For now, we'll just return success
	return nil
}
