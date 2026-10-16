package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/geraldiadityo/go-backend/internal/config"
	"github.com/geraldiadityo/go-backend/internal/helper"
	"github.com/geraldiadityo/go-backend/internal/modules/users"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials  = errors.New("username atau password salah")
	ErrInactiveUser        = errors.New("akun user tidak aktif")
	ErrInvalidRefreshToken = errors.New("refresh token tidak valid atau telah kedaluwarsa")
)

type ErrAccountLocked struct {
	WaitSecond int
}

type ErrWrongPassword struct {
	RemainingAttemps int
}

type Service interface {
	Login(ctx context.Context, req LoginRequest) (*TokenResponse, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (*TokenResponse, error)
	Logout(ctx context.Context, userID uint, rawToken string, expiresAt time.Time) error
}

type service struct {
	userService users.Service
	cfg         *config.AppConfig
	rdb         *redis.Client
}

func NewService(userService users.Service, cfg *config.AppConfig, rdb *redis.Client) Service {
	return &service{
		userService: userService,
		cfg:         cfg,
		rdb:         rdb,
	}
}

// error handler function
func (e *ErrAccountLocked) Error() string {
	return fmt.Sprintf("Akun terkunci semetara karena terlalu banyak percobaan gagal. Silahkan tunggu %d detik lagi", e.WaitSecond)
}

func (e *ErrWrongPassword) Error() string {
	return fmt.Sprintf("Username atau password salah. sisa kesempatan: %d kali sebelum akun dikunci", e.RemainingAttemps)
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*TokenResponse, error) {
	normalizedUsername := strings.ToLower(strings.TrimSpace(req.Username))
	blockedKey := fmt.Sprintf("auth:lockout:blocked:%s", normalizedUsername)
	attemptsKey := fmt.Sprintf("auth:lockout:attemps:%s", normalizedUsername)
	levelKey := fmt.Sprintf("auth:lockout:level:%s", normalizedUsername)

	ttl, err := s.rdb.TTL(ctx, blockedKey).Result()
	if err == nil && ttl > 0 {
		return nil, fmt.Errorf("Akun terkunci sementara karena terlalu banyak percobaan gagal tunggu %d detik lagi", int(ttl.Seconds()))
	}

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
		attemps, _ := s.rdb.Incr(ctx, attemptsKey).Result()
		s.rdb.Expire(ctx, attemptsKey, 15*time.Minute)
		if attemps >= 3 {
			level, _ := s.rdb.Incr(ctx, levelKey).Result()
			s.rdb.Expire(ctx, levelKey, 24*time.Hour)

			// Hitung durasi tunggu: 30s * 2^(level - 1)
			// Level 1: 30s, Level 2: 60s, Level 3: 120s...
			multiplier := 1 << (level - 1)
			waitSeconds := 30 * multiplier

			// Batas maksimal waktu tunggu (cap) misal 3600 detik (1 jam)

			if waitSeconds > 3600 {
				waitSeconds = 3600
			}

			lockDuration := time.Duration(waitSeconds) * time.Second

			s.rdb.Set(ctx, blockedKey, "locked", lockDuration)
			// Reset counter attempts ke 0 agar setelah waktu tunggu habis,
			// user kembali mendapatkan kesempatan 3 kali sebelum level-nya berlipat ganda lagi

			s.rdb.Del(ctx, attemptsKey)

			return nil, &ErrAccountLocked{
				WaitSecond: waitSeconds,
			}

		}
		return nil, &ErrWrongPassword{
			RemainingAttemps: int(3 - attemps),
		}
	}

	s.rdb.Del(ctx, blockedKey, attemptsKey, levelKey)

	accessToken, err := helper.GenerateAccessToken(user.ID, user.Username, s.cfg.JWTSecret, s.cfg.JWTAccessExpirationMinutes)
	if err != nil {
		return nil, err
	}

	refreshToken, err := helper.GenerateRefreshToken(user.ID, user.Username, s.cfg.JWTSecret, s.cfg.JWTRefreshExpirationDays)
	if err != nil {
		return nil, err
	}

	hashedRefreshToken := helper.HashTokenSHA256(refreshToken)

	err = s.userService.UpdateRefreshToken(ctx, user.ID, &hashedRefreshToken)
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
	user, err := s.userService.GetById(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, err
	}

	if user.RefreshToken == nil {
		return nil, ErrInvalidRefreshToken
	}

	if !helper.CompareTokenSHA256(*user.RefreshToken, req.RefreshToken) {
		return nil, ErrInvalidRefreshToken
	}

	if !user.Status {
		return nil, ErrInactiveUser
	}

	newAccessToken, err := helper.GenerateAccessToken(user.ID, user.Username, s.cfg.JWTSecret, s.cfg.JWTAccessExpirationMinutes)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := helper.GenerateRefreshToken(user.ID, user.Username, s.cfg.JWTSecret, s.cfg.JWTRefreshExpirationDays)
	if err != nil {
		return nil, err
	}

	hashedRefreshToken := helper.HashTokenSHA256(newRefreshToken)

	err = s.userService.UpdateRefreshToken(ctx, user.ID, &hashedRefreshToken)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *service) Logout(ctx context.Context, userID uint, rawToken string, expiresAt time.Time) error {
	remainingTTL := time.Until(expiresAt)
	if remainingTTL > 0 {
		tokenHash := helper.HashTokenSHA256(rawToken)
		blacklistKey := fmt.Sprintf("blacklist:token:%s", tokenHash)

		if err := s.rdb.Set(ctx, blacklistKey, "revoked", remainingTTL).Err(); err != nil {
			return err
		}
	}

	return s.userService.UpdateRefreshToken(ctx, userID, nil)
}
