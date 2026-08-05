package test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/utils/env"

	authv1 "auth/gen/go/auth/v1"
	"auth/internal/service/auth"
	"auth/test/suite"
)

func TestLogin_LoginExistUser_JWT(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	username := gofakeit.Username()
	password := gofakeit.Password(true, true, true, false, false, 12)
	email := gofakeit.Email()

	registerResp, err := st.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Username: username,
		Password: password,
		Email:    email,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, registerResp.GetUserId())

	loginResp, err := st.AuthClient.Login(ctx, &authv1.LoginRequest{
		Username: username,
		Password: password,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := loginResp.GetToken()
	assert.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(env.GetString("AUTH_JWT_SECRET", "default_secret_key")), nil
	})
	require.NoError(t, err)
	assert.NotNil(t, tokenParsed)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, username, claims["username"].(string))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, registerResp.GetUserId(), int64(claims["user_id"].(float64)))
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), int64(claims["exp"].(float64)), 1)
}

func TestLogin_WrongPassword_InvalidCredentials(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	username := gofakeit.Username()
	password := gofakeit.Password(true, true, true, false, false, 12)
	email := gofakeit.Email()

	registerResp, err := st.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Username: username,
		Password: password,
		Email:    email,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, registerResp.GetUserId())

	_, err = st.AuthClient.Login(ctx, &authv1.LoginRequest{
		Username: username,
		Password: password + "1",
	})
	require.Error(t, err)

	stCode := status.Code(err)
	assert.Equal(t, codes.Internal, stCode)
	assert.Contains(t, err.Error(), auth.ErrInvalidCredentials.Error())
}

func TestLogin_NonexistentUser_InvalidCredentials(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	username := gofakeit.Username()
	password := gofakeit.Password(true, true, true, false, false, 12)

	_, err := st.AuthClient.Login(ctx, &authv1.LoginRequest{
		Username: username,
		Password: password,
	})
	require.Error(t, err)

	stCode := status.Code(err)
	assert.Equal(t, codes.Internal, stCode)
	assert.Contains(t, err.Error(), auth.ErrInvalidCredentials.Error())
}
