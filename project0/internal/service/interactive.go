package service

import (
	"context"
	"golang.org/x/sync/errgroup"
	"project0/internal/domain"
	"project0/internal/repository"
)

// Created by Changer on 2024/2/9.
// Copyright 2024 programmer.

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string,id int64, uid int64) error
	CancelLike(ctx context.Context, biz string,id int64, uid int64) error
	Collect(ctx context.Context, biz string,id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error)
}

type interactiveService struct {
	repo repository.InteractiveRepository
}

// 获取互动数据
func (i *interactiveService) Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error) {
	//互动详情总数
	inter ,err :=i.repo.Get(ctx,biz,id)
	if err != nil {
		return domain.Interactive{},err
	}
    var eg errgroup.Group
	eg.Go(func() error {
		var er error
		liked,er := i.repo.Liked(ctx,biz,id,uid)
		inter.Liked = liked
		return er
	})

	eg.Go(func() error {
		var er error
		collected ,er := i.repo.Collected(ctx,biz,id,uid)
		inter.Liked = collected

		return er
	})

	return inter,eg.Wait()

}

func (i *interactiveService) Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error {

	return i.repo.AddCollectItem(ctx,biz,id,cid,uid)
}

func (i *interactiveService) CancelLike(ctx context.Context,biz string, id int64, uid int64) error {
	return i.repo.DecrLikeCnt(ctx,biz,id,uid)
}

func (i *interactiveService) Like(ctx context.Context,biz string, id int64, uid int64) error {
	return  i.repo.IncrLikeCnt(ctx,biz,id,uid)
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{repo: repo}
}

func (i *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return i.repo.IncrReadCnt(ctx, biz, bizId)
}
