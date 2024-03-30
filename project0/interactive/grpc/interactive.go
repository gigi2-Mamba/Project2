package grpc

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
Date: 2024/3/7  周四
Time: 11:50
*/

type  InteractiveServiceServer struct {
	intrv1.UnimplementedInteractiveServiceServer
	svc service.InteractiveService
}

func NewInteractiveServiceServer(svc service.InteractiveService) *InteractiveServiceServer {
	return &InteractiveServiceServer{svc: svc}
}
// 反向控制，在自身可调用方，将自身注册在grpc.server 里面
func (i *InteractiveServiceServer) Register(s *grpc.Server )  {
	intrv1.RegisterInteractiveServiceServer(s,i)
}

func (i *InteractiveServiceServer) IncrReadCnt(ctx context.Context, request *intrv1.IncrReadCntRequest) (*intrv1.IncrReadCntResponse, error) {
	err := i.svc.IncrReadCnt(ctx, request.GetBiz(), request.GetBizId())
	log.Println("---------------进入grpc调动")
	if err != nil {
		return nil, err
	}
	return &intrv1.IncrReadCntResponse{},nil
}

func (i *InteractiveServiceServer) Like(ctx context.Context, request *intrv1.LikeRequest) (*intrv1.LikeResponse, error) {
	err := i.svc.Like(ctx,request.GetBiz(),request.GetBizId(),request.GetUid())
	if err != nil {
		return nil, err
	}

	return &intrv1.LikeResponse{},nil
}

func (i *InteractiveServiceServer) CancelLike(ctx context.Context, request *intrv1.CancelLikeRequest) (*intrv1.CancelLikeResponse, error) {
	err := i.svc.CancelLike(ctx,request.GetBiz(),request.GetBizId(),request.Uid)
	if err != nil {
		return nil, err
	}

	return &intrv1.CancelLikeResponse{},nil
}

func (i *InteractiveServiceServer) Collect(ctx context.Context, request *intrv1.CollectRequest) (*intrv1.CollectResponse, error) {
	err :=i.svc.Collect(ctx,request.GetBiz(),request.GetBizId(),request.GetCid(),request.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.CollectResponse{},nil
}

func (i *InteractiveServiceServer) Get(ctx context.Context, request *intrv1.GetRequest) (*intrv1.GetResponse, error) {
	res,err := i.svc.Get(ctx,request.Biz,request.GetBizid(),request.GetUid())
	log.Println("远程的  。。。。。。 grpc ")
	if err != nil {
		return nil,err
	}
	intr := i.toDTO(res)

	return &intrv1.GetResponse{
		Intr: intr,
	},nil

}

func (i *InteractiveServiceServer) GetByIds(ctx context.Context, request *intrv1.GetByIdsRequest) (*intrv1.GetByIdsResponse, error) {
	res, err := i.svc.GetByIds(ctx, request.GetBiz(), request.GetIds())
	if err != nil {
		return nil, err
	}
	intrs := make(map[int64]*intrv1.Interactive,len(res))
	for k, v := range res {
		intrs[k] = i.toDTO(v)
	}
	return &intrv1.GetByIdsResponse{
		Intrs: intrs,
	}, nil
}

// 额外写一个方法做dto
func (i *InteractiveServiceServer) toDTO(intr domain.Interactive) *intrv1.Interactive{
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


