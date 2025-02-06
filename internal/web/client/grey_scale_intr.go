package client

import (
	"context"
	interactivev1 "github.com/Andras5014/webook/api/proto/gen/interactive/v1"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"google.golang.org/grpc"
	"math/rand"
)

type GreyScaleInteractiveServiceClient struct {
	local  interactivev1.InteractiveServiceClient
	remote interactivev1.InteractiveServiceClient

	threshold *atomicx.Value[int32]
}

func NewGreyScaleInteractiveServiceClient(local interactivev1.InteractiveServiceClient, remote interactivev1.InteractiveServiceClient, threshold int32) *GreyScaleInteractiveServiceClient {
	return &GreyScaleInteractiveServiceClient{
		local:     local,
		remote:    remote,
		threshold: atomicx.NewValueOf(threshold),
	}
}
func (g *GreyScaleInteractiveServiceClient) IncrReadCnt(ctx context.Context, in *interactivev1.IncrReadCntRequest, opts ...grpc.CallOption) (*interactivev1.IncrReadCntResponse, error) {
	return g.client().IncrReadCnt(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) Like(ctx context.Context, in *interactivev1.LikeRequest, opts ...grpc.CallOption) (*interactivev1.LikeResponse, error) {
	return g.client().Like(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) CancelLike(ctx context.Context, in *interactivev1.CancelLikeRequest, opts ...grpc.CallOption) (*interactivev1.CancelLikeResponse, error) {
	return g.client().CancelLike(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) Collect(ctx context.Context, in *interactivev1.CollectRequest, opts ...grpc.CallOption) (*interactivev1.CollectResponse, error) {
	return g.client().Collect(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) Get(ctx context.Context, in *interactivev1.GetRequest, opts ...grpc.CallOption) (*interactivev1.GetResponse, error) {
	return g.client().Get(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) GetByIds(ctx context.Context, in *interactivev1.GetByIdsRequest, opts ...grpc.CallOption) (*interactivev1.GetByIdsResponse, error) {
	return g.client().GetByIds(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) UpdateThreshold(threshold int) {
	g.threshold.Store(int32(threshold))
}

func (g *GreyScaleInteractiveServiceClient) client() interactivev1.InteractiveServiceClient {
	threshold := g.threshold.Load()
	num := rand.Int31n(100)
	if num < threshold {
		return g.local
	}
	return g.remote
}
