package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/joaopanucci/apsdigital/internal/config"
	"github.com/joaopanucci/apsdigital/internal/domain/entities"
	"github.com/joaopanucci/apsdigital/internal/domain/repositories"
	"github.com/joaopanucci/apsdigital/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo         repositories.UserRepository
	refreshTokenRepo repositories.RefreshTokenRepository
	roleRepo         repositories.RoleRepository
	config          *config.Config
}

type Claims struct {
	UserID       uuid.UUID `json:"user_id"`
	CPF          string    `json:"cpf"`
	Role         string    `json:"role"`
	Level        int       `json:"level"`
	IsAuthorized bool      `json:"is_authorized"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	CPF      string `json:"cpf" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	User         *entities.User `json:"user"`
}

type RegisterRequest struct {
	Email          string    `json:"email" binding:"required,email"`
	Password       string    `json:"password" binding:"required,min=6"`
	Name           string    `json:"name" binding:"required"`
	CPF            string    `json:"cpf" binding:"required"`
	Phone          string    `json:"phone"`
	RoleID         uuid.UUID `json:"role_id" binding:"required"`
	ProfessionID   *int      `json:"profession_id"`
	MunicipalityID *int      `json:"municipality_id" binding:"required"`
	Unit           string    `json:"unit"`
}

func NewAuthService(userRepo repositories.UserRepository, refreshTokenRepo repositories.RefreshTokenRepository, roleRepo repositories.RoleRepository, config *config.Config) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		roleRepo:         roleRepo,
		config:          config,
	}
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*entities.User, error) {
	// Validate CPF
	if !utils.ValidateCPF(req.CPF) {
		return nil, fmt.Errorf("CPF inválido")
	}
	
	// Clean CPF for storage
	req.CPF = utils.CleanCPF(req.CPF)
	
	// Validate role - prevent ADM registration
	role, err := s.roleRepo.GetByID(ctx, req.RoleID)
	if err != nil {
		return nil, fmt.Errorf("invalid role")
	}
	
	// Block ADM role registration (level 1 is ADM)
	if role.Level == 1 || role.Name == "ADM" {
		return nil, fmt.Errorf("ADM role registration is not allowed")
	}

	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	existingUserCPF, _ := s.userRepo.GetByCPF(ctx, req.CPF)
	if existingUserCPF != nil {
		return nil, fmt.Errorf("user with this CPF already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &entities.User{
		Email:          req.Email,
		Password:       string(hashedPassword),
		Name:           req.Name,
		CPF:            req.CPF,
		Phone:          req.Phone,
		RoleID:         req.RoleID,
		ProfessionID:   req.ProfessionID,
		MunicipalityID: req.MunicipalityID,
		Unit:           req.Unit,
		Status:         entities.UserStatusPendingAuthorization,
		IsAuthorized:   false, // Requer autorização por padrão
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Validate and clean CPF
	if !utils.ValidateCPF(req.CPF) {
		return nil, fmt.Errorf("credenciais inválidas")
	}
	req.CPF = utils.CleanCPF(req.CPF)
	
	// Get user by CPF
	user, err := s.userRepo.GetByCPF(ctx, req.CPF)
	if err != nil {
		return nil, fmt.Errorf("credenciais inválidas")
	}

	// Check if user is active
	if user.Status != entities.UserStatusActive {
		return nil, fmt.Errorf("conta de usuário não está ativa")
	}

	// Check if user is authorized (new requirement)
	if !user.IsAuthorized {
		return nil, fmt.Errorf("seu cadastro ainda não foi autorizado")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Clear password from response
	user.Password = ""

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	// Validate refresh token
	token, err := s.refreshTokenRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if user.Status != entities.UserStatusActive {
		return nil, fmt.Errorf("user account is not active")
	}

	// Revoke old refresh token
	if err := s.refreshTokenRepo.RevokeToken(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to revoke old token: %w", err)
	}

	// Generate new tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Clear password from response
	user.Password = ""

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.refreshTokenRepo.RevokeToken(ctx, refreshToken)
}

func (s *AuthService) generateAccessToken(user *entities.User) (string, error) {
	expirationTime, err := time.ParseDuration(s.config.JWT.Expiration)
	if err != nil {
		expirationTime = 24 * time.Hour
	}

	claims := &Claims{
		UserID:       user.ID,
		CPF:          user.CPF,
		Role:         user.Role.Name,
		Level:        user.Role.Level,
		IsAuthorized: user.IsAuthorized,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *AuthService) generateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	// Generate random token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	tokenString := hex.EncodeToString(bytes)

	expirationTime, err := time.ParseDuration(s.config.JWT.RefreshTokenExpiration)
	if err != nil {
		expirationTime = 7 * 24 * time.Hour
	}

	// Revoke existing refresh tokens for this user
	if err := s.refreshTokenRepo.RevokeByUserID(ctx, userID); err != nil {
		return "", err
	}

	// Save refresh token
	refreshToken := &entities.RefreshToken{
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(expirationTime),
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return "", err
	}

	return tokenString, nil
}
