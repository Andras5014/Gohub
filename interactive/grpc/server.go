package grpc

import (
	"context"
	interactivev1 "github.com/Andras5014/gohub/api/proto/gen/interactive/v1"
	"github.com/Andras5014/gohub/interactive/domain"
	"github.com/Andras5014/gohub/interactive/service"
	"google.golang.org/grpc"
)

type InteractiveServiceServer struct {
	svc service.InteractiveService
	interactivev1.UnimplementedInteractiveServiceServer
}

func (i *InteractiveServiceServer) Register(server *grpc.Server) {
	interactivev1.RegisterInteractiveServiceServer(server, i)
}
func NewInteractiveServiceServer(svc service.InteractiveService) *InteractiveServiceServer {
	return &InteractiveServiceServer{
		svc: svc,
	}
}
func (i *InteractiveServiceServer) IncrReadCnt(ctx context.Context, request *interactivev1.IncrReadCntRequest) (*interactivev1.IncrReadCntResponse, error) {
	err := i.svc.IncrReadCnt(ctx, request.Biz, request.BizId)
	if err != nil {
		return nil, err
	}
	return &interactivev1.IncrReadCntResponse{}, nil
}

func (i *InteractiveServiceServer) Like(ctx context.Context, request *interactivev1.LikeRequest) (*interactivev1.LikeResponse, error) {
	err := i.svc.Like(ctx, request.Biz, request.BizId, request.Uid)
	if err != nil {
		return nil, err
	}
	return &interactivev1.LikeResponse{}, nil
}

func (i *InteractiveServiceServer) CancelLike(ctx context.Context, request *interactivev1.CancelLikeRequest) (*interactivev1.CancelLikeResponse, error) {
	err := i.svc.CancelLike(ctx, request.Biz, request.BizId, request.Uid)
	if err != nil {
		return nil, err
	}
	return &interactivev1.CancelLikeResponse{}, nil
}

func (i *InteractiveServiceServer) Collect(ctx context.Context, request *interactivev1.CollectRequest) (*interactivev1.CollectResponse, error) {
	err := i.svc.Collect(ctx, request.Biz, request.BizId, request.Cid, request.Uid)
	if err != nil {
		return nil, err
	}
	return &interactivev1.CollectResponse{}, nil
}

func (i *InteractiveServiceServer) Get(ctx context.Context, request *interactivev1.GetRequest) (*interactivev1.GetResponse, error) {
	intr, err := i.svc.Get(ctx, request.Biz, request.BizId, request.Uid)
	if err != nil {
		return nil, err
	}
	return &interactivev1.GetResponse{
		Intr: i.toDTO(intr),
	}, nil
}

func (i *InteractiveServiceServer) GetByIds(ctx context.Context, request *interactivev1.GetByIdsRequest) (*interactivev1.GetByIdsResponse, error) {
	res, err := i.svc.GetByIds(ctx, request.Biz, request.BizIds)
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

func (i *InteractiveServiceServer) mustEmbedUnimplementedInteractiveServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (i *InteractiveServiceServer) toDTO(intr domain.Interactive) *interactivev1.Interactive {
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
