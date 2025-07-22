package usecase_test

import (
	"coffe/internal/auth"
	"coffe/internal/common"
	"coffe/internal/user/usecase"
	"coffe/internal/user/usecase/mocks"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login_Success(t *testing.T) {

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepository(ctrl)

	jwtService := auth.NewJWTService("your-test-secret", 5*time.Minute)

	authService := usecase.NewAuthService(mockUserRepo, jwtService, mockTokenRepo)

	userPassword := "pasword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)

	expectedUser := &common.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role: &common.Role{
			Name: "client",
		},
	}

	ctx := context.Background()

	mockUserRepo.EXPECT().GetUserByEmail(ctx, expectedUser.Email).Return(expectedUser, nil)

	mockTokenRepo.EXPECT().SetToken(ctx, expectedUser.ID.String(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	acessToken, refreshToken, userID, err := authService.Login(ctx, expectedUser.Email, userPassword)

	if err != nil {
		t.Errorf("Ожидалось отсутствияе отшибки а вернли : %v", err)
	}

	if acessToken == "" {
		t.Error("Ожидали access-токен, но получили пустую строку")
	}

	if refreshToken == "" {
		t.Error("Ожидали refresh-токен, но получили пустую строку")
	}

	if userID != expectedUser.ID {
		t.Errorf("Ожидали userID %v, но получили %v", expectedUser.ID, userID)

	}
}
