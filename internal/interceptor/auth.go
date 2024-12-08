package interceptor

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

type AuthInterceptor struct {
	jwtKey []byte
}

func NewAuthInterceptor(jwtKey string) *AuthInterceptor {
	return &AuthInterceptor{
		jwtKey: []byte(jwtKey),
	}
}

func (ai *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if strings.Contains(info.FullMethod, "RegisterUser") {
			return handler(ctx, req)
		}
		if strings.Contains(info.FullMethod, "LoginUser") {
			return handler(ctx, req)
		}

		claims, err := ai.validateToken(ctx)
		if err != nil {
			return nil, status.Errorf(401, "unauthorized: %v", err)
		}

		value, ok := (*claims)["user_id"]
		if !ok {
			return nil, status.Errorf(403, "forbidden")
		}

		userID, ok := value.(float64)
		if !ok {
			return nil, status.Errorf(403, "forbidden")
		}

		ctx = context.WithValue(ctx, "user_id", uint64(userID))

		return handler(ctx, req)
	}
}

func (ai *AuthInterceptor) validateToken(ctx context.Context) (*jwt.MapClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return ai.jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return &claims, nil
}
