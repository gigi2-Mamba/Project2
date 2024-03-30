package client

import (
	"context"
	"google.golang.org/grpc"
	"log"
	intrv1 "project0/api/proto/gen/api/proto/intr/v1"
	"project0/interactive/domain"
	"project0/interactive/service"
)

/*
User: society-programmer
Date: 2024/3/8  周五
Time: 19:41
*/

// 对接口装饰，然后又做一个适应
type LocalInteractiveServiceAdapter struct {
	svc service.InteractiveService
}

func NewLocalInteractiveServiceAdapter(svc service.InteractiveService) *LocalInteractiveServiceAdapter {
	return &LocalInteractiveServiceAdapter{svc: svc}
}

func (l *LocalInteractiveServiceAdapter) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCntResponse, error) {

	err := l.svc.IncrReadCnt(ctx, in.GetBiz(), in.GetBizId())

	return &intrv1.IncrReadCntResponse{},err
}

func (l *LocalInteractiveServiceAdapter) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	err := l.svc.Like(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrv1.LikeResponse{},err
}

func (l *LocalInteractiveServiceAdapter) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	err := l.svc.CancelLike(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrv1.CancelLikeResponse{},err
}

func (l *LocalInteractiveServiceAdapter) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	err :=l.svc.Collect(ctx,in.GetBiz(),in.GetBizId(),in.GetCid(),in.GetUid())

	return &intrv1.CollectResponse{},err
}

func (l *LocalInteractiveServiceAdapter) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	res, err := l.svc.Get(ctx, in.GetBiz(), in.GetBizid(), in.Uid)
	log.Println("本地的grpc")
	return &intrv1.GetResponse{
		Intr: l.toDTO(res),
	},err
}

func (l *LocalInteractiveServiceAdapter) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	res, err := l.svc.GetByIds(ctx, in.GetBiz(), in.GetIds())
	intrs := make(map[int64]*intrv1.Interactive,len(res))
	for k,v := range res {
		intr :=l.toDTO(v)
		intrs[k] = intr
	}
	return &intrv1.GetByIdsResponse{
		Intrs: intrs,
	},err
}


func (l *LocalInteractiveServiceAdapter) toDTO(intr domain.Interactive) *intrv1.Interactive{
	return &intrv1.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		ReadCnt:    intr.ReadCnt,
		CollectCnt: intr.CollectCnt,
		Collected:  intr.Collected,
		Liked:      intr.Liked,
		LikeCnt:    intr.LikeCnt,
	}
}

