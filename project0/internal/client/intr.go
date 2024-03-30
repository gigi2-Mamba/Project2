package client

import (
	"context"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	intrv1 "project0/api/proto/gen/api/proto/intr/v1"
)

/*
User: society-programmer
Date: 2024/3/8  周五
Time: 21:46
*/
// 为了方便迁移，做一个本地和远程切换
type InteractiveClient struct {
	remote  intrv1.InteractiveServiceClient
	local   intrv1.InteractiveServiceClient
    // 为了避免并发问题用原子类型
	threshold *atomicx.Value[int32]
}

func NewInteractiveClient(remote intrv1.InteractiveServiceClient, local intrv1.InteractiveServiceClient) *InteractiveClient {
	return &InteractiveClient{
		remote: remote,
		local: local,
	    threshold: atomicx.NewValue[int32](),
	}
}

func (i *InteractiveClient) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCntResponse, error) {
	return i.selectClient().IncrReadCnt(ctx, in, opts...)

}

func (i *InteractiveClient) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	return i.selectClient().Like(ctx, in, opts...)
}

func (i *InteractiveClient) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	return i.selectClient().CancelLike(ctx, in, opts...)
}

func (i *InteractiveClient) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	return i.selectClient().Collect(ctx, in, opts...)
}

func (i *InteractiveClient) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	return i.selectClient().Get(ctx, in, opts...)
}

func (i *InteractiveClient) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	return i.selectClient().GetByIds(ctx, in, opts...)
}

func (i *InteractiveClient) selectClient() intrv1.InteractiveServiceClient  {
   // [0,100]
    num := rand.Int31n(100)
	log.Println(" numllll",num)
	log.Println("i.threshold.Load() ",i.threshold.Load())
	// 初始化  threshold 为0的话岂不是一直local
	if num < i.threshold.Load() {
		log.Println("remote change success")
		return i.remote
	}
	return i.local

}


func (i *InteractiveClient) UpdateThreshold(val int32) {
    i.threshold.Store(val)
}




