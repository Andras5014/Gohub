package main

import (
	"context"
	interactivev1 "github.com/Andras5014/webook/api/proto/gen/interactive/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func TestGRPCClient(t *testing.T) {
	cc, err := grpc.Dial("localhost:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := interactivev1.NewInteractiveServiceClient(cc)
	resp, err := client.Get(context.Background(), &interactivev1.GetRequest{
		Biz:   "article",
		BizId: 1,
		Uid:   1,
	})
	require.NoError(t, err)
	t.Log(resp)
}
