package auth

import (
	"context"
	"errors"

	"github.com/geraldiadityo/go-backend/internal/config"
	"github.com/geraldiadityo/go-backend/internal/helper"
	"github.com/geraldiadityo/go-backend/internal/modules/users"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials  = errors.New("username atau password salah")
	ErrInactiveUser        = errors.New("akun user tidak aktif")
	ErrInvalidRefreshToken = errors.New("refresh token tidak valid atau telah kedaluwarsa")
)

type Service interface {
	Login(ctx context.Context, req LoginRequest) (*TokenResponse, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (*TokenResponse, error)
	Logout(ctx context.Context, userID uint) error
	LogoutByRefreshToken(ctx context.Context, token string) error
}

type service struct {
	userService users.Service
	cfg         *config.AppConfig
}

func NewService(userService users.Service, cfg *config.AppConfig) Service {
	return &service{
		userService: userService,
		cfg:         cfg,
	}
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	user, err := s.userService.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.Status {
		return nil, ErrInactiveUser
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := helper.GenerateAccessToken(user.ID, user.Username, s.cfg.JWTSecret, s.cfg.JWTAccessExpirationMinutes)
	if err != nil {
		return nil, err
	}

	refreshToken, err := helper.GenerateRefreshToken(user.ID, user.Username, s.cfg.JWTSecret, s.cfg.JWTRefreshExpirationDays)
	if err != nil {
		return nil, err
	}

	err = s.userService.UpdateRefreshToken(ctx, user.ID, &refreshToken)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*TokenResponse, error) {
	// 1. Validasi keabsahan token JWT refresh token
	claims, err := helper.ValidateToken(req.RefreshToken, s.cfg.JWTSecret)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	// 2. Ambil user berdasarkan refresh token yang tersimpan di database
	user, err := s.userService.GetByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, err
	}

	// Pastikan ID user pada klaim sesuai dengan data di DB
	if user.ID != claims.UserID {
		return nil, ErrInvalidRefreshToken
	}

	if !user.Status {
		return nil, ErrInactiveUser
	}

	// 3. Mekanisme Rotasi: Generate pasangan Access Token dan Refresh Token baru
	newAccessToken, err := helper.GenerateAccessToken(user.ID, user.Username, s.cfg.JWTSecret, s.cfg.JWTAccessExpirationMinutes)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := helper.GenerateRefreshToken(user.ID, user.Username, s.cfg.JWTSecret, s.cfg.JWTRefreshExpirationDays)
	if err != nil {
		return nil, err
	}

	// 4. Update refresh token baru ke database
	err = s.userService.UpdateRefreshToken(ctx, user.ID, &newRefreshToken)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *service) Logout(ctx context.Context, userID uint) error {
	return s.userService.UpdateRefreshToken(ctx, userID, nil)
}

func (s *service) LogoutByRefreshToken(ctx context.Context, token string) error {
	user, err := s.userService.GetByRefreshToken(ctx, token)
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			return nil
		}
		return err
	}
	return s.userService.UpdateRefreshToken(ctx, user.ID, nil)
}

