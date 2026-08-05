package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "auth/gen/go/auth/v1"
	"auth/internal/config"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient authv1.AuthClient
}

func NewSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	// t.Parallel() // sqlite не поддерживает параллельные тесты

	err := godotenv.Load("../.env")
	if err != nil {
		t.Logf("can't load .env file")
	}

	cfg, err := config.LoadByPath("../config/default_config.yaml")
	if err != nil {
		t.Fatalf("load configuration")
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelFunc()
	})

	clientConn, err := grpc.DialContext(
		context.Background(),
		net.JoinHostPort("127.0.0.1", strconv.Itoa(cfg.GRPC.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("dial gRPC server")
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: authv1.NewAuthClient(clientConn),
	}
}
