package interceptor

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Моковый gRPC-хендлер
type mockUnaryHandler struct {
	mock.Mock
}

func (m *mockUnaryHandler) Handle(ctx context.Context, req interface{}) (interface{}, error) {
	args := m.Called(ctx, req)
	return args.Get(0), args.Error(1)
}

func TestNewAuthInterceptor(t *testing.T) {
	jwtKey := "test_jwt_key"
	ai := NewAuthInterceptor(jwtKey)

	assert.NotNil(t, ai)
	assert.Equal(t, []byte(jwtKey), ai.jwtKey)
}

func TestAuthInterceptor_Unary_SkipMethods(t *testing.T) {
	jwtKey := "test_jwt_key"
	ai := NewAuthInterceptor(jwtKey)

	handler := &mockUnaryHandler{}
	handler.On("Handle", mock.Anything, mock.Anything).Return("response", nil)

	ctx := context.Background()
	req := "request"

	// Проверяем RegisterUser метод
	info := &grpc.UnaryServerInfo{FullMethod: "/service.Service/RegisterUser"}
	resp, err := ai.Unary()(ctx, req, info, handler.Handle)

	assert.NoError(t, err)
	assert.Equal(t, "response", resp)
	handler.AssertCalled(t, "Handle", ctx, req)

	// Проверяем LoginUser метод
	info.FullMethod = "/service.Service/LoginUser"
	resp, err = ai.Unary()(ctx, req, info, handler.Handle)

	assert.NoError(t, err)
	assert.Equal(t, "response", resp)
	handler.AssertCalled(t, "Handle", ctx, req)
}

func TestAuthInterceptor_Unary_Unauthorized(t *testing.T) {
	jwtKey := "test_jwt_key"
	ai := NewAuthInterceptor(jwtKey)

	handler := &mockUnaryHandler{}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})
	req := "request"

	info := &grpc.UnaryServerInfo{FullMethod: "/service.Service/ProtectedMethod"}
	resp, err := ai.Unary()(ctx, req, info, handler.Handle)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestAuthInterceptor_Unary_Forbidden(t *testing.T) {
	jwtKey := "test_jwt_key"
	ai := NewAuthInterceptor(jwtKey)

	handler := &mockUnaryHandler{}

	claims := jwt.MapClaims{"other_key": "value"}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(jwtKey))

	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctx := metadata.NewIncomingContext(context.Background(), md)
	req := "request"

	info := &grpc.UnaryServerInfo{FullMethod: "/service.Service/ProtectedMethod"}
	resp, err := ai.Unary()(ctx, req, info, handler.Handle)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestAuthInterceptor_Unary_Success(t *testing.T) {
	jwtKey := "test_jwt_key"
	ai := NewAuthInterceptor(jwtKey)

	handler := &mockUnaryHandler{}
	handler.On("Handle", mock.Anything, mock.Anything).Return("response", nil)

	claims := jwt.MapClaims{"user_id": float64(123)}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(jwtKey))

	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctx := metadata.NewIncomingContext(context.Background(), md)
	req := "request"

	info := &grpc.UnaryServerInfo{FullMethod: "/service.Service/ProtectedMethod"}
	resp, err := ai.Unary()(ctx, req, info, handler.Handle)

	assert.NoError(t, err)
	assert.Equal(t, "response", resp)
	handler.AssertCalled(t, "Handle", mock.Anything, req)
}

func TestAuthInterceptor_validateToken(t *testing.T) {
	jwtKey := "test_jwt_key"
	ai := NewAuthInterceptor(jwtKey)

	claims := jwt.MapClaims{"user_id": float64(123)}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(jwtKey))

	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	parsedClaims, err := ai.validateToken(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, parsedClaims)
	assert.Equal(t, float64(123), (*parsedClaims)["user_id"])
}

func TestAuthInterceptor_validateToken_InvalidToken(t *testing.T) {
	jwtKey := "test_jwt_key"
	ai := NewAuthInterceptor(jwtKey)

	md := metadata.Pairs("authorization", "Bearer invalid_token")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	parsedClaims, err := ai.validateToken(ctx)
	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
	assert.Contains(t, err.Error(), "token is malformed")
}
