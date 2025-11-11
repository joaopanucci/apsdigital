package services

import (
	"context"
	"errors"

	"github.com/joaopanucci/apsdigital/internal/domain/entities"
	"github.com/joaopanucci/apsdigital/internal/infra/repositories"
)

type PaymentService struct {
	paymentRepo *repositories.PaymentRepository
}

func NewPaymentService(paymentRepo *repositories.PaymentRepository) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
	}
}

func (s *PaymentService) CreatePayment(ctx context.Context, payment *entities.Payment) error {
	if payment.FileURL == "" {
		return errors.New("file URL is required")
	}
	
	if payment.Competence == "" {
		return errors.New("competence is required")
	}
	
	if payment.UploadedBy == 0 {
		return errors.New("uploaded by user is required")
	}
	
	return s.paymentRepo.Create(ctx, payment)
}

func (s *PaymentService) GetPaymentByID(ctx context.Context, id uint) (*entities.Payment, error) {
	if id == 0 {
		return nil, errors.New("invalid payment ID")
	}
	
	return s.paymentRepo.GetByID(ctx, id)
}

func (s *PaymentService) GetAllPayments(ctx context.Context, filters map[string]interface{}) ([]*entities.Payment, error) {
	return s.paymentRepo.GetAll(ctx, filters)
}

func (s *PaymentService) GetPaymentsByMunicipality(ctx context.Context, municipalityID uint, filters map[string]interface{}) ([]*entities.Payment, error) {
	if filters == nil {
		filters = make(map[string]interface{})
	}
	
	filters["municipality_id"] = municipalityID
	return s.paymentRepo.GetAll(ctx, filters)
}

func (s *PaymentService) UpdatePayment(ctx context.Context, payment *entities.Payment) error {
	if payment.ID == 0 {
		return errors.New("invalid payment ID")
	}
	
	// Check if payment exists
	_, err := s.paymentRepo.GetByID(ctx, payment.ID)
	if err != nil {
		return errors.New("payment not found")
	}
	
	return s.paymentRepo.Update(ctx, payment)
}

func (s *PaymentService) DeletePayment(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid payment ID")
	}
	
	// Check if payment exists
	_, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("payment not found")
	}
	
	return s.paymentRepo.Delete(ctx, id)
}

func (s *PaymentService) GetCompetences(ctx context.Context, municipalityID *uint) ([]string, error) {
	return s.paymentRepo.GetCompetences(ctx, municipalityID)
}

func (s *PaymentService) GetYears(ctx context.Context, municipalityID *uint) ([]int, error) {
	return s.paymentRepo.GetYears(ctx, municipalityID)
}

func (s *PaymentService) GetPaymentsWithPagination(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*entities.Payment, error) {
	if page <= 0 {
		page = 1
	}
	
	if limit <= 0 {
		limit = 10
	}
	
	offset := (page - 1) * limit
	
	if filters == nil {
		filters = make(map[string]interface{})
	}
	
	filters["limit"] = limit
	filters["offset"] = offset
	
	return s.paymentRepo.GetAll(ctx, filters)
}