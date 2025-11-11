package services

import (
	"context"
	"fmt"

	"apsdigital/internal/domain/entities"
	"apsdigital/internal/domain/repositories"
	
	"github.com/google/uuid"
)

type UserAuthorizationService struct {
	userRepo repositories.UserRepository
}

func NewUserAuthorizationService(userRepo repositories.UserRepository) *UserAuthorizationService {
	return &UserAuthorizationService{
		userRepo: userRepo,
	}
}

func (s *UserAuthorizationService) GetPendingUsers(ctx context.Context) ([]*entities.User, error) {
	users, err := s.userRepo.GetPendingAuthorization(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending authorization users: %w", err)
	}
	
	return users, nil
}

func (s *UserAuthorizationService) AuthorizeUser(ctx context.Context, userID uuid.UUID) error {
	// Check if user exists and is pending authorization
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	
	if user.Status != entities.UserStatusPendingAuthorization {
		return fmt.Errorf("user is not pending authorization")
	}
	
	// Authorize user
	if err := s.userRepo.AuthorizeUser(ctx, userID); err != nil {
		return fmt.Errorf("failed to authorize user: %w", err)
	}
	
	return nil
}

func (s *UserAuthorizationService) RejectUser(ctx context.Context, userID uuid.UUID) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	
	if user.Status != entities.UserStatusPendingAuthorization {
		return fmt.Errorf("user is not pending authorization")
	}
	
	// Set user as inactive (rejection)
	user.Status = entities.UserStatusInactive
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to reject user: %w", err)
	}
	
	return nil
}

func (s *UserAuthorizationService) GetUsersByMunicipality(ctx context.Context, municipalityID int) ([]*entities.User, error) {
	filters := map[string]interface{}{
		"municipality_id": municipalityID,
		"status": entities.UserStatusActive,
	}
	
	users, err := s.userRepo.List(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by municipality: %w", err)
	}
	
	return users, nil
}