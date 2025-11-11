package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/joaopanucci/apsdigital/internal/domain/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByCPF(ctx context.Context, cpf string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters map[string]interface{}) ([]*entities.User, error)
	GetPendingAuthorization(ctx context.Context) ([]*entities.User, error)
	AuthorizeUser(ctx context.Context, userID uuid.UUID) error
}

type MunicipalityRepository interface {
	Create(ctx context.Context, municipality *entities.Municipality) error
	GetByID(ctx context.Context, id int) (*entities.Municipality, error)
	GetByName(ctx context.Context, name string) (*entities.Municipality, error)
	List(ctx context.Context) ([]*entities.Municipality, error)
	Update(ctx context.Context, municipality *entities.Municipality) error
	Delete(ctx context.Context, id int) error
}

type TabletRepository interface {
	Create(ctx context.Context, tablet *entities.Tablet) error
	GetByID(ctx context.Context, id int) (*entities.Tablet, error)
	GetByUserCPF(ctx context.Context, cpf string) ([]*entities.Tablet, error)
	GetByAssignedUser(ctx context.Context, userID uuid.UUID) ([]*entities.Tablet, error)
	List(ctx context.Context, filters map[string]interface{}) ([]*entities.Tablet, error)
	Update(ctx context.Context, tablet *entities.Tablet) error
	Delete(ctx context.Context, id int) error
}

type RoleRepository interface {
	Create(ctx context.Context, role *entities.Role) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Role, error)
	GetByName(ctx context.Context, name string) (*entities.Role, error)
	List(ctx context.Context) ([]*entities.Role, error)
	Update(ctx context.Context, role *entities.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProfessionRepository interface {
	Create(ctx context.Context, profession *entities.Profession) error
	GetByID(ctx context.Context, id int) (*entities.Profession, error)
	GetByName(ctx context.Context, name string) (*entities.Profession, error)
	GetAll(ctx context.Context) ([]*entities.Profession, error)
	Update(ctx context.Context, profession *entities.Profession) error
	Delete(ctx context.Context, id int) error
}

type TabletRequestRepository interface {
	Create(ctx context.Context, request *entities.TabletRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.TabletRequest, error)
	List(ctx context.Context, filters map[string]interface{}) ([]*entities.TabletRequest, error)
	Update(ctx context.Context, request *entities.TabletRequest) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.TabletRequest, error)
}

type AuthorizationRepository interface {
	Create(ctx context.Context, auth *entities.Authorization) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Authorization, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.Authorization, error)
	List(ctx context.Context, status entities.AuthorizationStatus) ([]*entities.Authorization, error)
	Update(ctx context.Context, auth *entities.Authorization) error
}

type PaymentRepositoryInterface interface {
	Create(ctx context.Context, payment *entities.Payment) error
	GetByID(ctx context.Context, id uint) (*entities.Payment, error)
	GetAll(ctx context.Context, filters map[string]interface{}) ([]*entities.Payment, error)
	Update(ctx context.Context, payment *entities.Payment) error
	Delete(ctx context.Context, id uint) error
	GetCompetences(ctx context.Context, municipalityID *uint) ([]string, error)
	GetYears(ctx context.Context, municipalityID *uint) ([]int, error)
}

type ResolutionRepositoryInterface interface {
	Create(ctx context.Context, resolution *entities.Resolution) error
	GetByID(ctx context.Context, id uint) (*entities.Resolution, error)
	GetAll(ctx context.Context, filters map[string]interface{}) ([]*entities.Resolution, error)
	Update(ctx context.Context, resolution *entities.Resolution) error
	Delete(ctx context.Context, id uint) error
	GetTypes(ctx context.Context, municipalityID *uint) ([]string, error)
	GetYears(ctx context.Context, municipalityID *uint) ([]int, error)
	GetRecent(ctx context.Context, municipalityID *uint, limit int) ([]*entities.Resolution, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entities.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*entities.RefreshToken, error)
	RevokeByUserID(ctx context.Context, userID uuid.UUID) error
	RevokeToken(ctx context.Context, token string) error
	CleanupExpired(ctx context.Context) error
}
