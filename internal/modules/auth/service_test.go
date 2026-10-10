package auth

import (
	"context"
	"testing"
	"time"

	"github.com/geraldiadityo/go-backend/internal/config"
	"github.com/geraldiadityo/go-backend/internal/helper"
	"github.com/geraldiadityo/go-backend/internal/modules/users"
	"golang.org/x/crypto/bcrypt"
)

type mockUserService struct {
	user *users.User
}

func (m *mockUserService) CreateUser(ctx context.Context, req users.UserCreateDTO) (*users.UserResponse, error) {
	return nil, nil
}
func (m *mockUserService) GetAllUser(ctx context.Context, query users.UserQueryDTO) ([]users.UserResponse, helper.Meta, error) {
	return nil, helper.Meta{}, nil
}
func (m *mockUserService) GetByUsername(ctx context.Context, username string) (*users.User, error) {
	if m.user != nil && m.user.Username == username {
		return m.user, nil
	}
	return nil, users.ErrNotFound
}
func (m *mockUserService) GetByRefreshToken(ctx context.Context, token string) (*users.User, error) {
	if m.user != nil && m.user.RefreshToken != nil && *m.user.RefreshToken == token {
		return m.user, nil
	}
	return nil, users.ErrNotFound
}
func (m *mockUserService) GetById(ctx context.Context, id uint) (*users.User, error) {
	return nil, nil
}
func (m *mockUserService) UpdateUser(ctx context.Context, id uint, req users.UserUpdateRequest) (*users.UserResponse, error) {
	return nil, nil
}
func (m *mockUserService) UpdateRefreshToken(ctx context.Context, userID uint, token *string) error {
	if m.user != nil && m.user.ID == userID {
		m.user.RefreshToken = token
		return nil
	}
	return users.ErrNotFound
}
func (m *mockUserService) GetUserById(ctx context.Context, id uint) (*users.UserResponse, error) {
	return nil, nil
}
func (m *mockUserService) DeleteUser(ctx context.Context, id uint) error {
	return nil
}
func (m *mockUserService) ChangeStatus(ctx context.Context, id uint) (*users.UserResponse, error) {
	return nil, nil
}

func TestAuthService_Login_And_RefreshToken_Rotation(t *testing.T) {
	pwd, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	mockUser := &users.User{
		ID:           1,
		Username:     "testuser",
		Password:     string(pwd),
		Status:       true,
		RefreshToken: nil,
	}

	mockService := &mockUserService{user: mockUser}
	cfg := &config.AppConfig{
		JWTSecret:                  "testsecretkey",
		JWTAccessExpirationMinutes: 15,
		JWTRefreshExpirationDays:   7,
	}

	authSvc := NewService(mockService, cfg)
	ctx := context.Background()

	// 1. Test Login
	loginResp, err := authSvc.Login(ctx, LoginRequest{
		Username: "testuser",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("expected login success, got err: %v", err)
	}

	if loginResp.AccessToken == "" || loginResp.RefreshToken == "" {
		t.Fatalf("tokens should not be empty")
	}

	if mockUser.RefreshToken == nil || *mockUser.RefreshToken != loginResp.RefreshToken {
		t.Fatalf("expected user refresh token in db to be updated")
	}

	// Wait 1 second to ensure new JWTs have different iat
	time.Sleep(1 * time.Second)

	// 2. Test Refresh Token (Rotation)
	firstAccessToken := loginResp.AccessToken
	firstRefreshToken := loginResp.RefreshToken

	refreshResp, err := authSvc.RefreshToken(ctx, RefreshTokenRequest{
		RefreshToken: firstRefreshToken,
	})
	if err != nil {
		t.Fatalf("expected refresh token success, got err: %v", err)
	}

	if refreshResp.AccessToken == firstAccessToken {
		t.Errorf("expected new access token to be different from previous")
	}

	if refreshResp.RefreshToken == firstRefreshToken {
		t.Errorf("expected new refresh token to be different from previous (rotation)")
	}

	if mockUser.RefreshToken == nil || *mockUser.RefreshToken != refreshResp.RefreshToken {
		t.Errorf("expected user refresh token in db to match new rotated refresh token")
	}

	// 3. Test old refresh token should fail now (since db token was rotated)
	_, err = authSvc.RefreshToken(ctx, RefreshTokenRequest{
		RefreshToken: firstRefreshToken,
	})
	if err == nil {
		t.Errorf("expected old refresh token to fail, but it succeeded")
	}

	// 4. Test Logout
	err = authSvc.Logout(ctx, mockUser.ID)
	if err != nil {
		t.Fatalf("expected logout success, got: %v", err)
	}

	if mockUser.RefreshToken != nil {
		t.Errorf("expected refresh token to be null after logout")
	}
}

