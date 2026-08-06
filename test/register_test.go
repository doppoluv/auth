package test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authv1 "auth/gen/go/auth/v1"
	"auth/internal/service/auth"
	"auth/test/suite"
)

func TestRegister_CorrectCreds_Success(t *testing.T) {
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
}

func TestRegister_DuplicateUsername_UsernameAlreadyExists(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	username := gofakeit.Username()
	password := gofakeit.Password(true, true, true, false, false, 12)
	email1 := gofakeit.Email()

	registerResp, err := st.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Username: username,
		Password: password,
		Email:    email1,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, registerResp.GetUserId())

	email2 := gofakeit.Email()
	_, err = st.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Username: username,
		Password: password,
		Email:    email2,
	})
	require.Error(t, err)

	stCode := status.Code(err)
	assert.Equal(t, codes.Internal, stCode)
	assert.Contains(t, err.Error(), auth.ErrUsernameAlreadyExists.Error())
}

func TestRegister_DuplicateEmail_EmailAlreadyExists(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	username1 := gofakeit.Username()
	password := gofakeit.Password(true, true, true, false, false, 12)
	email := gofakeit.Email()

	registerResp, err := st.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Username: username1,
		Password: password,
		Email:    email,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, registerResp.GetUserId())

	username2 := gofakeit.Username()
	_, err = st.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Username: username2,
		Password: password,
		Email:    email,
	})
	require.Error(t, err)

	stCode := status.Code(err)
	assert.Equal(t, codes.Internal, stCode)
	assert.Contains(t, err.Error(), auth.ErrEmailAlreadyExists.Error())
}
