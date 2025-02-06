package client

import (
	"context"
	interactivev1 "github.com/Andras5014/webook/api/proto/gen/interactive/v1"
	"github.com/Andras5014/webook/interactive/domain"
	"github.com/Andras5014/webook/interactive/service"
	"google.golang.org/grpc"
)

// InteractiveServiceAdapter 本地实现伪装成grpc
type InteractiveServiceAdapter struct {
	svc service.InteractiveService
}

func NewInteractiveServiceAdapter(svc service.InteractiveService) *InteractiveServiceAdapter {
	return &InteractiveServiceAdapter{
		svc: svc,
	}
}
func (i *InteractiveServiceAdapter) IncrReadCnt(ctx context.Context, in *interactivev1.IncrReadCntRequest, opts ...grpc.CallOption) (*interactivev1.IncrReadCntResponse, error) {
	err := i.svc.IncrReadCnt(ctx, in.Biz, in.BizId)
	if err != nil {
		return nil, err
	}
	return &interactivev1.IncrReadCntResponse{}, nil
}

func (i *InteractiveServiceAdapter) Like(ctx context.Context, in *interactivev1.LikeRequest, opts ...grpc.CallOption) (*interactivev1.LikeResponse, error) {
	err := i.svc.Like(ctx, in.Biz, in.BizId, in.Uid)
	if err != nil {
		return nil, err
	}
	return &interactivev1.LikeResponse{}, nil
}

func (i *InteractiveServiceAdapter) CancelLike(ctx context.Context, in *interactivev1.CancelLikeRequest, opts ...grpc.CallOption) (*interactivev1.CancelLikeResponse, error) {
	err := i.svc.CancelLike(ctx, in.Biz, in.BizId, in.Uid)
	if err != nil {
		return nil, err
	}
	return &interactivev1.CancelLikeResponse{}, nil
}

func (i *InteractiveServiceAdapter) Collect(ctx context.Context, in *interactivev1.CollectRequest, opts ...grpc.CallOption) (*interactivev1.CollectResponse, error) {
	err := i.svc.Collect(ctx, in.Biz, in.BizId, in.Cid, in.Uid)
	if err != nil {
		return nil, err
	}
	return &interactivev1.CollectResponse{}, nil
}

func (i *InteractiveServiceAdapter) Get(ctx context.Context, in *interactivev1.GetRequest, opts ...grpc.CallOption) (*interactivev1.GetResponse, error) {
	intr, err := i.svc.Get(ctx, in.Biz, in.BizId, in.Uid)
	if err != nil {
		return nil, err
	}
	return &interactivev1.GetResponse{
		Intr: i.toDTO(intr),
	}, nil
}

func (i *InteractiveServiceAdapter) GetByIds(ctx context.Context, in *interactivev1.GetByIdsRequest, opts ...grpc.CallOption) (*interactivev1.GetByIdsResponse, error) {
	res, err := i.svc.GetByIds(ctx, in.Biz, in.BizIds)
	if err != nil {
		return nil, err
	}
	intrs := make(map[int64]*interactivev1.Interactive, len(res))
	for _, intr := range res {
		intrs[intr.BizId] = i.toDTO(intr)
	}
	return &interactivev1.GetByIdsResponse{
		Intrs: intrs,
	}, nil
}

func (i *InteractiveServiceAdapter) mustEmbedUnimplementedInteractiveServiceServer() {
	//TODO implement me
	panic("implement me")
}
func (i *InteractiveServiceAdapter) toDTO(intr domain.Interactive) *interactivev1.Interactive {
	return &interactivev1.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		CollectCnt: intr.CollectCnt,
		Collected:  intr.Collected,
		LikeCnt:    intr.LikeCnt,
		Liked:      intr.Liked,
		ReadCnt:    intr.ReadCnt,
	}
}
